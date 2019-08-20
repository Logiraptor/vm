package main

import "fmt"

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
	IPRegister Register = iota
	SPRegister
	ARegister
	BRegister
	CRegister
	DRegister
	ERegister
	FRegister
)

func (vm *VM) ReadOpCode() OpCode {
	opCode := OpCode(vm.instructions[vm.ip])
	vm.ip++
	Debugln("ReadOpCode", opCode)
	return opCode
}

func (vm *VM) ReadRegister() *Word {
	register := Register(vm.instructions[vm.ip])
	vm.ip++
	Debugln("ReadRegister", register)
	switch register {
	case SPRegister:
		return &vm.sp
	case IPRegister:
		return &vm.ip
	case ARegister:
		return &vm.a
	case BRegister:
		return &vm.b
	case CRegister:
		return &vm.c
	case DRegister:
		return &vm.d
	case ERegister:
		return &vm.e
	case FRegister:
		return &vm.f
	}
	panic(fmt.Sprintf("InvalidRegister: %d", register))
}


func (vm *VM) ReadWordLiteral() Word {
	val := vm.instructions.ReadWord(vm.ip)
	vm.ip += WordSize
	Debugln("ReadWordLiteral", val)
	return val
}
