package main

import (
	"flag"

	"github.com/0rax/apk"
)

var configPath = flag.String("config", "", "path to your config file")

func main() {

	flag.Parse()
	apk.NewServer(*configPath).Run()
}
