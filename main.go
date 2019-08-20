package main

import (
	"fmt"
	"os"
)

const Debug = false
const WordSize = 8

func Debugln(args ...interface{}) {
	if Debug {
		fmt.Fprintln(os.Stdout, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if Debug {
		fmt.Fprintf(os.Stdout, format, args...)
	}
}

type VM struct {
	sp, ip           Word
	a, b, c, d, e, f Word
	instructions     Memory
	stack            Memory
	heap             Memory
}

func CreateVM(stackSize, heapSize Word) VM {
	memory := AllocateMemory(stackSize + heapSize)
	return VM{
		stack: memory[:stackSize],
		heap:  memory[stackSize:],
	}
}

func (vm *VM) RunInstruction() {
	Debugf("VMState ip=%d, sp=%d, a=%d, b=%d, c=%d\n", vm.ip, vm.sp, vm.a, vm.b, vm.c)
	opCode := vm.ReadOpCode()
	switch opCode {
	case AddWord:
		dest := vm.ReadRegister()
		arg1 := vm.ReadRegister()
		arg2 := vm.ReadRegister()
		*dest = *arg1 + *arg2
	case Print:
		arg1 := vm.ReadRegister()
		fmt.Println(*arg1)
	case SetWord:
		dest := vm.ReadRegister()
		arg1 := vm.ReadWordLiteral()
		*dest = arg1
	case Exit:
		arg1 := vm.instructions[vm.ip]
		vm.ip++
		os.Exit(int(arg1))
	case Compare:
		dest := vm.ReadRegister()
		arg1 := vm.ReadRegister()
		arg2 := vm.ReadRegister()
		if *arg1 > *arg2 {
			*dest = 1
		} else if *arg1 < *arg2 {
			*dest = -1
		} else {
			*dest = 0
		}
	case BranchEqual:
		arg1 := vm.ReadRegister()
		arg2 := vm.ReadRegister()
		newIP := vm.ReadWordLiteral()
		if *arg1 == *arg2 {
			vm.ip = newIP
		}
	}
}

func main() {
	vm := CreateVM(10000, 10000)
	vm.instructions = Memory{
		// set b 1
		Byte(SetWord), Byte(BRegister), 1, 0, 0, 0, 0, 0, 0, 0,

		// loop:
		// print a
		Byte(Print), Byte(ARegister),

		// add a a b
		Byte(AddWord), Byte(ARegister), Byte(ARegister), Byte(BRegister),

		// swp a b
		//   set c 0
		//   add e b c
		//   add b a c
		//   add a e c
		Byte(SetWord), Byte(CRegister), 0, 0, 0, 0, 0, 0, 0, 0,
		Byte(AddWord), Byte(ERegister), Byte(BRegister), Byte(CRegister),
		Byte(AddWord), Byte(BRegister), Byte(ARegister), Byte(CRegister),
		Byte(AddWord), Byte(ARegister), Byte(ERegister), Byte(CRegister),

		// cmp c a 100
		//   set d 100
		//   cmp c a d
		Byte(SetWord), Byte(DRegister), 100, 0, 0, 0, 0, 0, 0, 0,
		Byte(Compare), Byte(CRegister), Byte(ARegister), Byte(DRegister),

		// beq c -1 loop
		//   set d -1
		//   beq c d loop
		Byte(SetWord), Byte(DRegister), 255, 255, 255, 255, 255, 255, 255, 255,
		Byte(BranchEqual), Byte(CRegister), Byte(DRegister), 10, 0, 0, 0, 0, 0, 0, 0,

		// exit 0
		Byte(Exit), Byte(0),
	}
	for {
		vm.RunInstruction()
	}
}
