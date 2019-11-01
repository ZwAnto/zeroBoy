package main

import (
	"github.com/zwanto/goBoy/core"
	"fmt"
	//"strconv"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
)

// go get -u github.com/faiface/pixel
// go get -u github.com/faiface/pixel/pixelgl

func Run(core *core.GbCore) {
	pixelgl.Run(func() {
		run(core)
	})
}
func run(core *core.GbCore) {

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
		core.ExitSignal = true
	}()

	pixelMap := pixel.MakePictureData(pixel.R(0, 0, 160, 144))

	for {

		// Exit when window closed
		if core.ExitSignal == true {break}

		<- core.GbPpu.Render
		for x := 0; x < 144; x++ {
			for y := 0; y < 160; y++ {
				colour := color.RGBA{R: core.GbPpu.Buffer[x][y][0], G: core.GbPpu.Buffer[x][y][1], B: core.GbPpu.Buffer[x][y][2], A: 0xFF}
				pixelMap.Pix[(143-x)*160+y] = colour
			}
		}

		graph := pixel.NewSprite(pixel.Picture(pixelMap), pixel.R(0, 0, 160, 144))
		mat := pixel.IM
		mat = mat.Moved(win.Bounds().Center())
		mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(3, 3))
		graph.Draw(win, mat)
		win.Update()
	}
}

func main() {
	core := new(core.GbCore)

	p := fmt.Println

	p("|========== goBoy Emulator ==========|")

	core.Init()

	fmt.Printf("| Clock Speed : %.2f Mhz\n",core.GbCpu.ClockSpeed)

	go core.CpuThread()
	go core.PpuThread()
	
	Run(core)

	// for i := 0; i < 1; i++ {
    //     <-operationDone
	// }

	// fmt.Println(core.GbMmu.Memory[0x9800:0x9bff])

	p("|================ END ===============|")
}
