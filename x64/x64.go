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

func Compile(arch Arch, bc []bytecode.Instruction) []byte {
	o := &output{}
	for _, b := range bc {
		switch b.Instruction() {
		case bytecode.Nop:
			// Do nothing
		case bytecode.MovImm:
			i := b.(bytecode.Imm)
			src := resolveReg(i.Src)
			o.rex(true, false, false, src.isExt())
			o.add(0xC7)
			o.modrm(0x03, 0, src.val())
			o.addImm(uint32(i.Imm))
		case bytecode.MovReg:
			r := b.(bytecode.Reg)
			src := resolveReg(r.Src)
			dst := resolveReg(r.Reg)
			o.rex(true, dst.isExt(), false, src.isExt())
			o.add(0x89)
			o.modrm(0x03, dst.val(), src.val())
		case bytecode.Inc:
			r := b.(bytecode.Imm)
			src := resolveReg(r.Src)
			o.rex(true, false, false, src.isExt())
			o.add(0xFF)
			o.modrm(0x03, 0, src.val())
		case bytecode.Dec:
			r := b.(bytecode.Imm)
			src := resolveReg(r.Src)
			o.rex(true, false, false, src.isExt())
			o.add(0xFF)
			o.modrm(0x03, 0x01, src.val())
		case bytecode.AddImm:
			i := b.(bytecode.Imm)
			src := resolveReg(i.Src)
			imm := uint32(i.Imm)

			// NOTE: src == RAX and imm == 32-bit, then special case
			if src == rax && imm >= 128 {
				// REX.W + 05 id	ADD RAX, imm32
				o.add(0x05)
			} else {
				if imm < 128 {
					// REX.W + 83 /0 ib    ADD r/m64, imm8
					o.add(0x83)
				} else {
					// REX.W + 81 /0 id	ADD r/m64, imm32
					o.add(0x81)
				}
				o.modrm(0x03, 0, src.val())
			}
			o.addImm(imm)
		case bytecode.AddReg:
			r := b.(bytecode.Reg)
			src := resolveReg(r.Src)
			dst := resolveReg(r.Reg)
			o.rex(true, dst.isExt(), false, src.isExt())
			o.add(0x01)
			o.modrm(0x03, dst.val(), src.val())
		case bytecode.SyscallExit:
			s := b.(bytecode.Syscall)
			switch arch {
			case Linux:
				// TODO: Move MovImm into a function to consolidate with bytecode.MovImm?
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
			case Macho:
				{
					src := rax
					o.rex(true, false, false, src.isExt())
					o.add(0xC7)
					o.modrm(0x03, 0, src.val())
					o.addImm(0x02000001)
				}
				{
					src := rdi
					o.rex(true, false, false, src.isExt())
					o.add(0xC7)
					o.modrm(0x03, 0, src.val())
					o.addImm(s.Arg1)
				}
				// Syscall
				o.add(0x0f, 0x05)
			}
		}
	}
	return o.data
}

func resolveReg(reg bytecode.RegType) reg {
	// TODO: Could we allocate registers semi-dynamically favouring
	// reg1 where possible?
	switch reg {
	case bytecode.Reg1:
		return rax
	case bytecode.Reg2:
		return rcx
	case bytecode.Reg3:
		return rdx
	case bytecode.Reg4:
		return rsi
	case bytecode.Reg5:
		return rdi
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
	panic(fmt.Sprintf("%v bytecode register doesn't map to an x64 register"))
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
