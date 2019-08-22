package main

import (
	"vm/basic_blocks"
	"vm/codegen"
	"vm/core"
)

func main() {
	bb1 := basic_blocks.NewBasicBlock("main")
	bb1.Instructions = []basic_blocks.Instruction{
		// set b 1
		basic_blocks.SetWordInstruction{Destination: core.BRegister, Operand1: 1},
	}
	bb2 := basic_blocks.NewBasicBlock("loop")
	bb2.Instructions = []basic_blocks.Instruction{
		// loop:
		// print a
		basic_blocks.PrintInstruction{Operand1: core.ARegister},

		// add a a b
		basic_blocks.ArithmeticInstruction{
			OpCode:      core.AddWord,
			Destination: core.ARegister,
			Operand1:    core.ARegister,
			Operand2:    core.BRegister,
		},

		// swp a b
		//   set c 0
		//   set e 0
		//   store e b
		//   add b a c
		//   load e a
		basic_blocks.SetWordInstruction{Destination: core.CRegister, Operand1: 0},
		basic_blocks.SetWordInstruction{Destination: core.ERegister, Operand1: 0},
		basic_blocks.StoreWordInstruction{Address: core.ERegister, Operand1: core.BRegister},
		basic_blocks.ArithmeticInstruction{
			OpCode:      core.AddWord,
			Destination: core.BRegister,
			Operand1:    core.ARegister,
			Operand2:    core.CRegister,
		},
		basic_blocks.LoadWordInstruction{Destination: core.ARegister, Address: core.ERegister},

		// cmp c a 100
		//   set d 100
		//   cmp c a d
		basic_blocks.SetWordInstruction{Destination: core.DRegister, Operand1: 100},
		basic_blocks.ArithmeticInstruction{
			OpCode:      core.Compare,
			Destination: core.CRegister,
			Operand1:    core.ARegister,
			Operand2:    core.DRegister,
		},

		// beq c -1 loop
		//   set d -1
		//   beq c d loop
		basic_blocks.SetWordInstruction{Destination: core.DRegister, Operand1: -1},
		basic_blocks.BranchEqualInstruction{Operand1: core.CRegister, Operand2: core.DRegister, BranchTarget: bb2},

		// exit 0
		basic_blocks.ExitInstruction{Operand1: 0},
	}
	machine := core.CreateVM(0, 8)
	machine.Instructions = codegen.CodeGen(basic_blocks.ResolveLabels([]*basic_blocks.BasicBlock{
		bb1, bb2,
	}))

	for {
		machine.RunInstruction()
	}
}
