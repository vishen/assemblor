package x64

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

var (
	buf = &binOutput{}
)

type binOutput []byte

func (b *binOutput) bytes() []byte {
	return []byte(*b)
}

func (b *binOutput) write(data ...byte) {
	*b = append(*b, data...)
}

func (b *binOutput) write32(v uint32) {
	tmpBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmpBuf, v)
	b.write(tmpBuf[:4]...)
}

func (b *binOutput) write64(vs ...uint64) {
	for _, v := range vs {
		tmpBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmpBuf, v)
		b.write(tmpBuf[:8]...)
	}
}

func (b *binOutput) writeString(s string, padding int) {
	b.write([]byte(s)...)
	if len(s) < padding {
		for i := 0; i < padding-len(s); i++ {
			b.write(0x00)
		}
	}

}

type MachoSect struct {
	name    string
	segname string
	addr    uint64
	size    uint64
	off     uint32
	align   uint32
	reloc   uint32
	nreloc  uint32
	flag    uint32
	res1    uint32
	res2    uint32
}

type MachoSeg struct {
	name       string
	vsize      uint64
	vaddr      uint64
	fileoffset uint64
	filesize   uint64
	prot1      uint32
	prot2      uint32
	// nsect      uint32
	// msect      uint32
	sect []MachoSect
	flag uint32
}

const (
	MACHO_CPU_AMD64  = 1<<24 | 7
	MACHO_SUBCPU_X86 = 3

	MH_MAGIC_64 = 0xfeedfacf
	MH_EXECUTE  = 0x2
	MH_NOUNDEFS = 0x1

	S_ATTR_PURE_INSTRUCTIONS = 0x80000000
	S_ATTR_SOME_INSTRUCTIONS = 0x00000400

	LC_UNIXTHREAD = 0x5
	LC_SEGMENT_64 = 0x19
)

func MachoWrite(code []byte) {
	// TODO: Will need to do __DATA

	origin := uint64(0x100000000)
	// TODO: Remove hardcoded
	codestart := uint64(0x168)
	loadsize := 328
	nsegs := 3

	// __PAGEZERO segment is used for null pointer deferences
	s1 := &MachoSeg{
		name:  "__PAGEZERO",
		vsize: origin,
	}

	// __TEXT segment contains the machine code
	s2 := &MachoSeg{
		name:       "__TEXT",
		vaddr:      s1.vsize,
		vsize:      0x1000,
		fileoffset: 0,
		filesize:   0x1000,
		prot1:      7,
		prot2:      5,
	}

	/*
		s2.sect = []MachoSect{
			{
				name:    "__text",
				segname: "__TEXT",

				// TODO: Remove hardcoded addr and off.
				addr: s2.vaddr + codestart, // 0x100 is hardcoded for where code starts
				off:  uint32(codestart),

				size: 0x1000 - codestart,
				flag: S_ATTR_PURE_INSTRUCTIONS | S_ATTR_SOME_INSTRUCTIONS,
			},
		}
	*/

	buf.write32(MH_MAGIC_64)      // Magic Number
	buf.write32(MACHO_CPU_AMD64)  // CPU Arch
	buf.write32(MACHO_SUBCPU_X86) // Subcpu?

	buf.write32(MH_EXECUTE) // Permissions

	buf.write32(uint32(nsegs))    // Number of load comamnds (includes segments)
	buf.write32(uint32(loadsize)) // Size of the load section

	buf.write32(uint32(MH_NOUNDEFS))
	buf.write32(0) //Reserved

	s := len(*buf)
	fmt.Println("start here", s)
	for _, s := range []*MachoSeg{s1, s2} {
		size := 72 + 80*len(s.sect)
		fmt.Println("1", len(*buf), size)

		buf.write32(LC_SEGMENT_64)
		buf.write32(uint32(size))
		buf.writeString(s.name, 16)
		buf.write64(s.vaddr)
		buf.write64(s.vsize)
		buf.write64(s.fileoffset)
		buf.write64(s.filesize)
		buf.write32(s.prot1)
		buf.write32(s.prot2)
		buf.write32(uint32(len(s.sect)))
		buf.write32(s.flag)

		for _, t := range s.sect {
			buf.writeString(t.name, 16)
			buf.writeString(t.segname, 16)
			buf.write64(t.addr)
			buf.write64(t.size)
			buf.write32(t.off)
			buf.write32(t.align)
			buf.write32(t.reloc)
			buf.write32(t.nreloc)
			buf.write32(t.flag)
			buf.write32(t.res1) // reserved
			buf.write32(t.res2) // reserved
			buf.write32(0)      // reserved
		}

		align := 8
		offset := len(*buf)
		padding := (align - (offset % align)) % align
		for i := 0; i < padding; i++ {
			buf.write(0x00)
		}

		fmt.Println("2", len(*buf), size, padding)
	}

	// Need to tell mac how to find and start the machine code
	ss := len(*buf)
	buf.write32(LC_UNIXTHREAD)
	buf.write32(184) // command size
	buf.write32(4)   // thread state: x86_THREAD_STATE64
	buf.write32(42)  // word count: x86_EXCEPTION_STATE64_COUNT

	// Setup the initial values for registers
	buf.write64(0x00, 0x00, 0x00, 0x00) // rax, rbx , rcx , rdx
	buf.write64(0x00, 0x00, 0x00, 0x00) // rdi, rsi, rbp, rsp
	buf.write64(0x00, 0x00, 0x00, 0x00) // r8, r9, r10, r11
	buf.write64(0x00, 0x00, 0x00, 0x00) // r12, r13, r14, r15

	// TODO: Don't hardcode code start address!
	buf.write64(uint64(origin+codestart), 0x00, 0x00, 0x00, 0x00) // rip=code, rflags, cs, fs, gs
	fmt.Printf("UNITHREAD: len=%d\n", len(*buf)-ss)

	fmt.Printf("len header: %d - %d = %d\n", len(*buf), s, len(*buf)-s)
	fmt.Printf("codestart=0x%x\n", len(*buf))

	// Write out code
	buf.write(code...)

	align := 4096
	offset := len(*buf)
	padding := (align - (offset % align)) % align
	for i := 0; i < padding; i++ {
		buf.write(0x00)
	}

	if err := os.WriteFile("./mac_bin", buf.bytes(), 0755); err != nil {
		log.Fatal(err)
	}
}
