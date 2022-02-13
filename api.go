// In case I need to use cf funcionality from another program.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

// Contest contains basic fields of a Codeforces contest
type Contest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ContestListResponse a list of contests, as Codeforces returns them.
type ContestListResponse struct {
	Status string    `json:"status"`
	Result []Contest `json:"result"`
}

// Lang contains values and functions for every supported language.
type Lang struct {
	Name   string
	Sample string
	Setup  func(string) error
	Run    func(string, string, string, string) (bool, error)
}

const (
	baseURL = "http://codeforces.com"

	validatorExact      = "exact"
	validatorLines      = "lines"
	validatorNumeric1e6 = "numeric1e6"
)

// validators functions for every supported validator
var validators = map[string]func(string, string) bool{
	validatorLines: func(expected, actual string) bool {
		l1 := strings.Split(strings.TrimRight(expected, "\n"), "\n")
		l2 := strings.Split(strings.TrimRight(actual, "\n"), "\n")
		if len(l1) != len(l2) {
			return false
		}
		for i := 0; i < len(l1); i++ {
			a := strings.TrimRight(l1[i], " ")
			b := strings.TrimRight(l2[i], " ")
			if a != b {
				log.Printf("Expected %q, got %q", a, b)
				return false
			}
		}
		return true
	},
	validatorExact: func(expected, actual string) bool {
		if expected != actual {
			log.Printf("Expected %q, got %q", expected, actual)
			return false
		}
		return true
	},
	validatorNumeric1e6: func(expected, actual string) bool {
		r1 := strings.NewReader(expected)
		r2 := strings.NewReader(actual)
		for {
			var a, b float64
			_, err1 := fmt.Fscan(r1, &a)
			_, err2 := fmt.Fscan(r2, &b)
			if err1 == io.EOF || err2 == io.EOF {
				if err1 != io.EOF {
					log.Printf("Output has less tokens than answer.")
					return false
				}
				if err2 != io.EOF {
					log.Printf("Output has more tokens than answer.")
					return false
				}
				break
			}
			if math.Abs(b-a) > 1e-6 {
				log.Printf("Expected %f, got %f: delta > 1e-6\n", a, b)
				return false
			}
		}
		return true
	},
}

// cppSetup compiles C++ solution
func cppSetup(sourceFile string) error {
	ext := filepath.Ext(sourceFile)
	if ext == "" {
		return errors.New("File has no extension")
	}
	out := strings.TrimSuffix(sourceFile, ext)
	ext = ext[1:]
	cmd := exec.Command("g++", "-std=c++14", "-W", "-o", out, sourceFile)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}
		return err
	}
	return nil
}

// runBinary tests C++ solution with given input and output.
func runBinary(sourceFile, inFile, outFile, validator string) (bool, error) {
	ext := filepath.Ext(sourceFile)
	if ext == "" {
		return false, errors.New("File has no extension")
	}
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	executable := wd + "/" + strings.TrimSuffix(sourceFile, ext)
	return runTest(exec.Command(executable), inFile, outFile, validator)
}

// runRuby tests  solution with given input and output.
func runRuby(sourceFile, inFile, outFile, validator string) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	src := wd + "/" + sourceFile
	return runTest(exec.Command("ruby", src), inFile, outFile, validator)
}

func runTest(cmd *exec.Cmd, inFile, outFile, validator string) (bool, error) {
	if validators[validator] == nil {
		return false, fmt.Errorf("Unknown validator: %q", validator)
	}

	// Read input
	buf, err := ioutil.ReadFile(inFile)
	if err != nil {
		return false, err
	}
	cmd.Stdin = bytes.NewReader(buf)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Start(); err != nil {
		return false, err
	}
	if err = cmd.Wait(); err != nil {
		return false, err
	}

	buf, err = ioutil.ReadFile(outFile)
	if err != nil {
		return false, err
	}
	expected := string(buf)
	actual := out.String()
	return validators[validator](expected, actual), nil
}

// cSetup compiles C solution
func cSetup(sourceFile string) error {
	ext := filepath.Ext(sourceFile)
	if ext == "" {
		return errors.New("File has no extension")
	}
	out := strings.TrimSuffix(sourceFile, ext)
	ext = ext[1:]
	cmd := exec.Command("gcc", "-W", "-o", out, sourceFile)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}
		return err
	}
	return nil
}

// goSetup compiles Go solution
func goSetup(sourceFile string) error {
	cmd := exec.Command("go", "build", sourceFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}
		return err
	}
	return nil
}

