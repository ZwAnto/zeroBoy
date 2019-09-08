package core

import (
	"fmt"
	"strconv"
)

type GbPpu struct {
	Line byte
	Mode byte
	FPS float64
	Buffer [144][160]byte
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
				c.PpuGetMap()
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
		c.GbMmu.Set(0xff41, c.GbMmu.Get(0xff41) & 0xfc + c.GbPpu.Mode)

		ppu <- true
	}
}

func (c *GbCore) PpuGetMap() {

	var l byte
	l = c.GbPpu.Line + c.GbMmu.Get(0xff42)

	mapline_ := int64(l / 8)
	var map_ []byte

	switch getbyte(c.GbMmu.Get(0xff40),6) {
	case 1 : 
		map_ = c.GbMmu.Memory[(0x9c00 + mapline_ * 32):(0x9c00 + mapline_ * 32 + 32)]
	case 0 : 
		map_ = c.GbMmu.Memory[(0x9800 + mapline_ * 32):(0x9800 + mapline_ * 32 + 32)]
	}

	tile_ := c.GbMmu.Memory[0x8000:0x8fff]
	
	var line_ [256]byte

	for j, v := range map_ {

		tile_down := fmt.Sprintf("%08v",strconv.FormatInt(int64(tile_[16*v + (l % 8) * 2])    , 2))
		tile_up   := fmt.Sprintf("%08v",strconv.FormatInt(int64(tile_[16*v + (l % 8) * 2 + 1]), 2))

		for i,_ := range(tile_down){

			a, _ := strconv.ParseInt(string(tile_down[i]),10,64)
			b, _ := strconv.ParseInt(string(tile_up[i]),10,64)

			line_[j*8 + i] = byte(a + b * 2)
		}
	}

	for k,v := range line_[c.GbMmu.Get(0xff43):(c.GbMmu.Get(0xff43) + 160)] {
		c.GbPpu.Buffer[c.GbPpu.Line][k] = v 
	}
}