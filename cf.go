package main

import (
	logger "log"
	"os"

	"github.com/codegangsta/cli"
)

var log *logger.Logger

func init() {
	log = logger.New(os.Stderr, "", 0)
}

func main() {
	var app = cli.NewApp()
	app.Name = "cf"
	app.Version = "0.1.0"
	app.Usage = "Codeforces client"
	app.Commands = []cli.Command{
		{
			Name:      "parse",
			Usage:     "Parses codeforces problem",
			ArgsUsage: "<Problem URL>",
			Action:    parse,
		},
	}

	app.Run(os.Args)
}
