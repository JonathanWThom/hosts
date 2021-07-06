package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("hosts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Host{})
	populateHosts()

	http.HandleFunc("/allow", allowHandler)
	http.ListenAndServe(":8080", nil)
}

type Host struct {
	gorm.Model
	host string
}

type Response struct {
	Allow bool `json:"allow"`
}

func allowHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func allowUrl(encodedUrl string) (bool, bool) {
	if encodedUrl == "" {
		return false, false
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedUrl)
	if err != nil {
		return false, false
	}

	u, err := url.Parse(string(decoded))
	if err != nil {
		return false, false
	}

	host := Host{host: u.Host}
	db.Find(&host)

	if host.ID != 0 {
		return true, false
	}

	return true, true
}

func populateHosts() {

}
