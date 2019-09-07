package core

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

func (c *GbCore) PpuThread(ppu chan bool, cpu chan bool ) {

	i := 0
	for ;i<1;{
		<- cpu

		switch {
		case c.GbPpu.Mode == 2 : 
				c.GbPpu.Mode = 3
		case c.GbPpu.Mode == 3 : 
				c.GbPpu.Mode = 0
		case c.GbPpu.Mode == 0 : 
			c.GbPpu.Line++
			if c.GbPpu.Line == 143 {
				c.GbPpu.Mode = 1
			} else {
				c.GbPpu.Mode = 2
			}
		case c.GbPpu.Mode == 1 : 
			if c.GbPpu.Line == 153 {
				c.GbPpu.Line = 0
				c.GbPpu.Mode = 2
			} else {
				c.GbPpu.Line++
			}
		}
		
		c.GbMmu.Set(0xff44, c.GbPpu.Line) 

		ppu <- true
	}
}
