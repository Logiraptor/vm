package codegen

import (
	"reflect"
	"unsafe"
	"vm/core"
	"vm/debug"
)

type Instruction interface {
	isCodeGenPassInstruction()
}

type ArithmeticInstruction struct {
	Instruction
	OpCode      core.OpCode
	Destination core.Register
	Operand1    core.Register
	Operand2    core.Register
}

type SetWordInstruction struct {
	Instruction
	Destination core.Register
	Operand1    core.Word
}

type BranchEqualInstruction struct {
	Instruction
	Operand1     core.Register
	Operand2     core.Register
	BranchTarget core.Word
}

type LoadWordInstruction struct {
	Instruction
	Destination core.Register
	Address     core.Register
}

type StoreWordInstruction struct {
	Instruction
	Address  core.Register
	Operand1 core.Register
}

type PrintInstruction struct {
	Instruction
	Operand1 core.Register
}

type ExitInstruction struct {
	Instruction
	Operand1 core.Byte
}

func wordToMem(word core.Word) core.Memory {
	mem := *(*core.Memory)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&word)),
		Len:  core.WordSize,
		Cap:  core.WordSize,
	}))
	debug.Debugln(word, "==", mem)
	return mem
}

func CodeGen(instructions []Instruction) core.Memory {
	type branchOffset struct{ branch, instruction int }
	var output core.Memory
	var instructionOffsets []int
	var branchOffsets []branchOffset

	for _, instruction := range instructions {
		instructionOffsets = append(instructionOffsets, len(output))
		switch v := instruction.(type) {
		case ArithmeticInstruction:
			output = append(output,
				core.Byte(v.OpCode),
				core.Byte(v.Destination),
				core.Byte(v.Operand1),
				core.Byte(v.Operand2),
			)
		case SetWordInstruction:
			output = append(append(output,
				core.Byte(core.SetWord),
				core.Byte(v.Destination),
			),
				wordToMem(v.Operand1)...,
			)
		case BranchEqualInstruction:
			output = append(append(output,
				core.Byte(core.BranchEqual),
				core.Byte(v.Operand1),
				core.Byte(v.Operand2),
			),
				// write a 0 here for now since we don't know where the
				// instruction is yet.
				wordToMem(0)...,
			)
			branchOffsets = append(branchOffsets, branchOffset{
				branch:      len(output) - core.WordSize,
				instruction: int(v.BranchTarget),
			})
		case LoadWordInstruction:
			output = append(output,
				core.Byte(core.LoadWord),
				core.Byte(v.Destination),
				core.Byte(v.Address),
			)
		case StoreWordInstruction:
			output = append(output,
				core.Byte(core.StoreWord),
				core.Byte(v.Address),
				core.Byte(v.Operand1),
			)
		case PrintInstruction:
			output = append(output,
				core.Byte(core.Print),
				core.Byte(v.Operand1),
			)
		case ExitInstruction:
			output = append(output,
				core.Byte(core.Exit),
				v.Operand1,
			)
		}
	}

	for _, bo := range branchOffsets {
		branchTarget := instructionOffsets[bo.instruction]
		debug.Debugln(bo.branch, bo.instruction, branchTarget)
		copy(output[bo.branch:], wordToMem(core.Word(branchTarget)))
	}

	return output
}
