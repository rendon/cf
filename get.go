// Parse sample tests of Codeforces problem.
package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func get(c *cli.Context) {
	var problemURL string
	if len(c.Args()) == 1 {
		problemURL = c.Args()[0]
	} else {
		url, err := getUrlFromSettings()
		if err != nil {
			log.Fatalf("Unable to get URL from the settings. Please provide a URL")
		}
		problemURL = url
	}

	var ins, outs, err = ParseProblem(problemURL)
	if err != nil {
		log.Fatalf("Failed to parse problem: %s", err)
	}
	// TODO: Add test
	for i := 1; i <= len(ins); i++ {
		if err = WriteTest(ins[i-1], outs[i-1], ".", i); err != nil {
			log.Fatalf("%s", err)
		}
	}

	settings, err := ReadKeyValueYamlFile(".settings.yml")
	if err != nil {
		log.Printf("Unable to read settings file: %s", err)
		settings = make(map[string]interface{})
	}

	// The problem description determines these values. Override the settings if needed.
	settings["problemUrl"] = problemURL
	settings["tests"] = len(ins)

	if err = WriteKeyValueYamlFile(".settings.yml", settings); err != nil {
		log.Fatalf("Unable to write settings file: %s", err)
	}
}
