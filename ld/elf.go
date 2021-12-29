package ld

const (
	virtualStartAddr    uint64 = 0x400000
	bssVirtualStartAddr uint64 = 0x600000
	alignment           uint64 = 0x200000
)

type Elf struct{}

func NewElf() Elf {
	return Elf{}
}

func (e Elf) BssAddr() uint64 {
	return bssVirtualStartAddr
}

func (e Elf) Link(code []byte, bssSize uint64) []byte {
	textSize := uint64(len(code))
	// Size of ELF header + 2 * size program header. The size of
	// the ELF header is always 0x40 bytes, and the size of each
	// program header is always 0x38 bytes.
	textOffset := uint64(0x40 + (2 * 0x38))

	// Build ELF Header
	buf.write(0x7f, 0x45, 0x4c, 0x46) // ELF magic value

	buf.write(0x02) // 64-bit executable
	buf.write(0x01) // Little endian
	buf.write(0x01) // ELF version
	buf.write(0x03) // Target OS ABI
	buf.write(0x00) // Further specify ABI version

	buf.write(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Unused bytes

	buf.write(0x02, 0x00)             // Executable type
	buf.write(0x3e, 0x00)             // x86-64 target architecture
	buf.write(0x01, 0x00, 0x00, 0x00) // ELF version

	// 64-bit virtual offsets always start at 0x400000?? https://stackoverflow.com/questions/38549972/why-elf-executables-have-a-fixed-load-address
	// This seems to be a convention set in the x86_64 system-v abi: https://refspecs.linuxfoundation.org/elf/x86_64-SysV-psABI.pdf P26
	buf.write64(virtualStartAddr + textOffset)

	buf.write(0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Offset from file to program header
	buf.write(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) // Start of section header table
	buf.write(0x00, 0x00, 0x00, 0x00)                         // Flags
	buf.write(0x40, 0x00)                                     // Size of this header
	buf.write(0x38, 0x00)                                     // Size of a program header table entry - This should always be the same for 64-bit
	buf.write(0x02, 0x00)                                     // Length of sections: data and text for now
	buf.write(0x00, 0x00)                                     // Size of section header, which we aren't using
	buf.write(0x00, 0x00)                                     // Number of entries section header
	buf.write(0x00, 0x00)                                     // Index of section header table entry

	// Build Program Header
	// Text Segment
	buf.write(0x01, 0x00, 0x00, 0x00) // PT_LOAD, loadable segment. Both data and text segment use this.
	buf.write(0x05, 0x00, 0x00, 0x00) // Flags: 0x4 executable, 0x2 write, 0x1 read
	buf.write64(0)                    // textOffset
	buf.write64(virtualStartAddr)     // Offset from the beginning of the file. These values depend on how big the header and segment sizes are.
	buf.write64(virtualStartAddr)     // Physical address, irrelavnt on linux.
	buf.write64(textSize)             // Number of bytes in file image of segment, must be larger than or equal to the size of payload in segment. Should be zero for bss data.
	buf.write64(textSize)             // Number of bytes in memory image of segment, is not always same size as file image.
	buf.write64(alignment)

	// Build Program Header
	// Bss Segment
	buf.write(0x01, 0x00, 0x00, 0x00) // PT_LOAD, loadable segment. Both data and text segment use this.
	buf.write(0x07, 0x00, 0x00, 0x00) // Flags: 0x4 executable, 0x2 write, 0x1 read // TODO: Which flags are which values exactly???
	buf.write64(0)                    // Offset address.
	buf.write64(bssVirtualStartAddr)  // Virtual address.
	buf.write64(bssVirtualStartAddr)  // Physical address.
	buf.write64(0)                    // Number of bytes in file image.
	buf.write64(bssSize)              // Number of bytes in memory image.
	buf.write64(alignment)

	// Output the text segment
	buf.write(code...)

	return buf.bytes()
}
