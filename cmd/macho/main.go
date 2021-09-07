package main

import (
	"bytes"
	"debug/macho"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/arch/x86/x86asm"
)

func main() {
	flag.Parse()

	data, err := os.ReadFile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewReader(data)
	f, err := macho.NewFile(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("f=%#v\n", f)

	fmt.Println()
	fmt.Println()

	textStart := 0
	textAddr := uint64(0)
	fmt.Printf("Loads: %d\n", len(f.Loads))
	for _, l := range f.Loads {
		switch t := l.(type) {
		case *macho.Segment:
			fmt.Printf("\tsegment: %#v: %d bytes\n", t.SegmentHeader, len(t.LoadBytes))
			textStart += len(t.LoadBytes)
			if t.Name == "__TEXT" {
				textAddr = t.Addr
			}
		case macho.LoadBytes:
			fmt.Printf("\tload bytes: %d (%v)\n", len(t), t)
			textStart += len(t)
			textAddr += uint64(len(t))
		case *macho.Symtab:
			fmt.Printf("\tsymtab: %#v\n", t.SymtabCmd)
			for _, s := range t.Syms {
				fmt.Printf("\t\tsymbols: %#v\n", s)
			}
			textStart += len(t.LoadBytes)
		case *macho.Dysymtab:
			fmt.Printf("\tdysymtab: %#v\n", t.DysymtabCmd)
			textStart += len(t.LoadBytes)
		case *macho.Dylib:
			fmt.Printf("\tdylib: %#v\n", t)
			textStart += len(t.LoadBytes)
		default:
			log.Fatalf("unknown type=%T", t)
		}
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("Sections: %d\n", len(f.Sections))
	for _, s := range f.Sections {
		fmt.Printf("\t%#v\n", s)
		textStart += int(s.Size)
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("Assembly: start addr 0x%x\n", textAddr)
	nop := 0
	for i := textStart; i < len(data); {
		if data[i] == 0x00 {
			i++
			nop++
			continue
		}
		if nop > 0 {
			fmt.Printf("nop: %d\n", nop)
			nop = 0
		}
		inst, err := x86asm.Decode(data[i:], 64)
		if err != nil {
			fmt.Print(data[i])
			i++
			continue
		}
		codeStart := 30
		for j := 0; j < inst.Len; j++ {
			fmt.Printf("%x ", data[j+i])
			if data[j+i] >= 0x10 {
				codeStart -= 3
			} else {
				codeStart -= 2
			}
		}
		i += inst.Len
		for i := 0; i < codeStart; i++ {
			fmt.Printf(" ")
		}
		fmt.Println(inst)
	}
	if nop > 0 {
		fmt.Printf("nop: %d\n", nop)
		nop = 0
	}
}
