package main

import (
	"flag"

	arg "github.com/alexflint/go-arg"
	browser "github.com/blippar/balrog"
)

// Version is the software version injected at build time
var Version = "unknown"

// default... defines the default cli arguments for the software
const (
	defaultConfig = "config.json"
)

// cliArgs defines the list of potential cli arguments the software takes
type cliArgs struct {
	Config string `arg:"-c,help:path to your config file [env: BALROG_CONFIG],env:BALROG_CONFIG"`
}

func (cliArgs) Version() string {
	return "balrog " + Version
}

func main() {

	args := &cliArgs{
		Config: defaultConfig,
	}
	arg.MustParse(args)

	flag.Parse()
	browser.NewServer(args.Config).Run()
}
