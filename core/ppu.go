package core

import (
	"time"
	// "fmt"
)

type GbPpu struct {
	Line byte
	Mode byte
}

func (p *GbPpu) Init() {
	p.Line = 0
	p.Mode = 2
}

func (c *GbCore) PpuThread(op chan bool, step chan uint64) {
	// var t float64
	var start time.Time

	// freq := 1000 * 1 / c.GbCpu.ClockSpeed
	i := 0
	start = time.Now()
	for ;i <= 1; { 

		switch {
		case c.GbPpu.Mode == 2 : 
			if c.GbCpu.Timer >= 80 {
				c.GbCpu.Timer = 0
				c.GbPpu.Mode = 3
			}
		case c.GbPpu.Mode == 3 : 
			if c.GbCpu.Timer >= 172 {
				c.GbCpu.Timer = 0
				c.GbPpu.Mode = 0
			}
		case c.GbPpu.Mode == 0 : 
			if c.GbCpu.Timer >= 204 {
				c.GbCpu.Timer = 0
				c.GbPpu.Line++
				if c.GbPpu.Line == 143 {
					c.GbPpu.Mode = 1
				} else {
					c.GbPpu.Mode = 2
				} 
			}
		case c.GbPpu.Mode == 1 : 
			if c.GbCpu.Timer >= 456 {
				c.GbCpu.Timer = 0
				c.GbPpu.Line++
			}
			if c.GbPpu.Line == 153 {
				c.GbPpu.Line = 0
				c.GbPpu.Mode = 2

				// 60 FPS timer
				time.Sleep(time.Duration(16666) * time.Microsecond - time.Since(start) )
				start = time.Now()
			}	
		}
		
		c.GbMmu.Set(0xff44, c.GbPpu.Line) 

		// WAIT FOR NEXT CPU INSTRUCTION
		<- step
	}
	op <- true
}
