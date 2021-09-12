global  _main
  section  .text
_main:
  ;mov rax, qword [8589934592]
  ;mov r8, [qword 0xdeadbe]
  ;mov rax, ds:0x04
  ;mov rax, qword buffer
  mov rax, [qword buffer]
  ;mov rcx, [qword 0xdeadbe]
  ;mov qword [0xdeadbe], 0x1234
  mov    rbx, [qword 0x02000004]    ; system call for write
  mov    rax, 0x02000004    ; system call for write
  mov    rdi, 1             ; file descriptor 1 is stdout
  mov    rsi, qword message ; get string address
  mov    rdx, 13            ; number of bytes
  syscall                   ; execute syscall (write)
  mov    rax, 0x02000001    ; system call for exit
  mov    rdi, 0             ; exit code 0
  syscall                   ; execute syscall (exit)
message: db    "Hello, World!", 0Ah, 00h

  section .bss
buffer: resb 6
