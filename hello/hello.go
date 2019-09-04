package main

import (
	"github.com/zwanto/gogb/core"
	"fmt"
)

func main() {
	core := new(core.GbCore)

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
	core.Test(0x79)
	fmt.Println(core.GbCpu.GetA())

	fmt.Println("Hello, world.")
}
