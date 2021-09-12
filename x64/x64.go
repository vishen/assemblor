package x64

import (
	"fmt"

	"github.com/vishen/assemblor/bytecode"
)

type Arch int

const (
	Linux Arch = iota
	Macho
)

func Compile(arch Arch, bc []bytecode.Instruction, bssAddr uint64) ([]byte, uint64) {
	o := &output{}

	fmt.Println("BSS ADDR", bssAddr)

	movReg := func(dst, src reg) {
		o.rex(true, src.isExt(), false, dst.isExt())
		o.add(0x89)
		o.modrm(0x03, src.val(), dst.val())
	}

	movImm := func(dst reg, imm uint32) {
		o.rex(true, false, false, dst.isExt())
		o.add(0xC7)
		o.modrm(0x03, 0, dst.val())
		o.addImm(imm)
	}

	bssSize := uint64(0)

	labels := make(map[bytecode.LabelType]int)
	branches := make(map[int]int)
	for j, b := range bc {
		switch b.Instruction() {
		case bytecode.Invalid:
			fmt.Printf("invalid bytecode found at %d", j)
		case bytecode.Nop:
			// Do nothing
		case bytecode.Label:
			l := b.(bytecode.LabelType)
			labels[l] = o.offset()
		case bytecode.Jmp:
			b := b.(bytecode.Branch)
			// Currently everything is just assumed to be a 32 bit displacement
			// sign extended to 64 bits
			// TODO: Find the proper way to do this better
			o.add(0xe9, 0x00, 0x00, 0x00, 0x00)
			branches[b.ID] = o.offset()
		case bytecode.MovImm:
			i := b.(bytecode.Imm)
			dst := resolveReg(i.DstReg)
			movImm(dst, uint32(i.Imm))
		case bytecode.MovReg:
			r := b.(bytecode.Reg)
			dst := resolveReg(r.Dst)
			src := resolveReg(r.Reg)
			movReg(dst, src)
		case bytecode.MovAddr:
			a := b.(bytecode.Addr)
			dst := resolveReg(a.Dst)
			movImm(dst, uint32(uint64(a.Addr)+bssAddr))
			/*
				o.rex(true, dst.isExt(), false, false)
				o.add(0x8b)
				o.modrm(0x00, dst.val(), 0x04)
				o.add(0x25)
				o.addImm(uint32(uint64(a.Addr) + bssAddr))
			*/
		case bytecode.WriteImm:
			// mov qword [0xdeadbe], 0x1234
			// 48 c7 04 25 be ad de 00 34 12 00 00 	mov	qword ptr [14593470], 4660
			i := b.(bytecode.Imm)
			o.rex(true, false, false, false)
			o.add(0xC7)
			o.modrm(0x00, 0x00, 0x04)
			o.add(0x25)
			o.addImm(uint32(uint64(i.DstAddr) + bssAddr))
			o.addImm(uint32(i.Imm))
		case bytecode.Inc:
			r := b.(bytecode.Imm)
			dst := resolveReg(r.DstReg)
			o.rex(true, false, false, dst.isExt())
			o.add(0xFF)
			o.modrm(0x03, 0, dst.val())
		case bytecode.Dec:
			r := b.(bytecode.Imm)
			dst := resolveReg(r.DstReg)
			o.rex(true, false, false, dst.isExt())
			o.add(0xFF)
			o.modrm(0x03, 0x01, dst.val())
		case bytecode.AddImm:
			i := b.(bytecode.Imm)
			dst := resolveReg(i.DstReg)
			imm := uint32(i.Imm)

			// NOTE: dst == RAX and imm == 32-bit, then special case
			if dst == rax && imm >= 128 {
				// REX.W + 05 id	ADD RAX, imm32
				o.add(0x05)
			} else {
				/*
					// TODO: Fix? AddImm always does a uint32, but for the <128 case
					// it should only be a single byte...
					if imm < 128 {
						// REX.W + 83 /0 ib    ADD r/m64, imm8
						o.add(0x83)
					} else {
						// REX.W + 81 /0 id	ADD r/m64, imm32
						o.add(0x81)
					}
				*/
				o.add(0x81)
				o.modrm(0x03, 0, dst.val())
			}
			o.addImm(imm)
		case bytecode.AddReg:
			r := b.(bytecode.Reg)
			dst := resolveReg(r.Dst)
			src := resolveReg(r.Reg)
			o.rex(true, src.isExt(), false, dst.isExt())
			o.add(0x01)
			o.modrm(0x03, src.val(), dst.val())
		case bytecode.SyscallExit:
			s := b.(bytecode.Syscall)
			switch arch {
			case Linux:
				// TODO: Move MovImm into a function to consolidate with bytecode.MovImm?
				/*
					// TODO: Make work on linux
					{
						src := rax
						o.rex(true, false, false, src.isExt())
						o.add(0xC7)
						o.modrm(0x03, 0, src.val())
						o.addImm(1)
					}
					{
						src := rbx
						o.rex(true, false, false, src.isExt())
						o.add(0xC7)
						o.modrm(0x03, 0, src.val())
						o.addImm(s.Arg1)
					}
					// Syscall
					o.add(0xcd, 0x80)
				*/
			case Macho:
				movReg(rdi, resolveReg(s.Reg1))
				movImm(rax, 0x02000001)
				// Syscall
				o.add(0x0f, 0x05)
			}
		case bytecode.SyscallWrite:
			s := b.(bytecode.Syscall)
			switch arch {
			case Linux:
				// TODO:
			case Macho:
				/*
					1   mov    rax, 0x02000004    ; system call for write
					2   mov    rdi, 1             ; file descriptor 1 is stdout
					3   mov    rsi, qword message ; get string address
					4   mov    rdx, 13            ; number of bytes
					5   syscall                   ; execute syscall (write)
				*/

				movReg(rsi, resolveReg(s.Reg1))
				movReg(rdx, resolveReg(s.Reg2))
				movImm(rax, 0x02000004)
				movImm(rdi, 0x01)
				// Syscall
				o.add(0x0f, 0x05)
			}
		case bytecode.ReserveBytes:
			r := b.(bytecode.Data)
			bssSize += uint64(r.Arg1)
		default:
			// TODO: What to do in case of missing instruction
			panic(fmt.Sprintf("unhandled bytecode instruction %v", b))
		}
	}

	// Resolve all the jmps
	for _, b := range bc {
		switch b.Instruction() {
		case bytecode.Jmp:
			b := b.(bytecode.Branch)
			offset := branches[b.ID]
			offsetToWrite := offset - 4 // len of the jmp instruction
			labelOffset := labels[b.Label]
			if jmpDiff := uint32(labelOffset - offset); jmpDiff > 0 {
				o.fill32(offsetToWrite, jmpDiff)
			} else {
				o.fill32(offsetToWrite, 0xffffffff+jmpDiff)
			}
		}
	}
	return o.data, bssSize
}