// langs supported languages
var langs = map[string]*Lang{
	"go": &Lang{
		Name:   "Golang",
		Sample: "package main\n\nimport ()\n\nfunc main() {\n}\n",
		Setup:  goSetup,
		Run:    runBinary,
	},
	"cpp": &Lang{
		Name:   "C++",
		Sample: "#include <bits/stdc++.h>\nint main() {\n    return 0;\n}\n",
		Setup:  cppSetup,
		Run:    runBinary,
	},
	"c": &Lang{
		Name:   "C",
		Sample: "#include <stdio.h>\nint main() {\n    return 0;\n}\n",
		Setup:  cSetup,
		Run:    runBinary,
	},
	"rb": &Lang{
		Name:   "Ruby",
		Sample: "#!/usr/bin/env ruby\nputs ''\n",
		Setup:  func(sourceFile string) error { return nil },
		Run:    runRuby,
	},
}

// traverse Walks through the DOM and collect test cases.
func traverse(node *html.Node, mode, in, out *string, ins, outs *[]string) {
	var found = false
	var data = node.Data
	if node.Type == html.ElementNode && (data == "div" || data == "br") {
		if data == "div" {
			for _, a := range node.Attr {
				if a.Key == "class" && a.Val == "sample-test" {
					found = true
					break
				}
			}
		} else if data == "br" {
			if *mode == "input" {
				*in += "\n"
			} else if *mode == "output" {
				*out += "\n"
			}
		}
		if found {
			*mode = "input"
		}
	} else if node.Type == html.TextNode && *mode != "" {
		if data == "Input" {
			if *in != "" || *out != "" {
				*ins = append(*ins, *in)
				*outs = append(*outs, *out)
			}
			*in = ""
			*mode = "input"
		} else if data == "Output" {
			*out = ""
			*mode = "output"
		} else if *mode == "input" {
			*in += data
		} else if *mode == "output" {
			*out += data
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, mode, in, out, ins, outs)
	}
	if found {
		if len(*in) > 0 || len(*out) > 0 {
			*ins = append(*ins, *in)
			*outs = append(*outs, *out)
		}
	}
}

// ParseProblem Parses problem and extracts tests cases.
func ParseProblem(url string) ([]string, []string, error) {
	var client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Not found")
		},
	}
	var resp, err = client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("Error retrieving page: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, errors.New("Non-okay response")
	}

	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing document: %s", err)
	}

	var mode string
	var in string
	var out string
	var ins []string
	var outs []string
	traverse(root, &mode, &in, &out, &ins, &outs)
	return ins, outs, nil
}

// WriteTest Writes tests cases to files, one for input, and another for output.
func WriteTest(in, out, dir string, id int) error {
	var inFile = fmt.Sprintf("%s/.in%d.txt", dir, id)
	if err := ioutil.WriteFile(inFile, []byte(in), 0644); err != nil {
		return fmt.Errorf("Failed to write test input: %s", err)
	}
	var outFile = fmt.Sprintf("%s/.out%d.txt", dir, id)
	if err := ioutil.WriteFile(outFile, []byte(out), 0644); err != nil {
		return fmt.Errorf("Failed to write test output: %s", err)
	}
	return nil
}

// ReadKeyValueYamlFile reads YAML file and returns a map of string -> interface
func ReadKeyValueYamlFile(file string) (map[string]interface{}, error) {
	var buffer, err = ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var doc map[interface{}]interface{}
	err = yaml.Unmarshal(buffer, &doc)
	var kv = make(map[string]interface{})
	for k, v := range doc {
		if _, ok := k.(string); !ok {
			return nil, fmt.Errorf("%v is not a string", k)
		}
		kv[k.(string)] = v
	}
	return kv, nil
}

// WriteKeyValueYamlFile writes YAML representation of `doc` to a file.
func WriteKeyValueYamlFile(dir string, doc map[string]interface{}) error {
	var buf, err = yaml.Marshal(doc)
	if err != nil {
		return err
	}
	var file = fmt.Sprintf("%s", dir)
	return ioutil.WriteFile(file, buf, 0664)
}

// GetContestName returns the name of contest with ID `id`.
func GetContestName(id int) (string, error) {
	var resp, err = http.Get("http://codeforces.com/api/contest.list")
	if err != nil {
		return "", err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var contests ContestListResponse
	if err = json.Unmarshal(buf, &contests); err != nil {
		return "", err
	}

	for i := 0; i < len(contests.Result); i++ {
		if id == contests.Result[i].ID {
			return contests.Result[i].Name, nil
		}
	}
	return "", errors.New("Contest not found")
}

// GenerateSampleSolution generates sample solution from file name, programming
// language will be determined from file extension.
func GenerateSampleSolution(sourceFile string) error {
	// Get file extension
	var ext = filepath.Ext(sourceFile)
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	if langs[ext] == nil {
		return fmt.Errorf("Language not supported: %q", ext)
	}

	// Get directory (if any)
	var dir = filepath.Dir(sourceFile)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	code, err := getCodeFromTemplate(ext)

	// In case of error use default template
	if err != nil {
		log.Warnf("Failed to generate code from template: %q", err)
		code = langs[ext].Sample
	}

	return ioutil.WriteFile(sourceFile, []byte(code), 0664)
}
