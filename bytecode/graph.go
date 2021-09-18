package bytecode

import "fmt"

type Graph struct {
	inst []Instruction

	labeln  int
	branchn int

	dataSize uint32
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) offset() int {
	return len(g.inst)
}

func (g *Graph) branchID() int {
	bid := g.branchn
	g.branchn++
	return bid
}

func (g *Graph) Print() {
	for i, b := range g.inst {
		fmt.Printf("%d: %v\n", i, b)
	}
}

func (g *Graph) Bytecode() ([]Instruction, error) {
	// TODO: Validate that there is no unresolved labels.
	return g.inst, nil
}

func (g *Graph) WriteImm(a AddrType, i ImmType) {
	g.inst = append(g.inst, Imm{
		Inst:    WriteImm,
		DstAddr: a,
		Imm:     i,
	})
}

func (g *Graph) WriteReg(a AddrType, r RegType) {
	g.inst = append(g.inst, Reg{
		Inst:    WriteReg,
		DstAddr: a,
		Reg:     r,
	})
}

func (g *Graph) MovImm(r RegType, i ImmType) {
	g.inst = append(g.inst, Imm{
		Inst:   MovImm,
		DstReg: r,
		Imm:    i,
	})
}

func (g *Graph) MovAddr(r RegType, a AddrType) {
	g.inst = append(g.inst, Addr{
		Inst: MovAddr,
		Dst:  r,
		Addr: a,
	})
}

func (g *Graph) MovReg(r RegType, r2 RegType) {
	g.inst = append(g.inst, Reg{
		Inst:   MovReg,
		DstReg: r,
		Reg:    r2,
	})
}

func (g *Graph) Inc(r RegType) {
	g.inst = append(g.inst, Imm{
		Inst:   Inc,
		DstReg: r,
	})
}

func (g *Graph) Dec(r RegType) {
	g.inst = append(g.inst, Imm{
		Inst:   Dec,
		DstReg: r,
	})
}

func (g *Graph) AddImm(r RegType, i ImmType) {
	g.inst = append(g.inst, Imm{
		Inst:   AddImm,
		DstReg: r,
		Imm:    i,
	})
}

func (g *Graph) AddReg(r RegType, r2 RegType) {
	g.inst = append(g.inst, Reg{
		Inst:   AddReg,
		DstReg: r,
		Reg:    r2,
	})
}
func (g *Graph) Label() LabelType {
	lt := LabelType(g.labeln)
	g.labeln++
	g.inst = append(g.inst, lt)
	return lt
}

func (g *Graph) FutureLabel() LabelType {
	lt := LabelType(g.labeln)
	g.labeln++
	return lt
}

func (g *Graph) ResolveLabel(l LabelType) {
	g.inst = append(g.inst, l)
}

func (g *Graph) CmpImm(r RegType, i ImmType) {
	g.inst = append(g.inst, Imm{
		Inst:   CmpImm,
		DstReg: r,
		Imm:    i,
	})
}

func (g *Graph) CmpReg(r RegType, r2 RegType) {
	g.inst = append(g.inst, Reg{
		Inst:   CmpReg,
		DstReg: r,
		Reg:    r2,
	})
}

func (g *Graph) Jmp(l LabelType) {
	g.inst = append(g.inst, Branch{
		ID:           g.branchID(),
		Inst:         Jmp,
		JmpTrueLabel: l,
	})
}

func (g *Graph) BranchCond(r1 RegType, c ConditionalType, r2 RegType, l LabelType) {
	g.inst = append(g.inst, Branch{
		ID:           g.branchID(),
		Inst:         BranchCond,
		Reg1:         r1,
		Cond:         c,
		Reg2:         r2,
		JmpTrueLabel: l,
	})
}

func (g *Graph) ReserveBytes(size uint32) AddrType {
	g.inst = append(g.inst, Data{
		Inst: ReserveBytes,
		Arg1: size,
	})

	bssSize := g.dataSize
	g.dataSize += size
	return AddrType(bssSize)
}

func (g *Graph) SyscallExit(statusCodePtr RegType) {
	g.inst = append(g.inst, Syscall{
		Inst: SyscallExit,
		Reg1: statusCodePtr,
	})
}

func (g *Graph) SyscallWrite(dataPtr RegType, length RegType) {
	g.inst = append(g.inst, Syscall{
		Inst: SyscallWrite,
		Reg1: dataPtr,
		Reg2: length,
	})
}
