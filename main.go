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
	data             Memory
}

func CreateVM(stackSize, heapSize Word) VM {
	memory := AllocateMemory(stackSize + heapSize)
	return VM{
		data: memory,
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
	case LoadWord:
		dest := vm.ReadRegister()
		arg1 := vm.ReadRegister()
		*dest = vm.data.ReadWord(*arg1)
	case StoreWord:
		dest := vm.ReadRegister()
		arg1 := vm.ReadRegister()
		vm.data.WriteWord(*dest, *arg1)
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
		//   set e 0
		//   store e b
		//   add b a c
		//   load e a
		Byte(SetWord), Byte(CRegister), 0, 0, 0, 0, 0, 0, 0, 0,
		Byte(SetWord), Byte(ERegister), 0, 0, 0, 0, 0, 0, 0, 0,
		Byte(StoreWord), Byte(ERegister), Byte(BRegister),
		Byte(AddWord), Byte(BRegister), Byte(ARegister), Byte(CRegister),
		Byte(LoadWord), Byte(ARegister), Byte(ERegister),

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
