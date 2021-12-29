package main

import (
	"flag"
	"log"
	"os"

	"github.com/vishen/assemblor/bytecode"
	"github.com/vishen/assemblor/compiler"
)

var (
	outputFlag = flag.String("o", "assemblored", "filename to output executable")
	osFlag     = flag.String("os", "macho", "OS to compile for: macho or linux")
)

func testGraph() *bytecode.Graph {
	g := bytecode.NewGraph()

	/*
		g.MovImm(bytecode.Reg1, 0x05)
		g.MovAddr(bytecode.Reg2, bytecode.AddrType(0xdeadbe))
		g.SyscallExit(123)
	*/

	/*
		l1 := g.FutureLabel()
		g.Jmp(l1)

		l2 := g.Label()
		g.AddImm(bytecode.Reg1, 0x03)
		g.SyscallExit(125)

		g.ResolveLabel(l1)
		g.AddImm(bytecode.Reg1, 0x02)
		g.Jmp(l2)
	*/

	// g.WriteReg(addr, bytecode.Reg1)

	/*
		addr := g.ReserveBytes(100 * 8)  // [100]int8
		g.WriteImm(addr, 0x48)           // H
		g.WriteImm(addr.Offset(1), 0x65) // e
		g.WriteImm(addr.Offset(2), 0x6c) // l
		g.WriteImm(addr.Offset(3), 0x6c) // l
		g.WriteImm(addr.Offset(4), 0x6f) // o
		g.MovAddr(bytecode.Reg1, addr)
		g.MovImm(bytecode.Reg2, 0x05)
		g.SyscallWrite(bytecode.Reg1, bytecode.Reg2)

		g.MovImm(bytecode.Reg1, 125)
		g.SyscallExit(bytecode.Reg1)

		g.CmpReg(bytecode.Reg1, bytecode.Reg2)
		g.CmpImm(bytecode.Reg1, 0x03)
		g.CmpImm(bytecode.Reg1, 0xdeadbee)
		g.CmpImm(bytecode.Reg10, 0x03)
		g.CmpImm(bytecode.Reg10, 0xdeadbee)
	*/

	/*
		g.MovImm(bytecode.Reg1, 0)
		g.MovImm(bytecode.Reg2, 5)
		l1 := g.Label()
		l2 := g.FutureLabel()
		g.BranchCond(bytecode.Reg1, bytecode.EQ, bytecode.Reg2, l2)
		g.Inc(bytecode.Reg1)
		g.Jmp(l1)

		g.ResolveLabel(l2)
		g.MovReg(bytecode.Reg10, bytecode.Reg1)

		g.AddImm(bytecode.Reg1, 48)     // Turn into ascii number
		addr := g.ReserveBytes(1 * 32)  // Reserve data: [1]int32
		g.WriteReg(addr, bytecode.Reg1) // Move ascii code in reg1 to addr

		g.MovAddr(bytecode.Reg1, addr) // Move the addr to reg1

		g.MovImm(bytecode.Reg2, 1) // Length to print
		g.SyscallWrite(bytecode.Reg1, bytecode.Reg2)

		g.SyscallExit(bytecode.Reg10) // Should be what started in reg2
	*/

	d1 := g.ReserveBytes(5) // [5]byte
	g.MovAddr(bytecode.Reg1, d1)

	for _, c := range "Hello" {
		g.MovImm(bytecode.Reg2, bytecode.ImmType(c))
		g.WriteRegToMem(bytecode.Reg1, bytecode.Reg2)
		g.AddImm(bytecode.Reg1, 1)
	}

	g.PushImm(0x1234)
	g.PushReg(bytecode.Reg12)
	g.PopReg(bytecode.Reg9)

	g.MovAddr(bytecode.Reg1, d1)
	g.MovImm(bytecode.Reg2, 5) // len
	g.SyscallWrite(bytecode.Reg1, bytecode.Reg2)

	g.MovImm(bytecode.Reg10, 0)
	g.SyscallExit(bytecode.Reg10)

	return g
}

func main() {
	flag.Parse()

	g := testGraph()
	// g.Print()

	bc, err := g.Bytecode()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Compiling for %s_x64", *osFlag)
	executable, err := compiler.CompileWithOptions(
		bc,
		compiler.TargetOS(*osFlag),
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
