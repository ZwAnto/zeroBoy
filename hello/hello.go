package main

import (
	"github.com/zwanto/gogb/core"
	"fmt"
)

func main() {
	core := new(core.GbCore)

	core.Init()
	core.GbMmu.Init()

	fmt.Println(core.GbMmu.Get(30))
	fmt.Println("Hello, world.")
}
