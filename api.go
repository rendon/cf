// In case I need to use cf funcionality from another program.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

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

func WriteKeyValueYamlFile(doc map[string]interface{}) error {
	var buf, err = yaml.Marshal(doc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(".settings.yml", buf, 0664)
}
