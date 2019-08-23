package main

import (
	"vm/basic_blocks"
	"vm/codegen"
	"vm/core"
)

func main() {
	bb1 := basic_blocks.NewBasicBlock("main")
	bb1.SetWord(core.B, 1)

	bb2 := basic_blocks.NewBasicBlock("loop")
	bb2.Print(core.A)

	bb2.Arithmetic(core.AddWord, core.A, core.A, core.B)

	// swp a b
	bb2.SetWord(core.C, 0)
	bb2.SetWord(core.E, 0)
	bb2.StoreWord(core.E, core.B)
	bb2.Arithmetic(core.AddWord, core.B, core.A, core.C)
	bb2.LoadWord(core.A, core.E)

	// cmp c a 100
	bb2.SetWord(core.D, 100)
	bb2.Arithmetic(core.Compare, core.C, core.A, core.D)

	// beq c -1 loop
	bb2.SetWord(core.D, -1)
	bb2.BranchEqual(core.C, core.D, bb2)

	bb3 := basic_blocks.NewBasicBlock("end")
	// exit 0
	bb3.Exit(0)

	machine := core.CreateVM(0, 8)
	machine.Instructions = codegen.CodeGen(basic_blocks.ResolveLabels([]*basic_blocks.BasicBlock{
		bb1, bb2, bb3,
	}))

	for {
		machine.RunInstruction()
	}
}
