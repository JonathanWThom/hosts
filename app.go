package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	routes []Route
}

type Route struct {
	path       string
	handleFunc func(http.ResponseWriter, *http.Request)
}

func newApp() *App {
	return &App{
		db:     newDB(),
		routes: []Route{{path: "/allow", handleFunc: allowHandler}},
	}
}

func (app *App) Serve() {
	for _, route := range app.routes {
		http.HandleFunc(route.path, route.handleFunc)
	}
	addr := getEnv("PORT", "8080")
	log.Infof("Listening on port %v...", addr)
	http.ListenAndServe(":"+addr, nil)
}

func newDB() *gorm.DB {
	log.Info("Connecting to database...")
	db, err := gorm.Open(sqlite.Open("hosts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Info("Connect to database.")
	log.Info("Migrating database...")
	db.AutoMigrate(&Host{})
	log.Info("Migrated database.")

	return db
}
