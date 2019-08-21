package core

import (
	"fmt"
	"os"
	"vm/debug"
)

const WordSize = 8
 
type VM struct {
	sp, ip           Word
	a, b, c, d, e, f Word
	Instructions     Memory
	data             Memory
}

func CreateVM(stackSize, heapSize Word) VM {
	memory := AllocateMemory(stackSize + heapSize)
	return VM{
		data: memory,
	}
}

func (vm *VM) RunInstruction() {
	debug.Debugf("VMState ip=%d, sp=%d, a=%d, b=%d, c=%d\n", vm.ip, vm.sp, vm.a, vm.b, vm.c)
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
		arg1 := vm.Instructions[vm.ip]
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
