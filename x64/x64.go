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

	movReg := func(dst, src reg) {
		o.rex(true, src.isExt(), false, dst.isExt())
		o.add(0x89)
		o.modrm(modVal, src.val(), dst.val())
	}

	movImm := func(dst reg, imm uint32) {
		o.rex(true, false, false, dst.isExt())
		o.add(0xC7)
		o.modrm(modVal, 0, dst.val())
		o.addImm(imm)
	}

	cmpReg := func(r1, r2 reg) {
		o.rex(true, r2.isExt(), false, r1.isExt())
		o.add(0x39)
		o.modrm(modVal, r2.val(), r1.val())
	}

	pushReg := func(rs ...reg) {
		for _, r1 := range rs {
			if r1.isExt() {
				o.add(0x41)
			}
			o.add(0x50 + r1.val())
		}
	}

	popReg := func(rs ...reg) {
		for _, r1 := range rs {
			if r1.isExt() {
				o.add(0x41)
			}
			o.add(0x58 + r1.val())
		}
	}

	bssSize := uint64(0)

	type branch struct {
		offset, offsetToWrite int
	}

	labels := make(map[bytecode.LabelType]int)
	branches := make(map[int]branch)
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
			offset := o.offset()
			branches[b.ID] = branch{offset, offset - 4} // -4 is the length of 0x00 to fill in later
		case bytecode.BranchCond:
			b := b.(bytecode.Branch)
			// Currently everything is just assumed to be a 32 bit displacement
			// sign extended to 64 bits
			// TODO: Find the proper way to do this better
			r1 := resolveReg(b.Reg1)
			r2 := resolveReg(b.Reg2)
			cmpReg(r1, r2)

			var opcode byte
			switch b.Cond {
			case bytecode.EQ:
				opcode = 0x84
			case bytecode.NEQ:
				opcode = 0x85
			default:
				panic(fmt.Sprintf("unknown conditional type for branch conditional: %v", b))
			}
			o.add(0x0F, opcode, 0x00, 0x00, 0x00, 0x00)
			offset := o.offset()
			branches[b.ID] = branch{offset, offset - 4} // -4 is the length of 0x00 to fill in later
		case bytecode.MovImm:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			movImm(dst, uint32(i.Imm))
		case bytecode.MovReg:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			src := resolveReg(i.Reg)
			movReg(dst, src)
		case bytecode.MovAddr:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			movImm(dst, uint32(uint64(i.Addr)+bssAddr))
		case bytecode.MovMem:
			// 48 8b 0b                    	mov	rcx, qword ptr [rbx]
			// 48 8b 08                    	mov	rcx, qword ptr [rax]
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			src := resolveReg(i.Mem)
			o.rex(true, dst.isExt(), false, src.isExt())
			o.add(0x8b)
			o.modrm(modMemory, dst.val(), src.val())
		case bytecode.WriteRegToAddr:
			// mov qword [0xdead], rbx
			// 48 89 1c 25 ad de 00 00     	mov	qword ptr [57005], rbx
			i := b.(bytecode.Inst)
			r1 := resolveReg(i.Reg)
			o.rex(true, r1.isExt(), false, false)
			o.add(0x89)
			o.modrm(modMemory, r1.val(), 0x04) // 0x00 = [rax]
			o.add(0x25)
			o.addImm(uint32(uint64(i.DstAddr) + bssAddr))
		case bytecode.WriteRegToMem:
			// mov [rbx], rcx
			// 48 89 0b                    	mov	qword ptr [rbx], rcx
			i := b.(bytecode.Inst)
			r1 := resolveReg(i.DstMem)
			r2 := resolveReg(i.Reg)
			o.rex(true, r2.isExt(), false, r1.isExt())
			o.add(0x89)
			o.modrm(modMemory, r2.val(), r1.val()) // 0x00 = [rax]
		case bytecode.WriteImm:
			// mov qword [0xdeadbe], 0x1234
			// 48 c7 04 25 be ad de 00 34 12 00 00 	mov	qword ptr [14593470], 4660
			i := b.(bytecode.Inst)
			o.rex(true, false, false, false)
			o.add(0xC7)
			o.modrm(modMemory, 0x00, 0x04)
			o.add(0x25)
			o.addImm(uint32(uint64(i.DstAddr) + bssAddr))
			o.addImm(uint32(i.Imm))
		case bytecode.Inc:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			o.rex(true, false, false, dst.isExt())
			o.add(0xFF)
			o.modrm(modVal, 0, dst.val())
		case bytecode.Dec:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			o.rex(true, false, false, dst.isExt())
			o.add(0xFF)
			o.modrm(modVal, 0x01, dst.val())
		case bytecode.SubImm:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			imm := uint32(i.Imm)
			o.rex(true, false, false, dst.isExt())
			if dst == rax && imm >= 128 {
				o.add(0x2d)
			} else {
				o.add(0x81)
				o.modrm(modVal, 0x05, dst.val())
			}
			o.addImm(imm)
		case bytecode.AddImm:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			imm := uint32(i.Imm)
			o.rex(true, false, false, dst.isExt())
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
				o.modrm(modVal, 0, dst.val())
			}
			o.addImm(imm)
		case bytecode.AddReg:
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			src := resolveReg(i.Reg)
			o.rex(true, src.isExt(), false, dst.isExt())
			o.add(0x01)
			o.modrm(modVal, src.val(), dst.val())
		case bytecode.CmpReg:
			i := b.(bytecode.Inst)
			r1 := resolveReg(i.DstReg)
			r2 := resolveReg(i.Reg)
			cmpReg(r1, r2)
		case bytecode.CmpImm:
			// TODO: Refactor with AddImm?
			i := b.(bytecode.Inst)
			dst := resolveReg(i.DstReg)
			imm := uint32(i.Imm)

			o.rex(true, false, false, dst.isExt())
			// NOTE: dst == RAX and imm == 32-bit, then special case
			if dst == rax && imm >= 128 {
				// REX.W + 3D id	CMP RAX, imm32
				o.add(0x3d)
			} else {
				o.add(0x81)
				o.modrm(modVal, 0x07, dst.val())
			}
			o.addImm(imm)
		case bytecode.PushImm:
			i := b.(bytecode.Inst)
			o.add(0x68)
			o.addImm(uint32(i.Imm))
		case bytecode.PushReg:
			i := b.(bytecode.Inst)
			src := resolveReg(i.Reg)
			pushReg(src)
		case bytecode.PopReg:
			i := b.(bytecode.Inst)
			src := resolveReg(i.Reg)
			popReg(src)
		case bytecode.SyscallExit:
			s := b.(bytecode.Syscall)
			switch arch {
			case Linux:
				movReg(rbx, resolveReg(s.Reg1))
				movImm(rax, 0x01)
				o.add(0xcd, 0x80)
			case Macho:
				movReg(rdi, resolveReg(s.Reg1))
				movImm(rax, 0x02000001)
				o.add(0x0f, 0x05)
			}
		case bytecode.SyscallWrite:
			s := b.(bytecode.Syscall)
			switch arch {
			case Linux:
				// https://stackoverflow.com/questions/2535989/what-are-the-calling-conventions-for-unix-linux-system-calls-and-user-space-f
				pushReg(rcx, rdx, rax, rbx, rdi, rsi, r8, r9, r10, r11)

				movReg(rcx, resolveReg(s.Reg1)) // ptr
				movReg(rdx, resolveReg(s.Reg2)) // length
				movImm(rax, 0x04)               // sys_write
				movImm(rbx, 0x01)               // stdout
				o.add(0xcd, 0x80)

				popReg(r11, r10, r9, r8, rsi, rdi, rbx, rax, rdx, rcx)
			case Macho:
				/*
					1   mov    rax, 0x02000004    ; system call for write
					2   mov    rdi, 1             ; file descriptor 1 is stdout
					3   mov    rsi, qword message ; get string address
					4   mov    rdx, 13            ; number of bytes
					5   syscall                   ; execute syscall (write)
				*/

				pushReg(rcx, rdx, rax, rbx, rdi, rsi, r8, r9, r10, r11)

				movReg(rsi, resolveReg(s.Reg1))
				movReg(rdx, resolveReg(s.Reg2))
				movImm(rax, 0x02000004)
				movImm(rdi, 0x01)
				o.add(0x0f, 0x05)

				popReg(r11, r10, r9, r8, rsi, rdi, rbx, rax, rdx, rcx)
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
		case bytecode.Jmp, bytecode.BranchCond:
			b := b.(bytecode.Branch)
			br := branches[b.ID]
			labelOffset := labels[b.JmpTrueLabel]
			if jmpDiff := uint32(labelOffset - br.offset); jmpDiff > 0 {
				o.fill32(br.offsetToWrite, jmpDiff)
			} else {
				o.fill32(br.offsetToWrite, 0xffffffff+jmpDiff)
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
