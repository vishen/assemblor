package main

import "github.com/vishen/assemblor/x64"

func main() {
	code := []byte{
		0xb8, 0x01, 0x00, 0x00, 0x02, // movl	$33554433, %eax
		0xbf, 0x00, 0x00, 0x00, 0x00, // movl	$0, %edi
		0x0f, 0x05, // syscall
	}
	x64.MachoWrite(code)
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
