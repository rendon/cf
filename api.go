// In case I need to use cf funcionality from another program.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

type Contest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ContestListResponse struct {
	Status string    `json:"status"`
	Result []Contest `json:"result"`
}

type Lang struct {
	Name   string
	Sample string
	Setup  func(string) error
	Run    func(string, string, string) (bool, error)
}

const (
	baseURL = "http://codeforces.com"
)

func cppSetup(srcFile string) error {
	ext := filepath.Ext(srcFile)
	if ext == "" {
		return errors.New("File has no extension")
	}
	out := strings.TrimSuffix(srcFile, ext)
	ext = ext[1:]
	cmd := exec.Command("g++", "-W", "-o", out, srcFile)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func cppRun(srcFile, inFile, outFile string) (bool, error) {
	ext := filepath.Ext(srcFile)
	if ext == "" {
		return false, errors.New("File has no extension")
	}
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	executable := wd + "/" + strings.TrimSuffix(srcFile, ext)
	cmd := exec.Command(executable)

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
	return expected == actual, nil
}

var langs = map[string]*Lang{
	"go": &Lang{
		Name:   "Golang",
		Sample: "package main\n\nimport ()\n\nfunc main() {\n}\n",
	},
	"cpp": &Lang{
		Name:   "C++",
		Sample: "#include <bits/stdc++.h>\nint main() {\n    return 0;\n}\n",
		Setup:  cppSetup,
		Run:    cppRun,
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
	var inFile = fmt.Sprintf("%s.in_%d.txt", dir, id)
	if err := ioutil.WriteFile(inFile, []byte(in), 0644); err != nil {
		return fmt.Errorf("Failed to write test input: %s", err)
	}
	var outFile = fmt.Sprintf("%s.out_%d.txt", dir, id)
	if err := ioutil.WriteFile(outFile, []byte(out), 0644); err != nil {
		return fmt.Errorf("Failed to write test output: %s", err)
	}
	return nil
}

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
			return nil, fmt.Errorf("%v is not a string")
		}
		kv[k.(string)] = v
	}
	return kv, nil
}

func WriteKeyValueYamlFile(dir string, doc map[string]interface{}) error {
	var buf, err = yaml.Marshal(doc)
	if err != nil {
		return err
	}
	var file = fmt.Sprintf("%s.settings.yml", dir)
	return ioutil.WriteFile(file, buf, 0664)
}

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

func GenerateSampleSolution(srcFile string) error {
	// Get file extension
	var ext = filepath.Ext(srcFile)
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	if langs[ext] == nil {
		return fmt.Errorf("Language not supported: %q", ext)
	}

	// Get directory (if any)
	var dir = filepath.Dir(srcFile)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	code, err := getCodeFromTemplate(ext)
	if err != nil {
		code = langs[ext].Sample
	}

	if err := ioutil.WriteFile(srcFile, []byte(code), 0664); err != nil {
		return err
	}
	return nil
}
