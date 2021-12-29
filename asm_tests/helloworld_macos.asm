global  _main
  section  .text
_main:
  push rax
  pop rax
  push rsi
  pop rsi
  push r10
  pop r10
  MOV RCX, [RAX]
  MOV R12, qword [RAX]
  MOV R12, qword [RAX]
  mov [rbx], rcx
  mov qword [rbx], rcx
  mov qword [0xdead], rbx
  mov rax, 0
  mov rdx, 5
l1:
  cmp rax, rdx
  je exit
  inc rax
  jmp l1
write:
  mov    rax, 0x02000004    ; system call for write
  mov    rdi, 1             ; file descriptor 1 is stdout
  mov    rsi, qword message ; get string address
  mov    rdx, 13            ; number of bytes
  syscall                   ; execute syscall (write)
exit:
  mov    rdi, rax             ; exit code 0
  mov    rax, 0x02000001    ; system call for exit
  syscall                   ; execute syscall (exit)
message: db    "Hello, World!", 0Ah, 00h

  section .bss
buffer: resb 6
