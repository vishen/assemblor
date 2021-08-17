package ld

import "encoding/binary"

var (
	buf = &output{}
)

type output []byte

func (b *output) offset() int {
	return len(*b)
}

func (b *output) bytes() []byte {
	return []byte(*b)
}

func (b *output) align(a uint64) {
	offset := b.offset()
	padding := (a - (uint64(offset) % a)) % a
	for i := uint64(0); i < padding; i++ {
		buf.write(0x00)
	}
}

func (b *output) fill(pos int, data ...byte) {
	for i, d := range data {
		[]byte(*b)[pos+i] = d
	}
}

func (b *output) write(data ...byte) {
	*b = append(*b, data...)
}

func (b *output) fill32(pos int, v uint32) {
	tmpBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpBuf, v)
	b.fill(pos, tmpBuf...)
}

func (b *output) write32(v uint32) {
	tmpBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpBuf, v)
	b.write(tmpBuf[:4]...)
}

func (b *output) write64(vs ...uint64) {
	for _, v := range vs {
		tmpBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmpBuf, v)
		b.write(tmpBuf[:8]...)
	}
}

func (b *output) fill64(pos int, v uint64) {
	tmpBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(tmpBuf, v)
	b.fill(pos, tmpBuf...)
}

func (b *output) writeString(s string, padding int) {
	b.write([]byte(s)...)
	if len(s) < padding {
		for i := 0; i < padding-len(s); i++ {
			b.write(0x00)
		}
	}

}
