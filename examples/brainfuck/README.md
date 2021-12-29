# Brainfuck

Brainfuck implementation written in `assemblor`:

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
