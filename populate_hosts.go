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

const defaultHostsUrl = "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"

func populateHosts() {
	log.Info("Populating hosts...")

	res, err := http.Get(hostsUrl)
	if err != nil {
		panic("unable to fetch hosts source")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("unable to read hosts source")
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		writeHashedKey(line)
	}

	log.Info("Populated hosts.")
}

func writeHashedKey(line string) {
	if strings.HasPrefix(line, "0.0.0.0") {
		rawHost := strings.ReplaceAll(line, "0.0.0.0 ", "")
		hashedHost := hashAndEncodeHost(rawHost)
		host := Host{Hostname: hashedHost}
		err := app.db.Where(Host{Hostname: hashedHost}).FirstOrCreate(&host).Error
		if err != nil {
			panic("unable to write host to database: " + rawHost)
		}
		fmt.Println("wrote: " + rawHost + " as " + hashedHost)
	}
}

func hashAndEncodeHost(rawHost string) string {
	hashKey := getEnv("HASH_KEY", "")
	hash := hmac.New(sha256.New, []byte(hashKey))
	hash.Write([]byte(rawHost))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
