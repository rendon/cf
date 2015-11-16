package main

import (
	"fmt"
	"os/user"
	"strconv"

	"github.com/codegangsta/cli"
)

func config(c *cli.Context) {
	if !c.Bool("list") && len(c.Args()) != 2 {
		cli.ShowCommandHelp(c, "config")
		return
	}

	var file = ".settings.yml"
	if c.Bool("global") {
		user, err := user.Current()
		if err != nil {
			log.Fatalf("Failed to get user's home directory: %s", err)
		}
		file = user.HomeDir + "/.cf.yml"
	}
	settings, err := ReadKeyValueYamlFile(file)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	if c.Bool("list") {
		for k, v := range settings {
			fmt.Printf("%s=%v\n", k, v)
		}
		return
	}

	key := c.Args()[0]
	var value interface{}
	if v, err := strconv.ParseInt(c.Args()[1], 10, 64); err == nil {
		value = v
	} else {
		value = c.Args()[1]
	}
	settings[key] = value
	if err = WriteKeyValueYamlFile(file, settings); err != nil {
		log.Fatalf("Failed to write settings file: %s", err)
	}
}
