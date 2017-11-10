package main

import (
	"flag"

	"github.com/blippar/alpine-package-browser"
)

var configPath = flag.String("config", "", "path to your config file")

func main() {

	flag.Parse()
	browser.NewServer(*configPath).Run()
}
