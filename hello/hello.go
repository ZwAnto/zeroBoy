package main

import (
	"github.com/zwanto/goBoy/core"
	"fmt"
	//"strconv"
	// "time"
)

func main() {
	core := new(core.GbCore)

	p := fmt.Println

	p("|========== goBoy Emulator ==========|")

	core.Init()
	
	operationDone := make(chan bool)
	ppuDone := make(chan bool)
	ppuStart := make(chan bool)

	fmt.Printf("| Clock Speed : %.2f Mhz\n",core.GbCpu.ClockSpeed)

	go core.CpuThread(operationDone, ppuDone, ppuStart)
	go core.PpuThread(ppuDone, ppuStart)

	for i := 0; i < 1; i++ {
        <-operationDone
	}
	
	fmt.Println(core.GbPpu.Buffer)

	p("|================ END ===============|")
}
