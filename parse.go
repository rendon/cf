// Parse sample tests of Codeforces problem.
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/cli"
)

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

func parse(c *cli.Context) {
	if len(c.Args()) != 1 {
		log.Fatalf("USAGE: parse <URL>")
		return
	}
	var url = c.Args()[0]

	var resp, err = http.Get(url)
	if err != nil {
		log.Fatalf("Error retrieving page: %s", err)
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing document: %s", err)
	}
	var mode string
	var in string
	var out string
	var ins []string
	var outs []string
	traverse(root, &mode, &in, &out, &ins, &outs)

	for i := 0; i < len(ins); i++ {
		var inFile = fmt.Sprintf(".in_%d.txt", i)
		if err = ioutil.WriteFile(inFile, []byte(ins[i]), 0644); err != nil {
			log.Fatalf("Failed to write input for test case %d: %s", i, err)
		}
		var outFile = fmt.Sprintf(".out_%d.txt", i)
		if err = ioutil.WriteFile(outFile, []byte(outs[i]), 0644); err != nil {
			log.Fatalf("Failed to write output for test case %d: %s", i, err)
		}
	}
}
