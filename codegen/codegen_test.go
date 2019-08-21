package codegen_test

import (
	"testing"
	"vm/codegen"
	"vm/core"
) 

func TestCodeGen(t *testing.T) {
	actual := codegen.CodeGen([]codegen.CodeGenPassInstruction{
		// set b 1
		codegen.SetWordInstruction{Destination: core.BRegister, Operand1: 1},
		// core.Byte(core.SetWord), core.Byte(core.BRegister), 1, 0, 0, 0, 0, 0, 0, 0,

		// loop:
		// print a
		codegen.PrintInstruction{Operand1: core.ARegister},
		// core.Byte(core.Print), core.Byte(core.ARegister),

		// add a a b
		codegen.ArithmeticInstruction{
			OpCode:      core.AddWord,
			Destination: core.ARegister,
			Operand1:    core.ARegister,
			Operand2:    core.BRegister,
		},
		// core.Byte(core.AddWord), core.Byte(core.ARegister), core.Byte(core.ARegister), core.Byte(core.BRegister),

		// swp a b
		//   set c 0
		//   set e 0
		//   store e b
		//   add b a c
		//   load e a
		codegen.SetWordInstruction{Destination: core.CRegister, Operand1: 0},
		// core.Byte(core.SetWord), core.Byte(core.CRegister), 0, 0, 0, 0, 0, 0, 0, 0,
		codegen.SetWordInstruction{Destination: core.ERegister, Operand1: 0},
		// core.Byte(core.SetWord), core.Byte(core.ERegister), 0, 0, 0, 0, 0, 0, 0, 0,
		codegen.StoreWordInstruction{Address: core.ERegister, Operand1: core.BRegister},
		// core.Byte(core.StoreWord), core.Byte(core.ERegister), core.Byte(core.BRegister),
		codegen.ArithmeticInstruction{
			OpCode:      core.AddWord,
			Destination: core.BRegister,
			Operand1:    core.ARegister,
			Operand2:    core.CRegister,
		},
		// core.Byte(core.AddWord), core.Byte(core.BRegister), core.Byte(core.ARegister), core.Byte(core.CRegister),
		codegen.LoadWordInstruction{Destination: core.ARegister, Address: core.ERegister},
		// core.Byte(core.LoadWord), core.Byte(core.ARegister), core.Byte(core.ERegister),

		// cmp c a 100
		//   set d 100
		//   cmp c a d
		codegen.SetWordInstruction{Destination: core.DRegister, Operand1: 100},
		// core.Byte(core.SetWord), core.Byte(core.DRegister), 100, 0, 0, 0, 0, 0, 0, 0,
		codegen.ArithmeticInstruction{
			OpCode:      core.Compare,
			Destination: core.CRegister,
			Operand1:    core.ARegister,
			Operand2:    core.DRegister,
		},
		// core.Byte(core.Compare), core.Byte(core.CRegister), core.Byte(core.ARegister), core.Byte(core.DRegister),

		// beq c -1 loop
		//   set d -1
		//   beq c d loop
		codegen.SetWordInstruction{Destination: core.DRegister, Operand1: -1},
		// core.Byte(core.SetWord), core.Byte(core.DRegister), 255, 255, 255, 255, 255, 255, 255, 255,
		codegen.BranchEqualInstruction{Operand1: core.CRegister, Operand2: core.DRegister, BranchTarget: 1},
		// core.Byte(core.BranchEqual), core.Byte(core.CRegister), core.Byte(core.DRegister), 10, 0, 0, 0, 0, 0, 0, 0,

		// exit 0
		codegen.ExitInstruction{Operand1: 0},
	})

	expected := core.Memory{
		core.Byte(core.SetWord), core.Byte(core.BRegister), 1, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.Print), core.Byte(core.ARegister),
		core.Byte(core.AddWord), core.Byte(core.ARegister), core.Byte(core.ARegister), core.Byte(core.BRegister),
		core.Byte(core.SetWord), core.Byte(core.CRegister), 0, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.SetWord), core.Byte(core.ERegister), 0, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.StoreWord), core.Byte(core.ERegister), core.Byte(core.BRegister),
		core.Byte(core.AddWord), core.Byte(core.BRegister), core.Byte(core.ARegister), core.Byte(core.CRegister),
		core.Byte(core.LoadWord), core.Byte(core.ARegister), core.Byte(core.ERegister),
		core.Byte(core.SetWord), core.Byte(core.DRegister), 100, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.Compare), core.Byte(core.CRegister), core.Byte(core.ARegister), core.Byte(core.DRegister),
		core.Byte(core.SetWord), core.Byte(core.DRegister), 255, 255, 255, 255, 255, 255, 255, 255,
		core.Byte(core.BranchEqual), core.Byte(core.CRegister), core.Byte(core.DRegister), 10, 0, 0, 0, 0, 0, 0, 0,
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
