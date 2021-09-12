package x64

import "encoding/binary"

type output struct {
	data []byte
}

func (o *output) offset() int {
	return len(o.data)
}

func (o *output) fill32(index int, val uint32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, val)
	for i, d := range buf {
		o.data[index+i] = d
	}
}

func (o *output) add(b ...byte) {
	o.data = append(o.data, b...)
}

func (o *output) addImm(imm uint32) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, imm)
	o.data = append(o.data, buf...)
}

func (o *output) addImm64(imm uint64) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, imm)
	o.data = append(o.data, buf...)
}

func (o *output) rex(operand64Bit, regExt, sibIndexExt, rmExt bool) {
	var rex byte = 0x40 // REX prefix

	if operand64Bit {
		rex |= 1 << 3
	}
	if regExt {
		rex |= 1 << 2
	}
	if sibIndexExt {
		rex |= 1 << 1
	}
	if rmExt {
		rex |= 1
	}
	o.data = append(o.data, rex)
}

func (o *output) modrm(mod byte, reg byte, rm byte) {
	var modrm byte = 0x0
	modrm |= (rm | (reg << 3) | (mod << 6))
	o.data = append(o.data, modrm)
}
