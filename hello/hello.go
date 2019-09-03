package main

import (
	"github.com/zwanto/gogb/core"
	"fmt"
)

func main() {
	core := new(core.GbCore)

	core.Init()
	core.GbMmu.Init()

	core.GbCpu.SetAF(300)

	fmt.Println(core.GbCpu.A)

	fmt.Println(core.GbCpu.F)

	core.GbMmu.Set(30,50)
	fmt.Println(core.GbMmu.Get(30))

	fmt.Println(core.GbCpu.GetSP())

	fmt.Println("Hello, world.")
}
