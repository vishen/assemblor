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

func (b *binOutput) write64(v uint64) {
	tmpBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(tmpBuf, v)
	b.write(tmpBuf[:8]...)
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

	LC_SEGMENT_64 = 0x19
)

func MachoWrite(code []byte) {
	// TODO: Will need to do __DATA

	// Number of segments
	segs := 0

	// Number of sections
	sects := 0

	/*
		___pagezerostart:
		   dd 0x19         ; LC_SEGMENT_64
		   dd ___pagezeroend - ___pagezerostart    ; command size
		   db '__PAGEZERO',0,0,0,0,0,0 ; segment name (pad to 16 bytes)
		   dq 0            ; VM address
		   dq 0x100000000  ; VM size
		   dq 0            ; file offset
		   dq 0            ; file size
		   dd 0x0          ; VM_PROT_NONE (maximum protection)
		   dd 0x0          ; VM_PROT_NONE (inital protection)
		   dd 0            ; number of sections
		   dd 0x0          ; flags
		   align 8, db 0   ; pad with zero to 8-byte boundary
		___pagezeroend:
	*/

	// __PAGEZERO segment is used for null pointer deferences
	s1 := &MachoSeg{
		name:  "__PAGEZERO",
		vsize: 0x100000000,
	}
	segs++

	/*
		___TEXTstart:
			dd 0x19         ; LC_SEGMENT_64
			dd ___TEXTend - ___TEXTstart    ; command size
			db '__TEXT',0,0,0,0,0,0,0,0,0,0 ; segment name (pad to 16 bytes)
			dq 0x100000000  ; VM address
			dq 0x1000       ; VM size
			dq 0            ; file offset
			dq 0x1000       ; file size
			dd 0x7          ; VM_PROT_READ | VM_PROT_WRITE | VM_PROT_EXECUTE
			dd 0x5          ; VM_PROT_READ | VM_PROT_EXECUTE
			dd 6            ; number of sections
			dd 0x0          ; flags

			SIZE:
			dd * 6
			dq * 4
			16 * 1
			= 72
		___TEXTtextstart:
			db '__text',0,0,0,0,0,0,0,0,0,0 ; section name (pad to 16 bytes)
			db '__TEXT',0,0,0,0,0,0,0,0,0,0 ; segment name (pad to 16 bytes)
			dq 0x100000000 + ___codestart - ___TEXTload ; address
			dq ___codeend - ___codestart    ; size
			dd ___codestart ; offset
			dd 0            ; alignment as power of 2 (1)
			dd 0            ; relocations data offset
			dd 0            ; number of relocations
			dd 0x80000400   ; S_REGULAR | S_ATTR_PURE_INSTRUCTIONS | S_ATTR_SOME_INSTRUCTIONS
			dd 0            ; reserved1
			dd 0            ; reserved2
			dd 0            ; reserved3

			SIZE:
			16 * 2 = 32
			dq * 2 = 12
			dd * 8 = 32
			= 80
		___TEXTend:
	*/

	// __TEXT segment contains the machine code
	s2 := &MachoSeg{
		name:  "__TEXT",
		vaddr: s1.vsize,
		//vsize:      0x1000,
		vsize:      0x100 + uint64(len(code)),
		fileoffset: 0,
		// filesize:   0x1000,
		filesize: uint64(len(code)),
		prot1:    7,
		prot2:    5,
	}

	s2.sect = []MachoSect{
		{
			name:    "__text",
			segname: "__TEXT",
			//		addr:    s1.vaddr + 168, // Needs to be where the code starts...
			addr: s2.vaddr + 0x100,
			size: uint64(len(code)),
			off:  0x100,
			//		off:     loadsize + 168, // start of code
			flag: S_ATTR_PURE_INSTRUCTIONS | S_ATTR_SOME_INSTRUCTIONS,
		},
	}
	segs++
	sects++

	// loadsize += 18 * 4 * segs
	// loadsize += 20 * 4 * sects
	// size of segments + size of sections
	//loadsize += 72*2 + 80 // TODO: Include alginments

	// TODO: Clean up with alginment code down below
	loadsize := 0
	loadsize += 18 * 4 * segs
	loadsize += 20 * 4 * sects
	loadsize = 224

	// s2.sect[0].addr = s2.vaddr + uint64(loadsize)
	// s2.sect[0].off = uint32(loadsize)

	/*
		__mh_execute_header:
			dd 0xfeedfacf   ; MH_MAGIC_64
			dd 16777223     ; CPU_TYPE_X86 | CPU_ARCH_ABI64
			dd 0x80000003   ; CPU_SUBTYPE_I386_ALL | CPU_SUBTYPE_LIB64
			dd 2            ; MH_EXECUTE
			dd 16           ; number of load commands
			dd ___loadcmdsend - ___loadcmdsstart    ; size of load commands
			dd 0x00200085   ; MH_NOUNDEFS | MH_DYLDLINK | MH_TWOLEVEL | MH_PIE
			dd 0            ; reserved
		___loadcmdsstart:
	*/

	buf.write32(MH_MAGIC_64)      // Magic Number
	buf.write32(MACHO_CPU_AMD64)  // CPU Arch
	buf.write32(MACHO_SUBCPU_X86) // Subcpu?

	buf.write32(MH_EXECUTE) // Permissions

	buf.write32(uint32(segs))     // Number of load comamnds (includes segments)
	buf.write32(uint32(loadsize)) // Size of the load section

	buf.write32(uint32(MH_NOUNDEFS))
	buf.write32(0) //Reserved

	fmt.Println("start here", len(*buf))
	for _, s := range []*MachoSeg{s1, s2} {
		/*
			padding = (align - (offset mod align)) mod align
		*/

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

	/*
		align := 0x1000
		offset := len(*buf)
		padding := (align - (offset % align)) % align
		for i := 0; i < padding; i++ {
			buf.write(0x00)
		}
	*/

	fmt.Println("len header:", len(*buf))

	// Write out code
	buf.write(code...)

	if err := os.WriteFile("./mac_bin", buf.bytes(), 0755); err != nil {
		log.Fatal(err)
	}
}
