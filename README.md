# Assemblor

- Experimental for fun and learning
- X64 macho and ELF executable
- cross compile
- limited bytecode instruction set
- assemble and link in one command
- static binary


NOTE: This is still a heavy work in progress and isn't intended to be used for anything 
other than learning. This currently only provides a basic programmatic API.

## Missing

- static data
- functions
- other syscalls?

## Future features

- register allocation
- debug information

## Example

	# Running without specifying OS will result in the program being compiled for the current OS
	$ go run examples/hello_world/main.go
	2022/07/13 09:10:35 Compiling for linux_x64
	2022/07/13 09:10:35 Writing executable to assemblored
	$ ./assemblored
	Hello World
	$ file ./assemblored
	assemblored: ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, no section header

	# You can cross-compile by specifying the os flag
	$ go run examples/hello_world/main.go -os macho -o assemblored_macho
	2022/07/13 09:13:07 Compiling for macho_x64
	2022/07/13 09:13:07 Writing executable to assemblored_macho
	$ file ./assemblored_macho
	assemblored_macho: Mach-O 64-bit x86_64 executable, flags:<NOUNDEFS>


## Resources

### Syscalls

- https://en.wikipedia.org/wiki/X86_calling_conventions#System_V_AMD64_ABI

### X64

- https://defuse.ca/online-x86-assembler.htm#disassembly
- https://uica.uops.info/
- https://www.felixcloutier.com/x86/

### Linux

- https://blog.rchapman.org/posts/Linux_System_Call_Table_for_x86_64/

### Macho

- https://stackoverflow.com/questions/32453849/minimal-mach-o-64-binary/32659692#32659692
- https://stackoverflow.com/questions/39863112/what-is-required-for-a-mach-o-executable-to-load
- https://redmaple.tech/blogs/macho-files/
- https://golang.org/src/cmd/link/internal/ld/macho.go
- https://www.mikeash.com/pyblog/friday-qa-2012-11-30-lets-build-a-mach-o-executable.html
- https://gist.github.com/zliuva/1084476
- https://github.com/aidansteele/osx-abi-macho-file-format-reference
- https://opensource.apple.com/source/xnu/xnu-1504.7.4/osfmk/mach/i386/thread_status.h

### lldb

- https://lldb.llvm.org/use/map.html

### Other

- https://docs.microsoft.com/en-us/cpp/build/x64-software-conventions?view=msvc-170#register-usage

