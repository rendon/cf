// Parse sample tests of Codeforces problem.
package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"

	"github.com/codegangsta/cli"
)

func traverse(n *html.Node, mode, in, out *string, ins, outs *[]string) {
	var found = false
	if n.Type == html.ElementNode && (n.Data == "div" || n.Data == "br") {
		if n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "sample-test" {
					found = true
					break
				}
			}
		} else if n.Data == "br" {
			if *mode == "input" {
				*in += "\n"
			} else if *mode == "output" {
				*out += "\n"
			}
		}
		if found {
			*mode = "input"
		}
	} else if n.Type == html.TextNode && *mode != "" {
		if n.Data == "Input" {
			if *in != "" || *out != "" {
				*ins = append(*ins, *in)
				*outs = append(*outs, *out)
			}
			*in = ""
			*mode = "input"
		} else if n.Data == "Output" {
			*out = ""
			*mode = "output"
		} else if *mode == "input" {
			*in += n.Data
		} else if *mode == "output" {
			*out += n.Data
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
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
		log.Printf("USAGE: parse <URL>")
		return
	}
	var url = c.Args()[0]

	var resp, err = http.Get(url)
	if err != nil {
		log.Printf("Error retrieving page: %s", err)
		return
	}

	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("Error parsing document: %s", err)
		return
	}
	var mode string
	var in string
	var out string
	var ins []string
	var outs []string
	traverse(root, &mode, &in, &out, &ins, &outs)

	for i := 0; i < len(ins); i++ {
		fmt.Printf("Input\n")
		fmt.Printf("%s\n", ins[i])
		fmt.Printf("Output\n")
		fmt.Printf("%s\n", outs[i])
		fmt.Printf("\n")
	}
}
