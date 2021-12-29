package bytecode

import "fmt"

type AddrType uint64

func (a AddrType) Offset(offset int) AddrType {
	return AddrType(uint64(a) + uint64(offset))
}

type ImmType uint32

type RegType int

const (
	RegUnknown RegType = iota
	Reg1
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
	return "none"
}

type ConditionalType int

const (
	UnknownConditional ConditionalType = iota
	EQ
	NEQ
)

func (c ConditionalType) String() string {
	switch c {
	case EQ:
		return "=="
	case NEQ:
		return "!="
	}
	return "none"
}

type InstructionType int

const (
	Invalid InstructionType = iota
	Nop
	MovImm         // move imm -> reg
	MovReg         // move reg -> reg
	MovAddr        // move addr -> reg
	MovMem         // move [reg] -> reg
	WriteImm       // write imm -> addr
	WriteRegToAddr // write reg -> addr
	WriteRegToMem  // write reg -> [reg]
	Inc            // increment reg
	Dec            // decrement reg
	SubImm         // sub imm -> reg
	AddImm         // add imm -> reg
	AddReg         // add reg -> reg
	CmpImm         // cmp imm -> reg
	CmpReg         // cmp reg -> reg
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
	case MovMem:
		return "MovMem"
	case WriteImm:
		return "WriteImm"
	case WriteRegToMem:
		return "WriteRegToMem"
	case WriteRegToAddr:
		return "WriteRegToAddr"
	case Inc:
		return "Inc"
	case Dec:
		return "Dec"
	case SubImm:
		return "SubImm"
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
	return "none"
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

type Inst struct {
	inst InstructionType

	DstReg  RegType
	DstAddr AddrType
	DstMem  RegType

	Imm  ImmType
	Reg  RegType
	Addr AddrType
	Mem  RegType
}

func (i Inst) String() string {
	switch it := i.Instruction(); it {
	case MovImm, AddImm:
		return fmt.Sprintf("%v %v %v", it, i.DstReg, i.Imm)
	case MovReg, Inc, Dec, AddReg:
		return fmt.Sprintf("%v %v %v", it, i.DstReg, i.Reg)
	case MovAddr:
		return fmt.Sprintf("%v %v %v", it, i.DstReg, i.Addr)
	case MovMem:
		return fmt.Sprintf("%v %v %v", it, i.DstReg, i.Mem)
	case WriteImm:
		return fmt.Sprintf("%v %v %v", it, i.DstAddr, i.Imm)
	case WriteRegToAddr:
		return fmt.Sprintf("%v %v %v", it, i.DstAddr, i.Reg)
	case WriteRegToMem:
		return fmt.Sprintf("%v %v %v", it, i.Mem, i.Reg)
	}
	return fmt.Sprintf("%v: dst_reg=%v dst_addr=%v dst_mem=%v | imm=%v reg=%v addr=%v mem=%v", i.inst, i.DstReg, i.DstAddr, i.DstMem, i.Imm, i.Reg, i.Addr, i.Mem)
}

func (i Inst) Instruction() InstructionType { return i.inst }

type Branch struct {
	ID           int
	Inst         InstructionType
	Reg1         RegType
	Cond         ConditionalType
	Reg2         RegType
	JmpTrueLabel LabelType
}

func (b Branch) String() string {
	// TODO: Do this better
	// For 'jmp' and non-conditional branches
	if b.Reg1 == 0 {
		return fmt.Sprintf("%v (%d): %v", b.Inst, b.ID, b.JmpTrueLabel)
	}
	return fmt.Sprintf("%v (%d): %v %v %v -> %v", b.Inst, b.ID, b.Reg1, b.Cond, b.Reg2, b.JmpTrueLabel)
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
