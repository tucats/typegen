package main

import (
	"fmt"
	"os"
)

const helpText = `json-go: Generate Go data definitions from sample json

Usage: 
    json-go [options]

Options:
  --camel,  -c                   Use camel-case for struct member names
  --file,   -f    <filename>     Read sample json from named file instead of stdin
  --language, -l  <language>     "go" or "swift", default is "go"
  --omit-empty                   Add omitempty to json tags in Go structures
  --output, -o    <filename>     Write generated text to named file instead of stdout
  --type,   -t    <typename>     Specify base type name in generated text`

func help() {
	fmt.Println(helpText)
	os.Exit(0)
}
