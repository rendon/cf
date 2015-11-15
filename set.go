package main

import (
	"os/user"
	"strconv"

	"github.com/codegangsta/cli"
)

func set(c *cli.Context) {
	if len(c.Args()) != 2 {
		cli.ShowCommandHelp(c, "set")
		return
	}
	key := c.Args()[0]
	var value interface{}
	if v, err := strconv.ParseInt(c.Args()[1], 10, 64); err == nil {
		value = v
	} else {
		value = c.Args()[1]
	}
	if c.Bool("global") {
		user, err := user.Current()
		if err != nil {
			log.Fatalf("Failed to get user's home directory: %s", err)
		}
		file := user.HomeDir + "/.cf.yml"
		settings, err := ReadKeyValueYamlFile(file)
		if err != nil {
			log.Fatalf("Failed to read settings file: %s", err)
		}
		settings[key] = value
		if err = WriteKeyValueYamlFile(file, settings); err != nil {
			log.Fatalf("Failed to write settings file: %s", err)
		}
	} else {
		file := ".settings.yml"
		settings, err := ReadKeyValueYamlFile(file)
		if err != nil {
			log.Fatalf("Failed to read settings file: %s", err)
		}
		settings[key] = value
		if err = WriteKeyValueYamlFile(file, settings); err != nil {
			log.Fatalf("Failed to write settings file: %s", err)
		}
	}
}
