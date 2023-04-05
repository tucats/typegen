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
  --no-aliases                   Do not attempt to generate type names based on object keys
  --omit-empty                   Add omitempty to json tags in Go structures
  --output, -o    <filename>     Write generated text to named file instead of stdout
  --pretty                       For Swift code, align the declaration columns (off by default)
  --suffix, -s    <identifier>   Use this as the suffix for type names (default is "Type")
  --type,   -t    <typename>     Specify base type name in generated text
  --version, -v                  Display version number of command line tool`

func help() {
	fmt.Println(helpText)
	os.Exit(0)
}
