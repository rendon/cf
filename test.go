package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
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
	sourceFile, ok := settings["sourceFile"].(string)
	if !ok {
		log.Fatalf("No 'sourceFile' field found in settings file.")
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

	if err = langs[lang].Setup(sourceFile); err != nil {
		log.Fatalf("Failed to setup source file: %s", err)
	}

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	failureCount := 0
	for i := 1; i <= tests; i++ {
		in := fmt.Sprintf(".in%d.txt", i)
		out := fmt.Sprintf(".out%d.txt", i)
		fmt.Printf("Test #%d: ", i)
		if c.GlobalBool("verbose") {
			fmt.Println()
		}

		start := time.Now()
		passed, err := langs[lang].Run(sourceFile, in, out, validator)
		if err != nil {
			log.Fatalf("Test %d failed: %s", i, err)
		}
		end := time.Now()
		if passed {
			fmt.Printf("%s %.3fs\n", green("PASSED"), end.Sub(start).Seconds())
		} else {
			fmt.Printf("%s %.3fs\n", red("FAILED"), end.Sub(start).Seconds())
			failureCount++
		}
	}

	fmt.Println()
	if failureCount == 0 {
		fmt.Printf("%s\n", green("Well done!"))
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", red("Try again!"))
		os.Exit(1)
	}
}
