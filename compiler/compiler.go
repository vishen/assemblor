package compiler

import (
	"fmt"
	"runtime"

	"github.com/vishen/assemblor/bytecode"
	"github.com/vishen/assemblor/ld"
	"github.com/vishen/assemblor/x64"
)

type TargetOS string

const (
	Macho TargetOS = "macho"
	Linux TargetOS = "linux"
)

type TargetArch string

const (
	X64 TargetArch = "x64"
)

func Compile(bc []bytecode.Instruction) ([]byte, error) {
	var arch TargetArch
	var os TargetOS

	switch runtime.GOARCH {
	case "amd64":
		arch = X64
	}

	switch runtime.GOOS {
	case "darwin":
		os = Macho
	}

	return CompileWithOptions(bc, os, arch)
}

func CompileWithOptions(bc []bytecode.Instruction, os TargetOS, arch TargetArch) ([]byte, error) {
	switch arch {
	case X64:
		// valid arch
	default:
		return nil, fmt.Errorf("unsupported architecture %q", arch)
	}

	var (
		linker ld.Linker
		carch  x64.Arch
	)

	switch os {
	case Macho:
		linker = ld.NewMacho()
		carch = x64.Macho
	case Linux:
		linker = ld.NewElf()
		carch = x64.Linux
	default:
		return nil, fmt.Errorf("unsupported os %q", os)
	}

	code, bssSize := x64.Compile(carch, bc, linker.BssAddr())
	executable := linker.Link(code, bssSize)
	return executable, nil
}
