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
			ArgsUsage: "<ProblemURL>",
			Action:    parse,
		},
		{
			Name:      "setup",
			Usage:     "Setup environment for contest or single problem",
			ArgsUsage: "<ContestID | ProblemURL>",
			Action:    setup,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "lang",
					Usage: "programming language for sample solutions",
				},
			},
		},
		{
			Name:      "gen",
			Usage:     "Generates sample solution",
			ArgsUsage: "<source_file.ext>",
			Action:    gen,
		},
		{
			Name:   "test",
			Usage:  "Runs solution against test cases",
			Action: test,
		},
		{
			Name:      "set",
			Usage:     "Sets settings values",
			ArgsUsage: "<key> <value>",
			Action:    set,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "global",
					Usage: "Settings for ~/.cf.yml file",
				},
			},
		},
	}

	app.Run(os.Args)
}
