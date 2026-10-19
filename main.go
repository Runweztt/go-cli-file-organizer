// Command fm organizes a folder by file type and helps find what is in it.
package main

import (
	"fmt"
	"os"
)

type command struct {
	name    string
	args    string
	summary string
	run     func(args []string) error
}

var commands = []command{
	{"organize", "[-apply] [dir]", "sort files into folders by type (dry run without -apply)", runOrganize},
	{"undo", "[dir]", "reverse organize", runUndo},
	{"stats", "[dir]", "count files and bytes per type", runStats},
	{"find", "pattern [dir]", `list files matching a glob such as "*.pdf"`, runFind},
	{"dupes", "[dir]", "list files with identical content", runDupes},
}

func main() {
	if len(os.Args) > 1 {
		for _, c := range commands {
			if c.name == os.Args[1] {
				if err := c.run(os.Args[2:]); err != nil {
					fmt.Fprintln(os.Stderr, "fm:", err)
					os.Exit(1)
				}
				return
			}
		}
	}
	usage()
	os.Exit(2)
}

func usage() {
	fmt.Fprint(os.Stderr, "usage: fm <command> [arguments]\n\n")
	for _, c := range commands {
		fmt.Fprintf(os.Stderr, "  %-9s %-15s %s\n", c.name, c.args, c.summary)
	}
}
