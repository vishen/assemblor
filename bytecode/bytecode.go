package bytecode

import "fmt"

type ImmType uint32

type RegType int

const (
	Reg1 RegType = iota
	Reg2
	Reg3
	Reg4
	Reg5
	Reg6
	Reg7
	Reg8
	Reg9
	Reg10
	Reg11
	Reg12
)

func (r RegType) String() string {
	switch r {
	case Reg1:
		return "Reg1"
	case Reg2:
		return "Reg2"
	case Reg3:
		return "Reg3"
	case Reg4:
		return "Reg4"
	case Reg5:
		return "Reg5"
	case Reg6:
		return "Reg6"
	case Reg7:
		return "Reg7"
	case Reg8:
		return "Reg8"
	case Reg9:
		return "Reg9"
	case Reg10:
		return "Reg10"
	case Reg11:
		return "Reg11"
	case Reg12:
		return "Reg12"
	}
	return fmt.Sprintf("unknown register %d", r)
}

type ConditionalType int

const (
	EQ ConditionalType = iota
)

type InstructionType int

const (
	Invalid InstructionType = iota
	Nop
	MovImm
	MovReg
	Inc
	Dec
	AddImm
	AddReg
	Jmp
	Label
	SyscallExit
)

func (i InstructionType) String() string {
	switch i {
	case Invalid:
		return "Invalid"
	case Nop:
		return "Nop"
	case MovImm:
		return "MovImm"
	case MovReg:
		return "MovReg"
	case Inc:
		return "Inc"
	case Dec:
		return "Dec"
	case AddImm:
		return "AddImm"
	case AddReg:
		return "AddReg"
	case Jmp:
		return "Jmp"
	case Label:
		return "Label"
	case SyscallExit:
		return "SyscallExit"
	}
	return fmt.Sprintf("unknown instruction type %d", i)
}

type Instruction interface {
	Instruction() InstructionType
}

type LabelType uint32

func (l LabelType) Instruction() InstructionType {
	return Label
}

func (l LabelType) String() string {
	return fmt.Sprintf("Label %d", l)
}

type Imm struct {
	Inst InstructionType
	Src  RegType
	Imm  ImmType
}

func (b Imm) Instruction() InstructionType { return b.Inst }

type Reg struct {
	Inst InstructionType
	Src  RegType
	Reg  RegType
}

func (b Reg) Instruction() InstructionType { return b.Inst }

type Branch struct {
	ID    int
	Inst  InstructionType
	Label LabelType
}

func (b Branch) Instruction() InstructionType { return b.Inst }

type Syscall struct {
	Inst InstructionType
	Arg1 uint32
}

func (b Syscall) Instruction() InstructionType { return b.Inst }
