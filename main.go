package main

import (
	"flag"
	"log"
	"os"

	"github.com/vishen/assemblor/bytecode"
	"github.com/vishen/assemblor/ld"
	"github.com/vishen/assemblor/x64"
)

var (
	outputFlag = flag.String("o", "assemblored", "filename to output executable")
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

	return g
}

func main() {
	g := testGraph()
	g.Print()

	bc, err := g.Bytecode()
	if err != nil {
		log.Fatal(err)
	}

	linker := ld.NewMacho()
	code, bssSize := x64.Compile(x64.Macho, bc, linker.BssAddr())
	executable := linker.Link(code, bssSize)

	log.Printf("Writing executable to %s", *outputFlag)
	if err := os.WriteFile(*outputFlag, executable, 0755); err != nil {
		log.Fatal(err)
	}
}

/*
func example1() {
	g := bytecode.NewGraph()
	l1 := g.Label()
	g.MovImm(bytecode.Reg1, 0x05)
	g.MovReg(bytecode.Reg2, bytecode.Reg1)
	g.Jmp(l1)
	g.AddReg(bytecode.Reg3, bytecode.Reg1)
	g.AddImm(bytecode.Reg4, 0x100)
	g.AddImm(bytecode.Reg4, 0xdeadbee)
	g.Dec(bytecode.Reg4)
	g.Inc(bytecode.Reg4)
	g.MovImm(bytecode.Reg1, 0xdeadbeef)
	g.MovImm(bytecode.Reg1, 0xdeadbee)

	g.MovImm(bytecode.Reg1, 0x00)
	//l1 := g.Label()
	g.Inc(bytecode.Reg1)
	// g.CondImm(bytecode.Reg1, bytecode.EQ, 5, exit)
	exit := g.FutureLabel()
	g.Jmp(exit)
	g.MovImm(bytecode.Reg1, 0xdeadbeef)
	g.MovImm(bytecode.Reg1, 0xdeadbeef)
	g.ResolveLabel(exit)
	g.SyscallExit(125)

	g.Print()

	bc, err := g.Bytecode()
	if err != nil {
		log.Fatal(err)
	}
	code := x64.Compile(x64.Macho, bc)
	executable := ld.Macho(code)

	log.Printf("Writing executable to %s", *outputFlag)
	if err := os.WriteFile(*outputFlag, executable, 0755); err != nil {
		log.Fatal(err)
	}
}
*/

/*
func test() {
	// TODO: heap (static data), stack, other syscalls, functions?

	// Examples
	assem.MovImm(reg1, 10)
	assem.MovReg(reg1, reg2)
	// assem.MovMem(reg1, 0x1234)
	// assem.MovToMem(0x1234, reg1)

	assem.Inc(reg1)
	assem.Dec(reg1)

	// Add,Sub,Mul,Div
	assem.AddImm(reg1, 10)
	assem.AddReg(reg1, reg2)


	// Program 1
	l := assem.Label("l1")
	fl := assem.FutureLabel("exit")

	assem.MovImm(reg1, 10)
	assem.Dec(reg1)
	assem.CmpBranch(reg1, 0, l)
	assem.Jump(f1)


	fl.Resolve()
	assem.SyscallExit()

}
*/
