package main

import (
	"fmt"
	"os"
)

const helpText = `typegen: Generate language-specific data definitions from sample json

Usage: 
    json-go [options]

Options:
  --file, -f      <filename>     Read sample json from named file instead of stdin
  --help, -h                     Produce help output and exit
  --language, -l  <language>     "go" or "swift", default is "go"
  --omit-empty                   Add omitempty to json tags in Go structures
  --output, -o    <filename>     Write generated text to named file instead of stdout
  --package, -p   <name>         For Go code, generate a package statent of the given name.
  --pretty                       For Swift code, align the declaration columns (off by default)
  --quiet, -q                    Suppress status messages
  --suffix, -s    <identifier>   Use this as the suffix for type names (default is "Type")
  --type,   -t    <typename>     Specify base type name in generated text
  --version, -v                  Display version number of command line tool`

func help() {
	fmt.Println(helpText)
	os.Exit(0)
}
