package bytecode

import "fmt"

type AddrType uint64

func (a AddrType) Offset(offset int) AddrType {
	return AddrType(uint64(a) + uint64(offset))
}

type ImmType uint32

type RegType string

const (
	RegUnknown RegType = "RegUnknown"
	Reg1       RegType = "Reg1"
	Reg2       RegType = "Reg2"
	Reg3       RegType = "Reg3"
	Reg4       RegType = "Reg4"
	Reg5       RegType = "Reg5"
	Reg6       RegType = "Reg6"
	Reg7       RegType = "Reg7"
	Reg8       RegType = "Reg8"
	Reg9       RegType = "Reg9"
	Reg10      RegType = "Reg10"
	Reg11      RegType = "Reg11"
	Reg12      RegType = "Reg12"
)

type ConditionalType string

const (
	UnknownConditional ConditionalType = "UnknownConditional"
	EQ                 ConditionalType = "EQ"
	NEQ                ConditionalType = "NEQ"
)

type InstructionType string

const (
	Invalid        InstructionType = "Invalid"
	Nop            InstructionType = "Nop"
	MovImm         InstructionType = "MovImm"         // move imm -> reg
	MovReg         InstructionType = "MovReg"         // move reg -> reg
	MovAddr        InstructionType = "MovAddr"        // move addr -> reg
	MovMem         InstructionType = "MovMen"         // move [reg] -> reg
	WriteImm       InstructionType = "WriteImm"       // write imm -> addr
	WriteRegToAddr InstructionType = "WriteRegToAddr" // write reg -> addr
	WriteRegToMem  InstructionType = "WriteRegToMem"  // write reg -> [reg]
	Inc            InstructionType = "Inc"            // increment reg
	Dec            InstructionType = "Dec"            // decrement reg
	SubImm         InstructionType = "SubImm"         // sub imm -> reg
	AddImm         InstructionType = "AddImm"         // add imm -> reg
	AddReg         InstructionType = "AddReg"         // add reg -> reg
	CmpImm         InstructionType = "CmpImm"         // cmp imm -> reg
	CmpReg         InstructionType = "CmpReg"         // cmp reg -> reg
	Jmp            InstructionType = "Jmp"
	BranchCond     InstructionType = "BrandCond"
	Label          InstructionType = "Label"
	PushImm        InstructionType = "PushImm"
	PushReg        InstructionType = "PushReg"
	PopReg         InstructionType = "PopReg"
	SyscallExit    InstructionType = "SyscallExit"
	SyscallWrite   InstructionType = "SyscallWrite"
	ReserveBytes   InstructionType = "ReserveBytes"
)

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
	if b.Reg1 == RegUnknown {
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
