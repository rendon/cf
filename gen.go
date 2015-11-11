package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"

	"github.com/codegangsta/cli"
)

var langSamples = map[string]string{
	".go":  "package main\n\nimport ()\n\nfunc main() {\n}\n",
	".cpp": "#include <bits/stdc++.h>\nint main() {\n    return 0;\n}\n",
}

func getCodeFromTemplate(ext string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	settings, err := ReadKeyValueYamlFile(user.HomeDir + "/.cf.yml")
	if err != nil {
		return "", err
	}

	var key = "template" + ext
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
	var srcfile = c.Args()[0]

	// Get file extension
	var re = regexp.MustCompile(`\.([^.]+)$`)
	var ext = re.FindString(srcfile)
	if langSamples[ext] == "" {
		log.Fatalf("Language not supported: %q", ext)
	}

	// Get directory (if any)
	re = regexp.MustCompile(`([^/]*/)*`)
	var dir = re.FindString(srcfile)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %s", err)
		}
	}

	code, err := getCodeFromTemplate(ext)
	if err != nil {
		//log.Printf("No template file found, using minimal template.")
		code = langSamples[ext]
	}

	if err := ioutil.WriteFile(srcfile, []byte(code), 0664); err != nil {
		log.Fatalf("Failed to write file contents: %s", err)
	}
}
