# Assemblor

- Experimental assembler for fun and learning
- X64 macho and ELF executable
- cross compile
- limited bytecode instruction set
- assemble and link in one command
- static binary

## Missing

- static data
- functions
- other syscalls?

## Future features

- register allocation
- debug information

## Assem Spec

- reg1-reg12

```
reserve buf: 32 * 1000
reserve tmp: 64 # int64
data_string hw: "Hello, world"
data_int n1: 123


mov_imm reg1, 0x1234	# mov_imm <reg_dst> <imm>
mov_reg reg2, reg1		# mov_reg <reg_dst> <reg_src>
branch some_label, reg1 == reg2
branch some_label, reg1 < reg2
branch some_label, reg1 > reg2
branch some_label, reg1 <= reg2
branch: some_label, reg1 >= reg2
jmp: some_other_label

add_imm: reg12, 0xdeadbeef
add_reg: reg11, reg12

write_imm: tmp, 1234567

mov_addr: reg1, buf
write_mem: reg1, 10
add reg1, 32
write_mem: reg1, 11

label some_label:
inc reg1
call some_func

func some_func:
push reg1
dec reg1
pop reg1
return

func print_hw:
pusha
mov_mem reg1, hw
syscall_write reg1, len(hw)
popa
return

label some_other_label:
mov_mem reg1, hw
syscall_write reg1, len(hw)

syscall_exit 0
```

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

