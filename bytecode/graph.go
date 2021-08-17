package bytecode

type Graph struct {
	inst []Instruction
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) Bytecode() []Instruction {
	return g.inst
}

func (g *Graph) MovImm(r RegType, i ImmType) {
	g.inst = append(g.inst, Imm{
		Inst: MovImm,
		Src:  r,
		Imm:  i,
	})
}

func (g *Graph) SyscallExit(statusCode uint32) {
	g.inst = append(g.inst, Syscall{
		Inst: SyscallExit,
		Arg1: statusCode,
	})
}
