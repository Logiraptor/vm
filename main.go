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
	loopBuilder := basic_blocks.InstructionBuilder{}
	// loop:
	// print a
	loopBuilder.Print(core.ARegister)

	// add a a b
	loopBuilder.Arithmetic(core.AddWord, core.ARegister, core.ARegister, core.BRegister)

	// swp a b
	//   set c 0
	//   set e 0
	//   store e b
	//   add b a c
	//   load e a
	loopBuilder.SetWord(core.CRegister, 0)
	loopBuilder.SetWord(core.ERegister, 0)
	loopBuilder.StoreWord(core.ERegister, core.BRegister)
	loopBuilder.Arithmetic(core.AddWord, core.BRegister, core.ARegister, core.CRegister)
	loopBuilder.LoadWord(core.ARegister, core.ERegister)

	// cmp c a 100
	//   set d 100
	//   cmp c a d
	loopBuilder.SetWord(core.DRegister, 100)
	loopBuilder.Arithmetic(core.Compare, core.CRegister, core.ARegister, core.DRegister)

	// beq c -1 loop
	//   set d -1
	//   beq c d loop
	loopBuilder.SetWord(core.DRegister, -1)
	loopBuilder.BranchEqual(core.CRegister, core.DRegister, bb2)

	// exit 0
	loopBuilder.Exit(0)
	bb2.Instructions = loopBuilder.Instructions
	machine := core.CreateVM(0, 8)
	machine.Instructions = codegen.CodeGen(basic_blocks.ResolveLabels([]*basic_blocks.BasicBlock{
		bb1, bb2,
	}))

	for {
		machine.RunInstruction()
	}
}
