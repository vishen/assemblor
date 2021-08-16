package main

import (
	"debug/macho"
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()

	f, err := macho.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("f=%#v\n", f)

	fmt.Println()
	fmt.Println()

	fmt.Printf("Loads: %d\n", len(f.Loads))
	for _, l := range f.Loads {
		switch t := l.(type) {
		case *macho.Segment:
			fmt.Printf("\tsegment: %#v: %d bytes\n", t.SegmentHeader, len(t.LoadBytes))
		case macho.LoadBytes:
			fmt.Printf("\tload bytes: %d (%v)\n", len(t), t)
		case *macho.Symtab:
			fmt.Printf("\tsymtab: %#v\n", t.SymtabCmd)
			for _, s := range t.Syms {
				fmt.Printf("\t\tsymbols: %#v\n", s)
			}
		case *macho.Dysymtab:
			fmt.Printf("\tdysymtab: %#v\n", t.DysymtabCmd)
		case *macho.Dylib:
			fmt.Printf("\tdylib: %#v\n", t)
		default:
			log.Fatalf("unknown type=%T", t)
		}
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("Sections: %d\n", len(f.Sections))
	for _, s := range f.Sections {
		fmt.Printf("\t%#v\n", s)
	}
}
