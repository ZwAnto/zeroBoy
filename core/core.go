package core

import (
	// "fmt"
	"time"
	// "strconv"
)
type GbCore struct {
	GbMmu GbMmu
	GbCpu GbCpu
	GbPpu GbPpu
	ExitSignal bool
}

func (c *GbCore) Init() {
	c.GbMmu.Init()
	c.GbCpu.Init()
	c.GbPpu.Init()
	c.GbMmu.LoadROM()
	c.ExitSignal = false
}

func (c *GbCore) CpuThread() {

	ticker := time.NewTicker(time.Duration(1/c.GbPpu.FPS * 1e6) * time.Microsecond )
	stepByFrame := uint64(c.GbCpu.ClockSpeed / c.GbPpu.FPS * 1e6)
	
	for range ticker.C {

		// Exit when window closed
		if c.ExitSignal == true {break}

		c.GbCpu.Timer = 0
		for ;c.GbCpu.Timer <= stepByFrame; {

			var stepPpu uint64
			switch c.GbPpu.Mode {
			case 0: 
				stepPpu = 204
			case 1: 
				stepPpu = 456
			case 2: 
				stepPpu = 80
			case 3: 
				stepPpu = 172
			}

			var curClock uint64 

			c.GbCpu.Trigger <- true

			for curClock = 0;curClock <= stepPpu; {
				a := c.GbMmu.Get(c.GbCpu.PC)
			 	// fmt.Println(strconv.FormatInt(int64(c.GbCpu.PC),16) + ":" + strconv.FormatInt(int64(a),16))
				c.GbCpu.PC ++ 
				t := c.Opcode(a)

				if c.GbCpu.PC > 0x100 {
					c.GbMmu.FlagBios = false
				}

				// if c.GbCpu.PC >= 0x105 {
				// 	fmt.Println(c.GbCpu.SP)
				// 	op <- true
				// }
				curClock += uint64(t)
			}

			<- c.GbPpu.Trigger

			c.GbCpu.Timer += curClock
		}		
	}
}
