package generate

import (
	"fmt"
	"sort"

	"github.com/tucats/typegen/language"
	"github.com/tucats/typegen/parser"
)

// Generate produces a text representation of the declaration, given
// a target language such as language.GoLang or language.Swift.
func Generate(p *parser.Parser, target language.Language) string {
	if p.Debug {
		dumpTree(p)
	}

	switch target {
	case language.GoLang:
		return generateGo(p)

	case language.Swift:
		return generateSwift(p)

	default:
		return fmt.Sprintf("Unsupported language: %v", target)
	}
}

func dumpTree(p *parser.Parser) {
	fmt.Printf("\n------------\n")
	fmt.Printf("Types:\n")

	keys := []string{}

	for key := range p.Types {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("    type: %s\n", key)
		dumpElement(p.Types[key], 2)
	}

	fmt.Printf("Base type:\n")
	dumpElement(p.Type, 0)

	fmt.Printf("------------\n\n")
}

func dumpElement(t *parser.Type, depth int) {
	switch t.Kind {
	case parser.BoolType:
		fmt.Printf("%s bool\n", pad(" ", depth*4))

	case parser.IntType:
		fmt.Printf("%s int\n", pad(" ", depth*4))

	case parser.InterfaceType:
		fmt.Printf("%s interface\n", pad(" ", depth*4))

	case parser.FloatType:
		fmt.Printf("%s float\n", pad(" ", depth*4))

	case parser.StringType:
		fmt.Printf("%s string\n", pad(" ", depth*4))

	case parser.ArrayType:
		fmt.Printf("%s array of %s\n", pad(" ", depth*4), t.BaseType)

	case parser.GenericArrayType:
		fmt.Printf("%s generic array\n", pad(" ", depth*4))

	case parser.TypeType:
		fmt.Printf("%s generated: %s\n", pad(" ", depth*4), t.Name)

	case parser.GenericStructType:
		fmt.Printf("%s generic struct\n", pad(" ", depth*4))

	case parser.StructType:
		fmt.Printf("%s struct\n", pad(" ", depth*4))

		for _, field := range t.Fields {
			fmt.Printf("%s field: %s\n", pad(" ", (depth+1)*4), field.Name)
			fmt.Printf("%s type: %s\n", pad(" ", (depth+2)*4), field.Type)
		}

	default:
		fmt.Printf("Unknown type %#v\n", t)
	}
}
