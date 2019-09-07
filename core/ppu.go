package core

type GbPpu struct {
	Line byte
	Mode byte
	Clock uint64
	FPS float64
}

func (p *GbPpu) Init() {
	p.Line = 0
	p.Mode = 2
	p.FPS = 60
	p.Clock = 0
}

func (c *GbCore) PpuThread() {

	switch {
	case c.GbPpu.Mode == 2 : 
		if c.GbPpu.Clock >= 80 {
			c.GbPpu.Clock = 0
			c.GbPpu.Mode = 3
		}
	case c.GbPpu.Mode == 3 : 
		if c.GbPpu.Clock >= 172 {
			c.GbPpu.Clock = 0
			c.GbPpu.Mode = 0
		}
	case c.GbPpu.Mode == 0 : 
		if c.GbPpu.Clock >= 204 {
			c.GbPpu.Clock = 0
			c.GbPpu.Line++
			if c.GbPpu.Line == 143 {
				c.GbPpu.Mode = 1
			} else {
				c.GbPpu.Mode = 2
			} 
		}
	case c.GbPpu.Mode == 1 : 
		if c.GbPpu.Clock >= 456 {
			c.GbPpu.Clock = 0
			c.GbPpu.Line++
		}
		if c.GbPpu.Line == 153 {
			c.GbPpu.Clock = 0
			c.GbPpu.Line = 0
			c.GbPpu.Mode = 2
			
		}	
	}
	
	c.GbMmu.Set(0xff44, c.GbPpu.Line) 
}
