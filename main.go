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

func main() {
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
