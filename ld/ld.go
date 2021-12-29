package ld

type Linker interface {
	BssAddr() uint64
	Link(code []byte, bssSize uint64) []byte
}
