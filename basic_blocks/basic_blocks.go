package basic_blocks

import (
	"fmt"
	"strconv"
	"vm/codegen"
	"vm/core"
)

//go:generate dsl_generator -type Instruction
type Instruction interface {
	isBasicBlockPassInstruction()
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
	BranchTarget *BasicBlock
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

type BasicBlock struct {
	id string
	InstructionBuilder
}

var basicBlockIndex = 0

func NewBasicBlock(name string) *BasicBlock {
	basicBlockIndex++
	return &BasicBlock{
		id: name + strconv.Itoa(basicBlockIndex),
	}
}

func ResolveLabels(basicBlocks []*BasicBlock) []codegen.Instruction {
	type branchTarget struct {
		branchLocation int
		basicBlock     *BasicBlock
	}
	var output []codegen.Instruction
	var branchOffsets []branchTarget
	var blockOffsets = map[*BasicBlock]int{}
	for _, bb := range basicBlocks {
		blockOffsets[bb] = len(output)
		for _, inst := range bb.Instructions {
			codegenInst, isBranch := convertInstruction(inst)
			if isBranch {
				branchOffsets = append(branchOffsets, branchTarget{
					branchLocation: len(output),
					basicBlock:     inst.(BranchEqualInstruction).BranchTarget,
				})
			}
			output = append(output, codegenInst)
		}
	}

	for _, bo := range branchOffsets {
		inst := output[bo.branchLocation].(codegen.BranchEqualInstruction)
		inst.BranchTarget = core.Word(blockOffsets[bo.basicBlock])
		output[bo.branchLocation] = inst
	}

	return output
}

func convertInstruction(inst Instruction) (instruction codegen.Instruction, isBranch bool) {
	switch v := inst.(type) {
	case BranchEqualInstruction:
		return codegen.BranchEqualInstruction{
			Operand1:     v.Operand1,
			Operand2:     v.Operand2,
			BranchTarget: 0,
		}, true
	case ArithmeticInstruction:
		return codegen.ArithmeticInstruction{
			OpCode:      v.OpCode,
			Destination: v.Destination,
			Operand1:    v.Operand1,
			Operand2:    v.Operand2,
		}, false
	case SetWordInstruction:
		return codegen.SetWordInstruction{
			Destination: v.Destination,
			Operand1:    v.Operand1,
		}, false
	case LoadWordInstruction:
		return codegen.LoadWordInstruction{
			Destination: v.Destination,
			Address:     v.Address,
		}, false
	case StoreWordInstruction:
		return codegen.StoreWordInstruction{
			Address:  v.Address,
			Operand1: v.Operand1,
		}, false
	case PrintInstruction:
		return codegen.PrintInstruction{
			Operand1: v.Operand1,
		}, false
	case ExitInstruction:
		return codegen.ExitInstruction{
			Operand1: v.Operand1,
		}, false
	}

	panic(fmt.Sprintf("Unsupported instruction in basic_block#convertInstruction: %T", inst))
}
