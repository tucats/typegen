package main

import (
	"fmt"
	"os"
)

const helpText = `typegen: Generate language-specific data definitions from sample json

Usage: 
    json-go [options]

Options:
  --aliases, -a                  Attempt to generate type names based on object keys
  --file,   -f    <filename>     Read sample json from named file instead of stdin
  --language, -l  <language>     "go" or "swift", default is "go"
  --omit-empty                   Add omitempty to json tags in Go structures
  --output, -o    <filename>     Write generated text to named file instead of stdout
  --suffix, -s    <identifier>   Use this as the suffix for type names (default is "Type")
  --type,   -t    <typename>     Specify base type name in generated text
  --version, -v                  Display version number of command line tool`

func help() {
	fmt.Println(helpText)
	os.Exit(0)
}
