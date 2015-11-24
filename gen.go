package main

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func getCodeFromTemplate(ext string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	settings, err := ReadKeyValueYamlFile(user.HomeDir + "/.cf.yml")
	if err != nil {
		return "", err
	}

	var key = "template." + ext
	template, ok := settings[key].(string)
	if !ok {
		return "", fmt.Errorf("Failed to get %q from setting file.", key)
	}
	buf, err := ioutil.ReadFile(template)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func gen(c *cli.Context) {
	if len(c.Args()) != 1 {
		cli.ShowCommandHelp(c, "gen")
		return
	}
	var srcFile = c.Args()[0]
	if err := GenerateSampleSolution(srcFile); err != nil {
		log.Fatalf("Failed to generate sample solution: %s", err)
	}

	// At this point we know srcFile contains a valid extension
	ext := filepath.Ext(srcFile)[1:]
	settings, err := ReadKeyValueYamlFile(".settings.yml")
	if err != nil {
		log.Printf("Failed to read settings file: %s\n", err)
		settings = make(map[string]interface{})
	}
	settings["lang"] = ext
	settings["src_file"] = srcFile
	if err = WriteKeyValueYamlFile(".settings.yml", settings); err != nil {
		log.Printf("Failed to write settings file: %s\n", err)
	}

}
