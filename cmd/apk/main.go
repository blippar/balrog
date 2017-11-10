package main

import (
	"flag"

	browser "github.com/blippar/balrog"
)

var configPath = flag.String("config", "", "path to your config file")

func main() {

	flag.Parse()
	browser.NewServer(*configPath).Run()
}
