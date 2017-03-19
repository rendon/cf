package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func test(c *cli.Context) {
	settings, err := ReadKeyValueYamlFile(".settings.yml")
	if err != nil {
		log.Fatalf("No '.settings.yml' file found.")
	}
	lang, ok := settings["lang"].(string)
	if !ok {
		log.Fatalf("No 'lang' field found in settings file.")
	}
	srcFile, ok := settings["src_file"].(string)
	if !ok {
		log.Fatalf("No 'src_file' field found in settings file.")
	}
	tests, ok := settings["tests"].(int)
	if !ok {
		log.Fatalf("No 'tests' field found in settings file.")
	}
	validator, ok := settings["validator"].(string)
	if !ok {
		validator = validatorLines
	}

	if langs[lang] == nil {
		log.Fatalf("Language %q not supported.", lang)
	}

	if err = langs[lang].Setup(srcFile); err != nil {
		log.Fatalf("Failed to setup source file: %s", err)
	}
	for i := 1; i <= tests; i++ {
		in := fmt.Sprintf(".in%d.txt", i)
		out := fmt.Sprintf(".out%d.txt", i)
		fmt.Printf("Test #%d:", i)
		if c.GlobalBool("verbose") {
			fmt.Println()
		}

		passed, err := langs[lang].Run(srcFile, in, out, validator)
		if err != nil {
			log.Fatalf("Test %d failed: %s", i, err)
		}
		if passed {
			fmt.Printf("PASSED\n")
		} else {
			fmt.Printf("FAILED\n")
		}
	}
}
