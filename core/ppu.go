package core

import (
	"time"
	"fmt"
)

type GbPpu struct {
	Line byte
	Mode byte
	FPS float64
}

func (p *GbPpu) Init() {
	p.Line = 0
	p.Mode = 2
	p.FPS = 60
}

func (c *GbCore) PpuThread(op chan bool) {

	clock := 0

	ticker := time.NewTicker(time.Duration(230) * time.Nanosecond )

	for range ticker.C {
		clock++

		switch {
		case c.GbPpu.Mode == 2 : 
			if clock >= 80 {
				clock = 0
				c.GbPpu.Mode = 3
			}
		case c.GbPpu.Mode == 3 : 
			if clock >= 172 {
				clock = 0
				c.GbPpu.Mode = 0
			}
		case c.GbPpu.Mode == 0 : 
			if clock >= 204 {
				clock = 0
				c.GbPpu.Line++
				if c.GbPpu.Line == 143 {
					c.GbPpu.Mode = 1
				} else {
					c.GbPpu.Mode = 2
				} 
			}
		case c.GbPpu.Mode == 1 : 
			if clock >= 456 {
				clock = 0
				c.GbPpu.Line++
			}
			if c.GbPpu.Line == 153 {
				c.GbPpu.Line = 0
				c.GbPpu.Mode = 2
				fmt.Println("l,iduznxuez")
			}	
		}
		
		c.GbMmu.Set(0xff44, c.GbPpu.Line) 
	}
	op <- true
}
