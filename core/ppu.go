package core

import (
	"fmt"
	"strconv"
)

type GbPpu struct {
	Line byte
	Mode byte
	FPS float64
	Buffer [144][160][3]byte
	Trigger chan bool
	Render chan bool
}

func (p *GbPpu) Init() {
	p.Line = 0
	p.Mode = 2
	p.FPS = 60
	p.Trigger = make(chan bool)
	p.Render = make(chan bool)
}

func (c *GbCore) PpuThread() {

	i := 0
	for ;i<1;{
		<- c.GbCpu.Trigger

		switch {
		case c.GbPpu.Mode == 2 : 
				c.GbPpu.Mode = 3
		case c.GbPpu.Mode == 3 : 
				c.UpdateBuffer()
				c.GbPpu.Mode = 0
		case c.GbPpu.Mode == 0 : 
			c.GbPpu.Line++
			if c.GbPpu.Line == 143 {
				c.GbPpu.Render <- true
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

		c.GbPpu.Trigger <- true
	}
}

func (c *GbCore) UpdateBuffer() {

	var l uint64
	l = (uint64(c.GbPpu.Line) + uint64(c.GbMmu.Get(0xff42)))%256

	mapline_ := int64(l / 8)
	var map_ []byte
	var tile_ []byte

	switch getbyte(c.GbMmu.Get(0xff40),3) {
	case 1 : 
		map_ = c.GbMmu.Memory[(0x9c00 + mapline_ * 32):(0x9c00 + mapline_ * 32 + 32)]
	case 0 : 
		map_ = c.GbMmu.Memory[(0x9800 + mapline_ * 32):(0x9800 + mapline_ * 32 + 32)]
	}

	switch getbyte(c.GbMmu.Get(0xff40),4) {
		case 1 : 
			tile_ = c.GbMmu.Memory[0x8000:0x8fff]
		case 0 : 
			tile_ = c.GbMmu.Memory[0x8800:0x97FF]
		}

	var line_ [256]byte

	for j, v := range map_ {

		tile_down := fmt.Sprintf("%08v",strconv.FormatInt(int64(tile_[16*int64(v) + (int64(l) % 8) * 2])    , 2))
		tile_up   := fmt.Sprintf("%08v",strconv.FormatInt(int64(tile_[16*int64(v) + (int64(l) % 8) * 2 + 1]), 2))

		for i,_ := range(tile_down){

			a, _ := strconv.ParseInt(string(tile_down[i]),10,64)
			b, _ := strconv.ParseInt(string(tile_up[i]),10,64)

			line_[j*8 + i] = byte(a + b * 2)
		}
	}

	for k,_ := range [160]uint8{}{

		v := line_[(k+int(c.GbMmu.Get(0xff43)))%256]

		switch (c.GbMmu.Get(0xff47) >> (v*2)) & 3 {
		case 0:
			c.GbPpu.Buffer[c.GbPpu.Line][k][0] = 255 
			c.GbPpu.Buffer[c.GbPpu.Line][k][1] = 255
			c.GbPpu.Buffer[c.GbPpu.Line][k][2] = 255
		case 1:
			c.GbPpu.Buffer[c.GbPpu.Line][k][0] = 192 
			c.GbPpu.Buffer[c.GbPpu.Line][k][1] = 192
			c.GbPpu.Buffer[c.GbPpu.Line][k][2] = 192
		case 2:
			c.GbPpu.Buffer[c.GbPpu.Line][k][0] = 96 
			c.GbPpu.Buffer[c.GbPpu.Line][k][1] = 96
			c.GbPpu.Buffer[c.GbPpu.Line][k][2] = 96
		case 3:
			c.GbPpu.Buffer[c.GbPpu.Line][k][0] = 0 
			c.GbPpu.Buffer[c.GbPpu.Line][k][1] = 0
			c.GbPpu.Buffer[c.GbPpu.Line][k][2] = 0
		}
	}
}