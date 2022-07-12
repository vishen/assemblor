package main

import (
	"flag"
	"log"
	"os"
	"runtime"

	"github.com/vishen/assemblor/bytecode"
	"github.com/vishen/assemblor/compiler"
)

var (
	outputFlag = flag.String("o", "assemblored", "filename to output executable")
	osFlag     = flag.String("os", "", "OS to compile for: macho or linux")
)

func main() {
	flag.Parse()

	g := bytecode.NewGraph()
	// Hello world program
	{
		word := "Hello World"
		d1 := g.ReserveBytes(5) // [5]byte
		g.MovAddr(bytecode.Reg1, d1)

		for _, c := range word {
			g.MovImm(bytecode.Reg2, bytecode.ImmType(c))
			g.WriteRegToMem(bytecode.Reg1, bytecode.Reg2)
			g.AddImm(bytecode.Reg1, 1)
		}

		g.MovAddr(bytecode.Reg1, d1)
		g.MovImm(bytecode.Reg2, bytecode.ImmType(len(word))) // len
		g.SyscallWrite(bytecode.Reg1, bytecode.Reg2)

		g.MovImm(bytecode.Reg10, 0)
		g.SyscallExit(bytecode.Reg10)
	}

	bc, err := g.Bytecode()
	if err != nil {
		log.Fatal(err)
	}

	OS := *osFlag
	if OS == "" {
		OS = runtime.GOOS
	}
	log.Printf("Compiling for %s_x64", OS)
	executable, err := compiler.CompileWithOptions(
		bc,
		compiler.TargetOS(OS),
		compiler.X64,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Writing executable to %s", *outputFlag)
	if err := os.WriteFile(*outputFlag, executable, 0755); err != nil {
		log.Fatal(err)
	}
}
