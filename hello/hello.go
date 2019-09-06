package main

import (
	"github.com/zwanto/goBoy/core"
	"fmt"
	"strconv"
)

func main() {
	core := new(core.GbCore)

	fmt.Println("==== GoGB 3mul4t0r ===========")

	core.Init()
	core.GbMmu.Init()
	i := 0
	for ;i <= 1; { // initialisation and post are omitted
		a := core.GbMmu.Memory[core.GbCpu.PC]
		fmt.Println(strconv.FormatInt(int64(core.GbCpu.PC),16) + ":" + strconv.FormatInt(int64(a),16))
		core.GbCpu.PC ++ 
		core.Opcode(a)
		if core.GbCpu.PC == 0x6b {
			break
		}
    }
	fmt.Println("Hello, world.")
}
