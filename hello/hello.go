package main

import (
	"github.com/zwanto/gogb/core"
	"fmt"
)

func main() {
	core := new(core.GbCore)

	fmt.Println("==== GoGB 3mul4t0r ===========")

	core.Init()
	core.GbMmu.Init()



	// core.GbCpu.SetAF(300)

	// fmt.Println(core.GbCpu.A)

	// fmt.Println(core.GbCpu.F)

	// core.GbMmu.Set(30,50)
	// fmt.Println(core.GbMmu.Get(30))

	// fmt.Println(core.GbCpu.GetSP())

	core.GbCpu.SetC(30)
	// fmt.Println(core.Operand1(0x0d)())
	core.Opcode(0x22)
	fmt.Println(core.GbCpu.GetHL())

	fmt.Println("Hello, world.")
}
