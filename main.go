package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"net/http"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type App struct {
	db *gorm.DB
}

func newApp() *App {
	return &App{
		db: newDB(),
	}
}

func (app *App) Serve() {
	http.HandleFunc("/allow", allowHandler)
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

func main() {
	app := newApp()
	//populateHosts()
	app.Serve()
}

// Unused for now, maybe add a refresh endpoint at some point
func populateHosts() {
	log.Info("Populating hosts...")
	rawHostsUrl := "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn/hosts"
	res, err := http.Get(rawHostsUrl)
	defer res.Body.Close()

	if err != nil {
		panic("unable to fetch hosts source")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("unable to read hosts source")
	}
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "0.0.0.0") {
			hostname := strings.ReplaceAll(line, "0.0.0.0 ", "")
			// add shared salt key
			key := getEnv("HASH_KEY", "")
			hash := hmac.New(sha256.New, []byte(key))
			hash.Write([]byte(hostname))
			encoded := base64.StdEncoding.EncodeToString(hash.Sum(nil))
			host := Host{Hostname: encoded}
			err = db.Where(Host{Hostname: encoded}).FirstOrCreate(&host).Error
			if err != nil {
				panic("unable to write host to database: " + hostname)
			}
			fmt.Println("wrote: " + hostname + " as " + encoded)
		}
	}

	log.Info("Populated hosts.")
}
