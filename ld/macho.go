package ld

const (
	// VM start address
	machoOrigin uint64 = 0x100000000

	// Minimum code size for a mach-o executable
	minCodeSize int = 0x1000

	toFill uint64 = 0xdeadbeef
)

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

type machoSeg struct {
	name       string
	vsize      uint64
	vaddr      uint64
	fileoffset uint64
	filesize   uint64
	prot1      uint32
	prot2      uint32
	// nsect      uint32
	// msect      uint32
	flag uint32
}

// Macho will take machine code and turn it into an executable
// macho file
//
// Currently it only creates and x64 executable.
func Macho(code []byte) []byte {
	// TODO: Will need to do __DATA

	// Mach-o required at the code data to be at least 4096 bytes otherwise
	// it is unable to execute it. Ensure the code size if larger than
	// the required size, otherwise set it to the minimum required size
	// and add some padding down below
	codeSize := len(code)
	codePadding := 0
	if codeSize < minCodeSize {
		codePadding = int(minCodeSize) - codeSize
		codeSize = minCodeSize
	}

	// Offsets into the output where data needs to be back-patched
	var loadsizeOffset int = 0
	var codestartOffset int = 0

	// __PAGEZERO segment is used for null pointer deferences
	s1 := &machoSeg{
		name:  "__PAGEZERO",
		vsize: machoOrigin,
	}

	// __TEXT segment contains the setup and offsets to the machine code
	s2 := &machoSeg{
		name:       "__TEXT",
		vaddr:      s1.vsize,
		vsize:      uint64(codeSize),
		fileoffset: 0,
		filesize:   uint64(codeSize),
		prot1:      7,
		prot2:      5,
	}

	// Start writing mach-o binary data

	// Mach-o Headers
	buf.write32(MH_MAGIC_64)      // Magic Number
	buf.write32(MACHO_CPU_AMD64)  // CPU Arch
	buf.write32(MACHO_SUBCPU_X86) // Subcpu?

	buf.write32(MH_EXECUTE) // Permissions

	buf.write32(3) // Number of segments; __PAGEZERO, __TEXT, UNIXTHREAD

	loadsizeOffset = buf.offset()
	buf.write32(uint32(toFill)) // Size of the load section

	buf.write32(uint32(MH_NOUNDEFS)) // No unresolved symbols
	buf.write32(0)                   //Reserved

	loadsizeStart := buf.offset()
	for _, s := range []*machoSeg{s1, s2} {
		buf.write32(LC_SEGMENT_64)
		buf.write32(72) // Size if always 72
		buf.writeString(s.name, 16)
		buf.write64(s.vaddr)
		buf.write64(s.vsize)
		buf.write64(s.fileoffset)
		buf.write64(s.filesize)
		buf.write32(s.prot1)
		buf.write32(s.prot2)
		buf.write32(0)
		buf.write32(s.flag)
	}

	// Need to tell mac how to find and start the machine code
	buf.write32(LC_UNIXTHREAD)
	buf.write32(184) // command size
	buf.write32(4)   // thread state: x86_THREAD_STATE64
	buf.write32(42)  // word count: x86_EXCEPTION_STATE64_COUNT

	// Setup the initial values for registers
	buf.write64(0x00, 0x00, 0x00, 0x00) // rax, rbx , rcx , rdx
	buf.write64(0x00, 0x00, 0x00, 0x00) // rdi, rsi, rbp, rsp
	buf.write64(0x00, 0x00, 0x00, 0x00) // r8, r9, r10, r11
	buf.write64(0x00, 0x00, 0x00, 0x00) // r12, r13, r14, r15

	codestartOffset = buf.offset()
	buf.write64(toFill) // rip=codestart

	buf.write64(0x00, 0x00, 0x00, 0x00) // rflags, cs, fs, gs

	// Fill in the loadsize offset once we know the size of the
	// load commands
	buf.fill32(loadsizeOffset, uint32(buf.offset()-loadsizeStart))

	// Fill in the codestart offset when we know where the code
	// starts
	buf.fill64(codestartOffset, machoOrigin+uint64(buf.offset()))

	// Write out code
	buf.write(code...)

	// Pad out the code section if under the required minimum
	// size
	for i := 0; i < codePadding; i++ {
		buf.write(0x00)
	}
	return buf.bytes()
}
