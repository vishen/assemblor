# Brainfuck

Brainfuck implementation written in assemblor:

```
bss: 1000 x int64

c.x64.EmitMovRegImm(x64e.RAX, cells) // mov rax, cells ; current position in cells
c.x64.EmitMovRegImm(x64e.R15, 0)     // mov r15, 0 ; this is where the character to be outputted will be.
```

## Brianfuck spec

```
+: inc rax
-: dec rax
>: mov rax to next call
<: mov rax to prev cell
.: sys_write stdout rax
[: start loop
	- c.x64.EmitCmpMemImm(x64e.RAX, 0)
	- jeq end of ]
]: end loop
```

```golang
for _, ch := range c.program {
	switch ch {
	case '+':
		c.EmitInc()
	case '-':
		c.EmitDec()
	case '>':
		c.EmitNext()
	case '<':
		c.EmitPrev()
	case '.':
		c.EmitOutputChar()
	case '[':
		loopsCounter += 1
		c.EmitLoop()
	case ']':
		loopsFinished += 1
		c.EmitLoopJump()
	}
}
```




### Required? 

+ heap?: bss: 1000 x int64
+ cmp
+ jeq
- call
- ret

Can probably get away without call and return
