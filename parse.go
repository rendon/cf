// Parse sample tests of Codeforces problem.
package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func parse(c *cli.Context) {
	if len(c.Args()) != 1 {
		log.Fatalf("USAGE: parse <URL>")
		return
	}
	var url = c.Args()[0]
	var ins, outs, err = ParseProblem(url)
	if err != nil {
		log.Fatalf("Failed to parse problem: %s", err)
	}
	for i := 0; i < len(ins); i++ {
		if err = WriteTest(ins[i], outs[i], "", i); err != nil {
			log.Fatalf("%s", err)
		}
	}

	var settings = make(map[string]interface{})
	if _, err = os.Stat(".settings.yml"); err == nil {
		if settings, err = ReadKeyValueYamlFile(".settings.yml"); err != nil {
			log.Fatalf("Failed to read settings file: %s", err)
		}
	} else {
		settings["tests"] = len(ins)
	}
	if err = WriteKeyValueYamlFile(settings); err != nil {
		log.Fatalf("Failed to write settings file: %s", err)
	}
}
