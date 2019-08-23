package core

import (
	"fmt"
	"vm/debug"
)

//go:generate stringer -type=OpCode
type OpCode Byte

const (
	// AddWord | DESTINATION | OPERAND 1 | OPERAND 2
	AddWord OpCode = iota
	// Compare | DESTINATION | OPERAND 1 | OPERAND 2
	Compare
	// SetWord | Destination | LiteralWord
	SetWord
	// BranchEqual | OPERAND 1 | OPERAND 2 | LiteralWord
	BranchEqual
	// LoadWord | DESTINATION | OPERAND 1 (Address)
	LoadWord
	// StoreWord | DESTINATION (Address) | OPERAND 1 (Value)
	StoreWord
	// Print | OPERAND
	Print
	// Exit | ExitCode
	Exit
	ExtendedInstruction OpCode = 255
)

//go:generate stringer -type=Register
type Register Byte

const (
	IP Register = iota
	SP
	A
	B
	C
	D
	E
	F
)

func (vm *VM) ReadOpCode() OpCode {
	opCode := OpCode(vm.Instructions[vm.ip])
	vm.ip++
	debug.Debugln("ReadOpCode", opCode)
	return opCode
}

func (vm *VM) ReadRegister() *Word {
	register := Register(vm.Instructions[vm.ip])
	vm.ip++
	debug.Debugln("ReadRegister", register)
	switch register {
	case SP:
		return &vm.sp
	case IP:
		return &vm.ip
	case A:
		return &vm.a
	case B:
		return &vm.b
	case C:
		return &vm.c
	case D:
		return &vm.d
	case E:
		return &vm.e
	case F:
		return &vm.f
	}
	panic(fmt.Sprintf("InvalidRegister: %d", register))
}

func (vm *VM) ReadWordLiteral() Word {
	val := vm.Instructions.ReadWord(vm.ip)
	vm.ip += WordSize
	debug.Debugln("ReadWordLiteral", val)
	return val
}
