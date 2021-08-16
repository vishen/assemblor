; A minimal Mach-o x64 executable for OS X Sierra
; $ nasm -f bin -o tiny_hello tiny_hello.s
; $ chmod +x tiny_hello
; Constants (For readability)
%define MH_MAGIC_64                    0xfeedfacf
%define CPU_ARCH_ABI64                0x01000000
%define    CPU_TYPE_I386                0x00000007
%define CPU_TYPE_X86_64                CPU_ARCH_ABI64 | CPU_TYPE_I386
%define CPU_SUBTYPE_LIB64            0x80000000
%define CPU_SUBTYPE_I386_ALL        0x00000003
%define MH_EXECUTE                    0x2
%define MH_NOUNDEFS                    0x1
%define LC_SEGMENT_64                0x19
%define LC_UNIXTHREAD                0x5 
%define VM_PROT_READ                0x1
%define VM_PROT_WRITE                0x2
%define VM_PROT_EXECUTE                0x4
%define x86_THREAD_STATE64            0x4
%define    x86_EXCEPTION_STATE64_COUNT    42
%define SYSCALL_CLASS_SHIFT            24
%define SYSCALL_CLASS_MASK            (0xFF << SYSCALL_CLASS_SHIFT)
%define SYSCALL_NUMBER_MASK            (~SYSCALL_CLASS_MASK)  
%define SYSCALL_CLASS_UNIX            2
%define SYSCALL_CONSTRUCT_UNIX(syscall_number) \
            ((SYSCALL_CLASS_UNIX << SYSCALL_CLASS_SHIFT) | \
             (SYSCALL_NUMBER_MASK & (syscall_number)))
%define SYS_exit                    1
%define SYS_write                    4
; NASM directive, not compiled
; Use RIP-Relative addressing for x64
BITS    64
;DEFAULT    REL
%define __origin 0x100000000
org __origin
; Mach-O header
DD        MH_MAGIC_64                                        ; magic
DD        CPU_TYPE_X86_64                                    ; cputype
DD        CPU_SUBTYPE_LIB64 | CPU_SUBTYPE_I386_ALL        ; cpusubtype
DD        MH_EXECUTE                                        ; filetype
DD        3                                                ; ncmds
DD        __COMMANDSend  - __COMMANDSstart                ; sizeofcmds
DD        MH_NOUNDEFS                                        ; flags
DD        0x0                                                ; reserved
__COMMANDSstart:

___PAGEZEROstart:
        DD        LC_SEGMENT_64                                    ; cmd
        dd         ___PAGEZEROend - ___PAGEZEROstart                ; command size
hello_str:
        db         '__PAGEZERO',0x0,0,0,0,0,0 ; segment name (pad to 16 bytes)
        DQ        0x0                                                ; vmaddr
        DQ        __origin                                        ; vmsize
        DQ        0                                                ; fileoff
        DQ        0                                                ; filesize
        DD        0                                                 ; maxprot
        DD        0                                                ; initprot
        DD        0x0                                                ; nsects
        DD        0x0                                                ; flags
___PAGEZEROend:
; Segment and Sections
___TEXTstart:
        DD        LC_SEGMENT_64                                    ; cmd
        dd ___TEXTend - ___TEXTstart    ; command size

        db '__TEXT',0,0,0,0,0,0,0,0,0,0 ; segment name (pad to 16 bytes)
        DQ        __origin                                        ; vmaddr
        DQ        ___codeend - __origin                ; vmsize
        DQ        0                                                ; fileoff
        DQ        ___codeend - __origin                    ; filesize
        DD        VM_PROT_READ | VM_PROT_WRITE | VM_PROT_EXECUTE    ; maxprot
        DD        VM_PROT_READ | VM_PROT_EXECUTE                            ; initprot
        DD        0x0                                                ; nsects
        DD        0x0                                                ; flags
___TEXTend:
__UNIX_THREADstart:
; UNIX Thread Status
DD        LC_UNIXTHREAD                                    ; cmd
DD        __UNIX_THREADend - __UNIX_THREADstart             ; cmdsize
DD        x86_THREAD_STATE64                                ; flavor
DD        x86_EXCEPTION_STATE64_COUNT                        ; count
DQ        0x0, 0x0, 0x00, 0x0                                ; rax, rbx , rcx , rdx
DQ        0x01, hello_str, 0x00, 0x00                        ; rdi = STDOUT, rsi = address of hello_str,  rbp, rsp
DQ        0x00, 0x00                                        ; r8 and r9
DQ        0x00, 0x00, 0x00, 0x00, 0x00, 0x00                ; r10, r11, r12, r13, r14, r15
DQ         ___codestart, 0x00, 0x00, 0x00, 0x00            ; rip, rflags, cs, fs, gs
__UNIX_THREADend:
__COMMANDSend:
___codestart:                                                    ; 24 bytes
    ; rdi and rsi have already been set in the initial state
    mov        rdx, 11
    mov        rax, SYSCALL_CONSTRUCT_UNIX(SYS_write)
    syscall
    mov            rdi, rax
    mov            rax, SYSCALL_CONSTRUCT_UNIX(SYS_exit)
    syscall
___codeend:
    times 4096-($-$$) DB  0;
    filesize    EQU    $-$$
