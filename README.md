# typegen

Generate Go or Swift type definitions based on sample JSON data

## Introduction

Given a sample JSON file, this tool will generate language type definitions
for the types and structures needed to be able to represent the JSON in memory.
The type defintions include JSON tags on each item to support controlling how
the types are marshalled or unmarshalled.

The tool will attempt to generate specific type definitions where possible. If
an array is not heterogeneous, for example, it will be defined as []interface{}
when the language specified is `go`, but if all the elements of the array are 
the same, then it is an array of the given type.

Nesting is detected, such that the generated type may reference additional
generated types to represents objects defined within the JSON object.

## Command

By default, the program reads input from stdin that is the sample JSON file, and
sends the generated Go program code to stdout. These can be overridden by command
line options.

```text
Usage:

     typegen  [options]

Options:
     --camel, -c                      Force use of camel case in field names
     --file, -f       <filename>      Read from the named file instead of stdin
     --language, -l   <language>      Go or Swift; if not specified the default is Go
     --omit-empty                     For Go, add "omitempty" to the json tags
     --output, -o     <filename>      Write to the named file isntead of stdout
     --type, -t       <name>          Name of the base type to be created
     --help, -h                       Produce help output
```
