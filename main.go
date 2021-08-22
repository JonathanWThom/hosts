package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"net/http"
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
	//populateHosts()

	http.HandleFunc("/allow", allowHandler)
	addr := getEnv("PORT", "8080")
	log.Info("Listening on port %v...", addr)
	http.ListenAndServe(addr, nil)
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
	encodedUrl := r.URL.Query().Get("url") // sha256 + base64 encoded url
	finalUrl := strings.ReplaceAll(encodedUrl, " ", "+")
	ok, allow := allowUrl(finalUrl)

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

func allowUrl(encodedAndHashedHost string) (bool, bool) {
	//u, err := url.Parse(encodedAndHashedUrl)
	//if err != nil {
	//log.Error(err)
	//return false, false
	//}

	host := Host{}
	db.Where("hostname = ?", encodedAndHashedHost).First(&host)
	if host.ID != 0 {
		log.Info("host.id not equal to zero")
		return true, false
	}

	return true, true
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
			hash := hmac.New(sha256.New, nil)
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
