package bytecode

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

type InstructionType int

const (
	Nop InstructionType = iota
	MovImm
	MovReg
	Inc
	Dec
	AddImm
	AddReg
	SyscallExit
)

type Instruction interface {
	Instruction() InstructionType
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

type Syscall struct {
	Inst InstructionType
	Arg1 uint32
}

func (b Syscall) Instruction() InstructionType { return b.Inst }
