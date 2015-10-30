package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	log.SetFlags(0)
	var app = cli.NewApp()
	app.Name = "tp"
	app.Usage = "test your program"
	app.Commands = []cli.Command{
		{
			Name:   "parse",
			Usage:  "parse codeforces problem",
			Action: parse,
		},
	}

	app.Run(os.Args)
}
