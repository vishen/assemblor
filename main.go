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
	g.MovImm(bytecode.Reg1, 0xdeadbeef)
	g.SyscallExit(125)

	code := x64.Compile(x64.Macho, g.Bytecode())
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
