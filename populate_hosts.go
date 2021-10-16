package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const defaultHostsUrl = "https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/porn/hosts"

func populateHosts() {
	log.Info("Populating hosts...")
	res, err := http.Get(hostsUrl)
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
			err = app.db.Where(Host{Hostname: encoded}).FirstOrCreate(&host).Error
			if err != nil {
				panic("unable to write host to database: " + hostname)
			}
			fmt.Println("wrote: " + hostname + " as " + encoded)
		}
	}

	log.Info("Populated hosts.")
}