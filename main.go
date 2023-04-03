package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/tucats/typegen/generate"
	"github.com/tucats/typegen/parser"
)

func main() {
	var (
		err      error
		p        *parser.Parser
		outfile  string
		text     string
		typeName string
		camel    bool
		omit     bool
	)

	input := os.Stdin
	language := generate.GoLang

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-h", "--help":
			help()

		case "--language", "-l":
			i++
			name := strings.ToLower(os.Args[i])

			switch name {
			case "go", "golang":
				language = generate.GoLang

			case "swift":
				language = generate.Swift

			default:
				err = fmt.Errorf("unrecognized or unsupported language: %s", name)
			}

		case "--omit", "--omit-empty":
			omit = true

		case "-c", "--camel-case", "--camel":
			camel = true

		case "-f", "--file":
			name := os.Args[i+1]
			input, err = os.Open(name)
			i++

		case "-o", "--output":
			outfile = os.Args[i+1]
			i++

		case "-t", "--type":
			typeName = os.Args[i+1]
			i++

		default:
			err = fmt.Errorf("unrecognized option: %s", arg)
		}

		if err != nil {
			break
		}
	}

	b := []byte{}

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

			b = append(b, d[:count]...)
		}
	}

	if err == nil {
		p = parser.New().Named(typeName).CamelCase(camel).OmitEmpty(omit)

		err = p.Parse(b)
	}

	if err == nil {
		text = generate.Generate(p, language)
	}

	if err == nil {
		if outfile == "" {
			fmt.Println(text)
		} else {
			err = os.WriteFile(outfile, []byte(text), fs.ModePerm)
		}
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
