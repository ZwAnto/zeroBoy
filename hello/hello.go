package main

import (
	"github.com/zwanto/goBoy/core"
	"fmt"
	//"strconv"
	// "time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// go get -u github.com/faiface/pixel/pixel
// go get -u github.com/faiface/pixel/pixelgl

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Test",
		Bounds: pixel.R(0, 0, 160*3, 144*3),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		for !win.Closed() {
		}
	}()
}

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

	
	pixelgl.Run(func() {
		run()
	})
	
	p("|================ END ===============|")
}
