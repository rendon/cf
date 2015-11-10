package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

func setup(c *cli.Context) {
	if len(c.Args()) != 1 {
		cli.ShowCommandHelp(c, "setup")
		return
	}

	var cid, err = strconv.Atoi(c.Args()[0])
	if err != nil {
		log.Fatalf("%q is not a valid contest ID: %s", c.Args()[0], err)
	}
	for p := 'A'; p <= 'Z'; p++ {
		var url = fmt.Sprintf("http://codeforces.com/contest/%d/problem/%c", cid, p)
		var ins, outs, err = ParseProblem(url)
		if err != nil {
			break
		}
		var dir = fmt.Sprintf("%c/", p)
		if err = os.MkdirAll(dir, 0775); err != nil {
			log.Fatalf("Failed to create directory: %s", err)
		}
		for i := 0; i < len(ins); i++ {
			if err = WriteTest(ins[i], outs[i], dir, i); err != nil {
				log.Fatalf("Problem %c: %s", p, err)
			}
		}
		var settings = map[string]interface{}{"tests": len(ins)}
		if err = WriteKeyValueYamlFile(dir, settings); err != nil {
			log.Fatalf("%c: Failed to write settings file: %s", p, err)
		}
		fmt.Printf("Problem %c is ready!\n", p)
	}
}
