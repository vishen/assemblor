package bytecode

import "fmt"

type AddrType uint64

func (a AddrType) Offset(offset int) AddrType {
	return AddrType(uint64(a) + uint64(offset))
}

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
	MovAddr
	WriteImm
	WriteReg
	Inc
	Dec
	AddImm
	AddReg
	CmpImm
	CmpReg
	Jmp
	BranchCond
	Label
	SyscallExit
	SyscallWrite
	ReserveBytes
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
	case MovAddr:
		return "MovAddr"
	case WriteImm:
		return "WriteImm"
	case WriteReg:
		return "WriteReg"
	case Inc:
		return "Inc"
	case Dec:
		return "Dec"
	case AddImm:
		return "AddImm"
	case AddReg:
		return "AddReg"
	case CmpImm:
		return "CmpImm"
	case CmpReg:
		return "CmpReg"
	case Jmp:
		return "Jmp"
	case BranchCond:
		return "BranchCond"
	case Label:
		return "Label"
	case SyscallExit:
		return "SyscallExit"
	case SyscallWrite:
		return "SyscallWrite"
	case ReserveBytes:
		return "ReserveBytes"
	}
	return fmt.Sprintf("Unknown(%d)", i)
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
	Inst    InstructionType
	DstReg  RegType
	DstAddr AddrType
	Imm     ImmType
}

func (b Imm) Instruction() InstructionType { return b.Inst }

type Reg struct {
	Inst    InstructionType
	DstReg  RegType
	DstAddr AddrType
	Reg     RegType
}

func (b Reg) Instruction() InstructionType { return b.Inst }

type Addr struct {
	Inst InstructionType
	Dst  RegType
	Addr AddrType
}

func (b Addr) Instruction() InstructionType { return b.Inst }

type Branch struct {
	ID           int
	Inst         InstructionType
	Reg1         RegType
	Cond         ConditionalType
	Reg2         RegType
	JmpTrueLabel LabelType
}

func (b Branch) Instruction() InstructionType { return b.Inst }

type Syscall struct {
	Inst InstructionType
	Reg1 RegType
	Reg2 RegType
}

func (b Syscall) Instruction() InstructionType { return b.Inst }

type Data struct {
	Inst InstructionType
	Arg1 uint32
}

func (b Data) Instruction() InstructionType { return b.Inst }
