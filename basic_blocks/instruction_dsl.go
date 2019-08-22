package basic_blocks

import "vm/core"

type InstructionBuilder struct {
	Instructions []Instruction
}

func (b *InstructionBuilder) Arithmetic(opcode core.OpCode, destination core.Register, operand1 core.Register, operand2 core.Register) {
	b.Instructions = append(b.Instructions, ArithmeticInstruction{
		OpCode:      opcode,
		Destination: destination,
		Operand1:    operand1,
		Operand2:    operand2,
	})
}

func (b *InstructionBuilder) SetWord(destination core.Register, operand1 core.Word) {
	b.Instructions = append(b.Instructions, SetWordInstruction{
		Destination: destination,
		Operand1:    operand1,
	})
}

func (b *InstructionBuilder) BranchEqual(operand1 core.Register, operand2 core.Register, branchtarget *BasicBlock) {
	b.Instructions = append(b.Instructions, BranchEqualInstruction{
		Operand1:     operand1,
		Operand2:     operand2,
		BranchTarget: branchtarget,
	})
}

func (b *InstructionBuilder) LoadWord(destination core.Register, address core.Register) {
	b.Instructions = append(b.Instructions, LoadWordInstruction{
		Destination: destination,
		Address:     address,
	})
}

func (b *InstructionBuilder) StoreWord(address core.Register, operand1 core.Register) {
	b.Instructions = append(b.Instructions, StoreWordInstruction{
		Address:  address,
		Operand1: operand1,
	})
}

func (b *InstructionBuilder) Print(operand1 core.Register) {
	b.Instructions = append(b.Instructions, PrintInstruction{
		Operand1: operand1,
	})
}

func (b *InstructionBuilder) Exit(operand1 core.Byte) {
	b.Instructions = append(b.Instructions, ExitInstruction{
		Operand1: operand1,
	})
}
