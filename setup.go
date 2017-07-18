package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
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

	var lang = c.String("lang")
	if lang != "" {
		if langs[lang] == nil {
			log.Fatalf("No language with code %q was found.", lang)
		}
	}

	cname, err := GetContestName(cid)
	if err != nil {
		log.Printf("Failed to get contest name, using ID instead.")
		cname = fmt.Sprintf("%d", cid)
	}

	// Sanitize contest name to avoid problems with directory names
	cname = strings.Replace(cname, " ", "", -1)
	cname = strings.Replace(cname, "#", "", -1)
	cname = strings.Replace(cname, "(", "", -1)
	cname = strings.Replace(cname, ")", "", -1)

	if err = os.MkdirAll(cname, 0755); err != nil {
		log.Fatalf("Failed to create contest directory: %s", err)
	}
	for p := 'A'; p <= 'Z'; p++ {
		var url = fmt.Sprintf("%s/contest/%d/problem/%c", baseURL, cid, p)
		var ins, outs, err = ParseProblem(url)
		if err != nil {
			break
		}
		var dir = fmt.Sprintf("%s/%c", cname, p)
		if err = os.MkdirAll(dir, 0775); err != nil {
			log.Fatalf("Failed to create directory: %s", err)
		}
		for i := 1; i <= len(ins); i++ {
			if err = WriteTest(ins[i-1], outs[i-1], dir, i); err != nil {
				log.Fatalf("Problem %c: %s", p, err)
			}
		}

		var settings = map[string]interface{}{"tests": len(ins)}
		if lang != "" {
			var srcFile = fmt.Sprintf("%s/%c.%s", dir, p, lang)
			if err = GenerateSampleSolution(srcFile); err != nil {
				log.Printf("Failed to generate sample solution: %s", err)
			}
			settings["lang"] = lang
			settings["src_file"] = filepath.Base(srcFile)
		}

		err = WriteKeyValueYamlFile(dir+"/.settings.yml", settings)
		if err != nil {
			log.Fatalf("%c: Failed to write settings file: %s", p, err)
		}
		fmt.Printf("Problem %c is ready!\n", p)
	}
}
