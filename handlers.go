package main

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	Hostname string `gorm:"index" db:"hostname"`
}

type Response struct {
	Allow bool `json:"allow"`
}

func allowUrl(encodedAndHashedHost string) (bool, bool) {
	//u, err := url.Parse(encodedAndHashedUrl)
	//if err != nil {
	//log.Error(err)
	//return false, false
	//}

	host := Host{}
	app.db.Where("hostname = ?", encodedAndHashedHost).First(&host)
	if host.ID != 0 {
		log.Info("host.id not equal to zero")
		return true, false
	}

	return true, true
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
