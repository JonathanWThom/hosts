package main

import (
	"flag"
)

var app *App
var populate bool
var hostsUrl string
var hashKey string

func main() {
	app = newApp()

	flag.BoolVar(&populate, "p", false, "Add -p to refresh blocked hosts database")
	flag.StringVar(&hostsUrl, "h", defaultHostsUrl, "Url of hosts list to block, only useful if also passing p flag")
	flag.StringVar(&hashKey, "key", "", "Key used to hash host values with - client will need this as well. Only useful if also passing p flag.")
	flag.Parse()

	if populate {
		populateHosts()
	}

	app.Serve()
}
