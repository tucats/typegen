# typegen

Generate Go or Swift type definitions based on sample JSON data

## Introduction

The typegen tool is a command-line interface, which given a sample JSON file,
will generate language type definitions
for the types and structures needed to be able to represent the JSON in memory.
The type defintions include JSON tags on each item to support controlling how
the types are marshalled or unmarshalled.

The tool will attempt to generate specific type definitions where possible. If
an array is not heterogeneous, for example, it will be defined as []interface{}
when the language specified is `go`, but if all the elements of the array are
the same, then it is an array of the given type.

Nesting is detected, such that the generated type may reference additional
generated types to represents objects defined within the JSON object.

If the command generates names that conflict with object member names, you can
use the `--suffix` option to specify the string suffix attached to type names.
If not specified, the default is to append "Type" to the type names.

## Command

By default, the program reads input from stdin that is the sample JSON file, and
sends the generated Go program code to stdout. These can be overridden by command
line options.

```text
Usage:

     typegen  [options]

Options:
     --file, -f       <filename>      Read from the named file instead of stdin
     --help, -h                       Produce help output and exit
     --language, -l   <language>      Go or Swift; if not specified the default is Go
     --omit-empty                     For Go, add "omitempty" to the json tags
     --output, -o     <filename>      Write to the named file instead of stdout 
     --package, -p    <name>          For Go code, generate a package statement with the given name
     --pretty                         For Swift code, align the declaration columns (off by default)
     --quiet, -q                      Suppress status messages to stdout
     --type, -t       <name>          Name of the base type to be created
     --version, -v                    Display the typegen version number and exit
```

## Example using Go

Consider the following JSON data file, which we will assume is named "data.json":

```json
{
    "team": {
        "name": "All Stars",
        "members": [{
            "name": "Dick",
            "age": 55,
            "gender": "m"
        }, {
            "name": "Jane",
            "age": 52,
            "gender": "f"
        }]
    }
}
```

To generate the associated Go code, you can use a command line like this:

```sh
testgen --file data.json --language go 
```

The program will run, and generate it's output to stdout. (You can redirect it to a file
using the `--output` command line option). The use of `--language` is optional in this case
since Go is the default language.

```go
type MembersType struct {
  Age      int     `json:"age"`
  Gender   string  `json:"gender"`
  Name     string  `json:"name"`
}

type TeamType struct {
  Members   []MembersType  `json:"members"`
  Name      string         `json:"name"`
}

type jsonData struct {
  Team   TeamType  `json:"team"`
}
```

## Example using Swift

Consider the following JSON data file, which we will assume is named "data.json":

```json
{
    "glossary": {
        "title": "example glossary",
        "GlossDiv": {
            "title": "S",
            "GlossList": {
                "GlossEntry": {
                    "ID": "SGML",
                    "SortAs": "SGML",
                    "GlossTerm": "Standard Generalized Markup Language",
                    "Acronym": "SGML",
                    "Abbrev": "ISO 8879:1986",
                    "GlossDef": {
                        "para": "A meta-markup language, used to create markup languages such as DocBook.",
                        "GlossSeeAlso": ["GML", "XML"]
                    },
                    "GlossSee": "markup"
                }
            }
        }
    }
}
```

To generate the associated Swift code, you can use a command line like this:

```sh
testgen --file data.json --language swift 
```

The program will run, and generate it's output to stdout. (You can redirect it to a file
using the `--output` command line option). The use of `--language` indicates that the output
is generated using Swift syntax.

```swift
class GlossDefType: Codable {
  var GlossSeeAlso:     [String]()
  var para:             String    
}

class GlossEntryType: Codable {
  var Abbrev:        String      
  var Acronym:       String      
  var GlossDef:      GlossDefType
  var GlossSee:      String      
  var GlossTerm:     String      
  var ID:            String      
  var SortAs:        String      
}

class GlossListType: Codable {
  var GlossEntry:     GlossEntryType
}

class GlossDivType: Codable {
  var GlossList:     GlossListType
  var title:         String       
}

class GlossaryType: Codable {
  var GlossDiv:     GlossDivType
  var title:        String      
}

class JsonData: Codable {
  var glossary:     GlossaryType
}

```
