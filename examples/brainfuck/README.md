# Brainfuck

Brainfuck implementation written in `assemblor`.

## Brianfuck spec

```
char array[30000] = {0}

+: inc rax
-: dec rax
>: mov rax to next call
<: mov rax to prev cell
.: sys_write stdout rax
[: start loop
]: end loop
```

## Running

```
$ go run main.go -f examples/helloworld.bf
wrote executable to helloworld

$ file helloworld
helloworld: Mach-O 64-bit executable x86_64

./helloworld
Hello World!
```

## Cross compiling

```
$ go run main.go -f examples/helloworld.bf -os linux
wrote executable to helloworld

$ file helloworld
helloworld: ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, no section header

./helloworld
Hello World!
```

