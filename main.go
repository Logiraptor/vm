package main

import (
	"vm/codegen"
	"vm/core"
)

func main() {
	machine := core.CreateVM(0, 8)
	machine.Instructions = codegen.CodeGen([]codegen.CodeGenPassInstruction{
		// set b 1
		codegen.SetWordInstruction{Destination: core.BRegister, Operand1: 1},

		// loop:
		// print a
		codegen.PrintInstruction{Operand1: core.ARegister},

		// add a a b
		codegen.ArithmeticInstruction{
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
		codegen.SetWordInstruction{Destination: core.CRegister, Operand1: 0},
		codegen.SetWordInstruction{Destination: core.ERegister, Operand1: 0},
		codegen.StoreWordInstruction{Address: core.ERegister, Operand1: core.BRegister},
		codegen.ArithmeticInstruction{
			OpCode:      core.AddWord,
			Destination: core.BRegister,
			Operand1:    core.ARegister,
			Operand2:    core.CRegister,
		},
		codegen.LoadWordInstruction{Destination: core.ARegister, Address: core.ERegister},

		// cmp c a 100
		//   set d 100
		//   cmp c a d
		codegen.SetWordInstruction{Destination: core.DRegister, Operand1: 100},
		codegen.ArithmeticInstruction{
			OpCode:      core.Compare,
			Destination: core.CRegister,
			Operand1:    core.ARegister,
			Operand2:    core.DRegister,
		},

		// beq c -1 loop
		//   set d -1
		//   beq c d loop
		codegen.SetWordInstruction{Destination: core.DRegister, Operand1: -1},
		codegen.BranchEqualInstruction{Operand1: core.CRegister, Operand2: core.DRegister, BranchTarget: 1},

		// exit 0
		codegen.ExitInstruction{Operand1: 0},
	})

	for {
		machine.RunInstruction()
	}
}
