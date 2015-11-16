package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

func main() {
	var app = cli.NewApp()
	app.Name = "cf"
	app.Version = "0.1.0"
	app.Usage = "Codeforces client"
	app.Action = func(c *cli.Context) {
		if c.Bool("version") {
			cli.ShowVersion(c)
			return
		}
		cli.ShowAppHelp(c)
	}
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("verbose") {
			log.SetLevel(log.InfoLevel)
		}
		return nil
	}
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
			Name:      "config",
			Usage:     "Set or show settings",
			ArgsUsage: "[<key> <value>]",
			Action:    config,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "list",
					Usage: "Display current settings",
				},
				cli.BoolFlag{
					Name:  "global",
					Usage: "Settings for ~/.cf.yml file",
				},
			},
		},
	}
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Display detailed information of the operations",
		},
		cli.BoolFlag{
			Name:  "version, V",
			Usage: "print the version",
		},
	}

	app.Run(os.Args)
}