func resolveReg(reg bytecode.RegType) reg {
	// TODO: Could we allocate registers semi-dynamically favouring
	// reg1 where possible?
	// https://wiki.cdot.senecacollege.ca/wiki/X86_64_Register_and_Instruction_Quick_Start
	switch reg {
	case bytecode.Reg1:
		return rax
	case bytecode.Reg2:
		return rcx
	case bytecode.Reg3:
		return rdx
	case bytecode.Reg4:
		return rbx
	case bytecode.Reg5:
		return r15
	case bytecode.Reg6:
		return r8
	case bytecode.Reg7:
		return r9
	case bytecode.Reg8:
		return r10
	case bytecode.Reg9:
		return r11
	case bytecode.Reg10:
		return r12
	case bytecode.Reg11:
		return r13
	case bytecode.Reg12:
		return r14
	}
	panic(fmt.Sprintf("%v bytecode register doesn't map to an x64 register", reg))
}

type reg int8

func (r reg) isExt() bool {
	return r&8 == 8
}
func (r reg) val() byte {
	// TODO: I am sure there is a bit manipulation way to do this?
	if r < 8 {
		return byte(r)
	} else {
		return byte(r ^ 8)
	}
}

const (
	rax reg = iota
	rcx
	rdx
	rbx
	rsp
	rbp
	rsi
	rdi
	r8
	r9
	r10
	r11
	r12
	r13
	r14
	r15
)
