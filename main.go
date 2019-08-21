package main

import "vm/core"

func main() {
	machine := core.CreateVM(10000, 10000)
	machine.Instructions = core.Memory{
		// set b 1
		core.Byte(core.SetWord), core.Byte(core.BRegister), 1, 0, 0, 0, 0, 0, 0, 0,

		// loop:
		// print a
		core.Byte(core.Print), core.Byte(core.ARegister),

		// add a a b
		core.Byte(core.AddWord), core.Byte(core.ARegister), core.Byte(core.ARegister), core.Byte(core.BRegister),

		// swp a b
		//   set c 0
		//   set e 0
		//   store e b
		//   add b a c
		//   load e a
		core.Byte(core.SetWord), core.Byte(core.CRegister), 0, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.SetWord), core.Byte(core.ERegister), 0, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.StoreWord), core.Byte(core.ERegister), core.Byte(core.BRegister),
		core.Byte(core.AddWord), core.Byte(core.BRegister), core.Byte(core.ARegister), core.Byte(core.CRegister),
		core.Byte(core.LoadWord), core.Byte(core.ARegister), core.Byte(core.ERegister),

		// cmp c a 100
		//   set d 100
		//   cmp c a d
		core.Byte(core.SetWord), core.Byte(core.DRegister), 100, 0, 0, 0, 0, 0, 0, 0,
		core.Byte(core.Compare), core.Byte(core.CRegister), core.Byte(core.ARegister), core.Byte(core.DRegister),

		// beq c -1 loop
		//   set d -1
		//   beq c d loop
		core.Byte(core.SetWord), core.Byte(core.DRegister), 255, 255, 255, 255, 255, 255, 255, 255,
		core.Byte(core.BranchEqual), core.Byte(core.CRegister), core.Byte(core.DRegister), 10, 0, 0, 0, 0, 0, 0, 0,

		// exit 0
		core.Byte(core.Exit), core.Byte(0),
	}

	for {
		machine.RunInstruction()
	}
}
