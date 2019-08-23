package codegen_test

import (
	"testing"
	"vm/codegen"
	"vm/core"
) 

func TestCodeGen(t *testing.T) {
	actual := codegen.CodeGen([]codegen.CodeGenPassInstruction{
		// set b 1
		codegen.SetWordInstruction{Destination: core.B, Operand1: 1},
		// core.Byte(core.SetWord), core.Byte(core.B), 1, 0, 0, 0, 0, 0, 0, 0,

		// loop:
		// print a
		codegen.PrintInstruction{Operand1: core.A},
		// core.Byte(core.Print), core.Byte(core.A),

		// add a a b
		codegen.ArithmeticInstruction{
			OpCode:      core.AddWord,
			Destination: core.A,
			Operand1:    core.A,
			Operand2:    core.B,
		},
		// core.Byte(core.AddWord), core.Byte(core.A), core.Byte(core.A), core.Byte(core.B),

		// swp a b
		//   set c 0
		//   set e 0
		//   store e b
		//   add b a c
		//   load e a
		codegen.SetWordInstruction{Destination: core.C, Operand1: 0},
		// core.Byte(core.SetWord), core.Byte(core.C), 0, 0, 0, 0, 0, 0, 0, 0,
		codegen.SetWordInstruction{Destination: core.E, Operand1: 0},
		// core.Byte(core.SetWord), core.Byte(core.E), 0, 0, 0, 0, 0, 0, 0, 0,
		codegen.StoreWordInstruction{Address: core.E, Operand1: core.B},
		// core.Byte(core.StoreWord), core.Byte(core.E), core.Byte(core.B),
		codegen.ArithmeticInstruction{
			OpCode:      core.AddWord,
			Destination: core.B,
			Operand1:    core.A,
			Operand2:    core.C,
		},
		// core.Byte(core.AddWord), core.Byte(core.B), core.Byte(core.A), core.Byte(core.C),
		codegen.LoadWordInstruction{Destination: core.A, Address: core.E},
		// core.Byte(core.LoadWord), core.Byte(core.A), core.Byte(core.E),

		// cmp c a 100
		//   set d 100
		//   cmp c a d
		codegen.SetWordInstruction{Destination: core.D, Operand1: 100},
		// core.Byte(core.SetWord), core.Byte(core.D), 100, 0, 0, 0, 0, 0, 0, 0,
		codegen.ArithmeticInstruction{
			OpCode:      core.Compare,
			Destination: core.C,
			Operand1:    core.A,
			Operand2:    core.D,
		},
		// core.Byte(core.Compare), core.Byte(core.C), core.Byte(core.A), core.Byte(core.D),

		// beq c -1 loop
		//   set d -1
		//   beq c d loop
		codegen.SetWordInstruction{Destination: core.D, Operand1: -1},
		// core.Byte(core.SetWord), core.Byte(core.D), 255, 255, 255, 255, 255, 255, 255, 255,
		codegen.BranchEqualInstruction{Operand1: core.C, Operand2: core.D, BranchTarget: 1},
		// core.Byte(core.BranchEqual), core.Byte(core.C), core.Byte(core.D), 10, 0, 0, 0, 0, 0, 0, 0,

		// exit 0
		codegen.ExitInstruction{Operand1: 0},
	})

	expected := core.Memory{
		core.Byte(core.SetWord), core.Byte(core.B), 1, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.Print), core.Byte(core.A),
		core.Byte(core.AddWord), core.Byte(core.A), core.Byte(core.A), core.Byte(core.B),
		core.Byte(core.SetWord), core.Byte(core.C), 0, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.SetWord), core.Byte(core.E), 0, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.StoreWord), core.Byte(core.E), core.Byte(core.B),
		core.Byte(core.AddWord), core.Byte(core.B), core.Byte(core.A), core.Byte(core.C),
		core.Byte(core.LoadWord), core.Byte(core.A), core.Byte(core.E),
		core.Byte(core.SetWord), core.Byte(core.D), 100, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.Compare), core.Byte(core.C), core.Byte(core.A), core.Byte(core.D),
		core.Byte(core.SetWord), core.Byte(core.D), 255, 255, 255, 255, 255, 255, 255, 255,
		core.Byte(core.BranchEqual), core.Byte(core.C), core.Byte(core.D), 10, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.Exit), core.Byte(0),
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("ERROR at %d, %d != %d", i, actual[i], expected[i])
		}
	}

	if len(expected) != len(actual) {
		t.Errorf("Mismatched lengths: %d != %d", len(expected), len(actual))
	}
}
