package ld

type Linker interface {
	BssAddr() uint32
	Link() []byte
}
