package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"net/http"
	"net/url"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	log.Info("Connecting to database...")
	db, err = gorm.Open(sqlite.Open("hosts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	log.Info("Connect to database.")
	log.Info("Migrating database...")
	db.AutoMigrate(&Host{})
	log.Info("Migrated database.")
	populateHosts()

	http.HandleFunc("/allow", allowHandler)
	log.Info("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

type Host struct {
	gorm.Model
	Hostname string `gorm:"index" db:"hostname"`
}

type Response struct {
	Allow bool `json:"allow"`
}

func allowHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(r)
	w.Header().Set("Content-Type", "application/json")
	encodedUrl := r.URL.Query().Get("url") // base64 encoded url
	ok, allow := allowUrl(encodedUrl)

	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	res := Response{Allow: allow}
	js, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
	log.Info(w)
}

func allowUrl(encodedUrl string) (bool, bool) {
	if encodedUrl == "" {
		return false, false
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedUrl)
	if err != nil {
		log.Error(err)
		return false, false
	}

	u, err := url.Parse(string(decoded))
	if err != nil {
		log.Error(err)
		return false, false
	}

	host := Host{}
	db.Where("hostname = ?", u.Host).First(&host)
	if host.ID != 0 {
		log.Info("host.id not equal to zero")
		return true, false
	}

	return true, true
}

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
			host := Host{Hostname: hostname}
			err := db.Where(Host{Hostname: hostname}).FirstOrCreate(&host).Error
			if err != nil {
				panic("unable to write host to database: " + hostname)
			}
			fmt.Println("wrote: " + hostname)
		}
	}

	log.Info("Populated hosts.")
}
