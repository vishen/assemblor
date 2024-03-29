package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vishen/assemblor/bytecode"
	"github.com/vishen/assemblor/compiler"
)

func usage() {
	fmt.Printf("usage: brainfuck -f /path/to/brainfuck-file -o <output-binary>\n")
}

var (
	inputFilenameFlag  = flag.String("f", "", "path to bainfuck program to compile")
	outputFilenameFlag = flag.String("o", "", "binary executable output name. Defaults to the passed in filename")
	compileOSFlag      = flag.String("os", "macho", "Compile to OS: macho or linux")
)

func main() {
	flag.Parse()

	fileToCompile := *inputFilenameFlag
	if fileToCompile == "" {
		fmt.Printf("missing required flag -f <path/to/brainfuck program>\n")
		usage()
		return
	}

	program, err := os.ReadFile(fileToCompile)
	if err != nil {
		log.Fatalf("unable to open file %q: %v", fileToCompile, err)
	}

	// Registers used for running brainfuck
	var (
		dataReg    = bytecode.Reg8
		scratchReg = bytecode.Reg9
		zeroReg    = bytecode.Reg10
		oneReg     = bytecode.Reg11
	)

	// Create a new bytecode graph
	g := bytecode.NewGraph()

	// Reserve a 30000 int64 array
	size := uint32(64)
	gData := g.ReserveBytes(30000 * size)
	g.MovAddr(dataReg, gData) // Where we are in gData

	g.MovImm(scratchReg, 0) // Scratch register
	g.MovImm(oneReg, 1)     // Len of syswrite
	g.MovImm(zeroReg, 0)    // TODO: cmp reg because we don't have cmp imm

	loopsCounter := 0
	loopsFinished := 0

	type label struct {
		startLabel bytecode.LabelType
		endLabel   bytecode.LabelType
	}
	labels := map[int]label{}
	id := 0

	loopPos := map[int]int{}
	pos := 0

	shouldWriteToMem := false

	// Only write to memory when we need do.
	writeToMem := func() {
		if shouldWriteToMem {
			g.WriteRegToMem(dataReg, scratchReg)
			shouldWriteToMem = false
		}
	}

	for _, sym := range program {
		switch sym {
		case '+':
			g.Inc(scratchReg)
			shouldWriteToMem = true
		case '-':
			g.Dec(scratchReg)
			shouldWriteToMem = true
		case '>':
			writeToMem()
			g.AddImm(dataReg, bytecode.ImmType(size))
			g.MovMem(scratchReg, dataReg)
		case '<':
			writeToMem()
			g.SubImm(dataReg, bytecode.ImmType(size))
			g.MovMem(scratchReg, dataReg)
		case '.':
			writeToMem()
			g.SyscallWrite(dataReg, oneReg)
		case '[':
			writeToMem()
			// If the byte at the data pointer is zero, then instead of
			// moving the instruction pointer forward to the next command,
			// jump it forward to the command after the matching ] command.
			loopsCounter += 1

			endLabel := g.FutureLabel()
			g.BranchCond(scratchReg, bytecode.EQ, zeroReg, endLabel)
			startLabel := g.Label()

			id++
			labels[id] = label{startLabel, endLabel}

			pos++
			loopPos[pos] = id

		case ']':
			writeToMem()
			// If the byte at the data pointer is nonzero, then instead of
			// moving the instruction pointer forward to the next command,
			// jump it back to the command after the matching [ command.
			loopsFinished += 1

			id := loopPos[pos]
			pos--
			l := labels[id]

			g.BranchCond(scratchReg, bytecode.NEQ, zeroReg, l.startLabel)
			g.ResolveLabel(l.endLabel)
		}
	}

	// Syscall exit with status code zero
	g.MovImm(bytecode.Reg1, 0)
	g.SyscallExit(bytecode.Reg1)

	if loopsCounter != loopsFinished {
		log.Fatalf("unbalanced []: %d opened, %d closed\n", loopsCounter, loopsFinished)
	}

	bc, err := g.Bytecode()
	if err != nil {
		log.Fatalf("unable to generate bytecode: %v", err)
	}

	// Compile to machine code for the desired target os
	data, err := compiler.CompileWithOptions(
		bc,
		compiler.TargetOS(*compileOSFlag),
		compiler.X64,
	)
	if err != nil {
		log.Fatalf("unable to compile bytecode: %v", err)
	}

	// Write the machine code to an output file
	var outputFilename string
	if *outputFilenameFlag != "" {
		outputFilename = *outputFilenameFlag
	} else {
		fileBase := filepath.Base(fileToCompile)
		outputFilename = strings.Replace(fileBase, filepath.Ext(fileBase), "", -1)
	}

	if err := os.WriteFile(outputFilename, data, 0755); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote executable to %s\n", outputFilename)
}
