global  _main
  section  .text
_main:
label2:
  jmp label1
  mov rax, 0xdeadbeef
label1:
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
jmp label1
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadb
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdead
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadb
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdead
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadb
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdead
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
jmp label2
label3:
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov rax, 0xdeadbeef
  mov    rax, 0x02000004    ; system call for write
  mov    rdi, 1             ; file descriptor 1 is stdout
  mov    rsi, qword message ; get string address
  mov    rdx, 13            ; number of bytes
  syscall                   ; execute syscall (write)
  mov    rax, 0x02000001    ; system call for exit
  mov    rdi, 0             ; exit code 0
  syscall                   ; execute syscall (exit)
message: db    "Hello, World!", 0Ah, 00h
