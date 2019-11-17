package core

import (
	"fmt"
	//"strconv"
)

type GbPpu struct {
	Line byte
	Mode byte
	FPS float64
	Buffer [144][160][3]byte
	Trigger chan bool
	Render chan bool
	Map []byte
	Tile []byte
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

		// Exit when window closed
		if c.ExitSignal == true {
		
			// fmt.Println(c.GbMmu.Memory[0x9800:0x9bff])
			// fmt.Println(c.GbMmu.Memory[0x9c00:0x9fff])

			break
		}

		<- c.GbCpu.Trigger

		switch {
		case c.GbPpu.Mode == 2 : 
				c.GbPpu.Mode = 3
		case c.GbPpu.Mode == 3 : 
				c.UpdateBuffer()
				c.GbPpu.Mode = 0
		case c.GbPpu.Mode == 0 : 
			c.GbPpu.Line++
			if c.GbPpu.Line == 144 {
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


func (c *GbCore) setMap(mapline int, maptype string){

	var mapbit byte

	switch maptype {
	case "bg":
		mapbit = 3
	case "win":
		mapbit = 6
	}

	switch getbyte(c.GbMmu.Get(0xff40),mapbit) {
		case 1 : 
			c.GbPpu.Map = c.GbMmu.Memory[(0x9c00 + mapline * 32):(0x9c00 + mapline * 32 + 32)]
		case 0 : 
			c.GbPpu.Map = c.GbMmu.Memory[(0x9800 + mapline * 32):(0x9800 + mapline * 32 + 32)]
	}

}

func (c *GbCore) setTile(){

	switch getbyte(c.GbMmu.Get(0xff40),4) {
		case 1 : 
			c.GbPpu.Tile = c.GbMmu.Memory[0x8000:0x8fff]
		case 0 : 
			c.GbPpu.Tile = c.GbMmu.Memory[0x8800:0x97FF]
		}
}

func (c *GbCore) UpdateBuffer() {

	x := c.GbMmu.Get(0xff43)
	y := c.GbMmu.Get(0xff42)

	screen_line := c.GbMmu.Get(0xff44)

	l := (uint64(screen_line) + uint64( y )) % 256

	mapline := int(l / 8)

	// var bg_enable = getbyte(c.GbMmu.Get(0xff40),0) == 1 
	// var obj_enable = getbyte(c.GbMmu.Get(0xff40),1) == 1
	// var win_enable = getbyte(c.GbMmu.Get(0xff40),5) == 1
	
	c.setTile()

	c.setMap(mapline,"bg")

	var b2i = map[byte]int{48: 0, 49: 1}
	
	// Array of current line 32 tile of 8*8 = 256 bytes
	var tile_line [256]byte

	// Background Loop
	for m, v := range c.GbPpu.Map {

		tile_row := uint64(l) % 8 * 2
		tile_index := uint64(v) * 16

		tile_down := fmt.Sprintf("%08b", c.GbPpu.Tile[tile_index + tile_row])
		tile_up   := fmt.Sprintf("%08b", c.GbPpu.Tile[tile_index + tile_row + 1])
		
		for i,_ := range tile_up {

			a := b2i[ tile_down[i] ]
			b := b2i[ tile_up[i] ]

			tile_line[m*8 + i] = byte(a+b*2)

		}
	}

	for k,_ := range [160]uint8{}{

		v := tile_line[(k+int( x ))%256]

		switch (c.GbMmu.Get(0xff47) >> (v*2)) & 3 {
		case 0:
			c.GbPpu.Buffer[screen_line][k][0] = 255 
			c.GbPpu.Buffer[screen_line][k][1] = 255
			c.GbPpu.Buffer[screen_line][k][2] = 255
		case 1:
			c.GbPpu.Buffer[screen_line][k][0] = 192 
			c.GbPpu.Buffer[screen_line][k][1] = 192
			c.GbPpu.Buffer[screen_line][k][2] = 192
		case 2:
			c.GbPpu.Buffer[screen_line][k][0] = 96 
			c.GbPpu.Buffer[screen_line][k][1] = 96
			c.GbPpu.Buffer[screen_line][k][2] = 96
		case 3:
			c.GbPpu.Buffer[screen_line][k][0] = 0 
			c.GbPpu.Buffer[screen_line][k][1] = 0
			c.GbPpu.Buffer[screen_line][k][2] = 0
		}
	}
}