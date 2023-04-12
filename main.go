package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"runtime"
	"strings"

	"github.com/tucats/typegen/generate"
	"github.com/tucats/typegen/language"
	"github.com/tucats/typegen/parser"
)

// Main entrypoint for the typegen CLI. This parses the command line options, creates
// a new Parser object, and sends the input file (or stdin) to the parser. If the parse
// is successful, this calls the generator for the user-requested language type (the
// default type is Go if not specified).
func main() {
	var (
		err      error
		p        *parser.Parser
		outfile  string
		text     string
		typeName string
		pkg      string
		camel    bool
		omit     bool
		debug    bool
		aliases  bool
		pretty   bool
		quiet    bool
	)

	input := os.Stdin
	target := language.GoLang
	aliases = true

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-h", "--help":
			help()

		case "--version", "-v":
			fmt.Printf("typegen %s (%s)\n", Version, runtime.Version())
			os.Exit(0)

		case "-q", "--quiet":
			quiet = true

		case "-p", "--package":
			i++
			if i >= len(os.Args) {
				err = fmt.Errorf("missing command line argument for %s", arg)

				break
			}

			pkg = os.Args[i]

		case "-d", "--debug":
			debug = true

		case "-a", "--no-alias", "--no-aliases":
			aliases = false

		case "-s", "--suffix":
			i++
			if i >= len(os.Args) {
				err = fmt.Errorf("missing command line argument for %s", arg)

				break
			}

			parser.AliasTypeSuffix = os.Args[i]

		case "--pretty", "--pretty-print":
			pretty = true

		case "--language", "-l":
			i++
			if i >= len(os.Args) {
				err = fmt.Errorf("missing command line argument for %s", arg)

				break
			}

			name := strings.ToLower(os.Args[i])

			switch name {
			case "go", "golang":
				target = language.GoLang
				pretty = true

			case "swift":
				target = language.Swift
				pretty = false

			default:
				err = fmt.Errorf("unrecognized or unsupported language: %s", name)
			}

		case "--omit", "--omit-empty":
			omit = true

		case "-c", "--camel-case", "--camel":
			camel = true

		case "-f", "--file":
			i++
			if i >= len(os.Args) {
				err = fmt.Errorf("missing command line argument for %s", arg)

				break
			}

			name := os.Args[i]
			input, err = os.Open(name)

		case "-o", "--output":
			i++
			if i >= len(os.Args) {
				err = fmt.Errorf("missing command line argument for %s", arg)

				break
			}

			outfile = os.Args[i]

		case "-t", "--type":
			i++
			if i >= len(os.Args) {
				err = fmt.Errorf("missing command line argument for %s", arg)

				break
			}

			typeName = os.Args[i+1]

		default:
			// Assume it was the input file name and see if it's readable. If it isn't
			// a valid file name, complain that its and invalid command line option.
			input, err = os.Open(arg)
			if err != nil {
				err = fmt.Errorf("unrecognized option: %s", arg)
			}
		}

		if err != nil {
			break
		}
	}

	inputData := []byte{}

	// If no errors so far, read the file. This is done in 16k chunks until the entire file has been
	// read into memory, by appending each chunk to the inputData byte array.
	if err == nil {
		for {
			d := make([]byte, 16384)
			count := 0

			count, err = input.Read(d)
			if count == 0 {
				if err == io.EOF {
					err = nil
				}

				break
			}

			inputData = append(inputData, d[:count]...)
		}
	}

	// If reading the file finished without error, set up a new parser and configure it
	// using the options we parsed already from the command line. Then do the parse, using
	// the byte array read from input.
	if err == nil {
		p = parser.New().
			Named(typeName).
			CamelCase(camel).
			OmitEmpty(omit).
			Language(target)

		p.Debug = debug
		p.UseAliases = aliases
		p.Pretty = pretty
		p.Package = pkg

		err = p.Parse(inputData)
	}

	// If parsing completed without error, generate the output using the requested target
	// language (if not specified, Go is the default). The output is in the form of a
	// single string with internal formatting to make the source human-readable.
	//
	// The string is then either written to the stdout if there was no output file
	// designation. If tehre was an output file, it is written using the text of the
	// formatted output. The default file permissions are applied to the newly-created
	// file.
	if err == nil {
		text = generate.Generate(p, target)

		if outfile == "" {
			fmt.Println(text)
		} else {
			err = os.WriteFile(outfile, []byte(text), fs.ModePerm)
			if err == nil && !quiet {
				lines := len(strings.Split(strings.TrimSuffix(text, "\n"), "\n"))
				targetLanguage := "undefined"

				switch target {
				case language.GoLang:
					targetLanguage = "Go"
				case language.Swift:
					targetLanguage = "Swift"
				}

				fmt.Printf("Generated %d lines of %s output\n", lines, targetLanguage)
			}
		}
	}

	// If after all that, there was an error, then print out the error to the user and
	// exit with a non-zero return code. If no errors, then just exit the main program
	// cleanly
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
