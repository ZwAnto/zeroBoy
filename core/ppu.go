package core

import (
	"fmt"
)

type GbPpu struct {
	Line    byte
	Mode    byte
	FPS     float64
	Buffer  [144][160][3]byte
	Trigger chan bool
	Render  chan bool

	Map     []byte
	Tile    []byte
	Sprites [10][4]byte
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
	for i < 1 {

		// Exit when window closed
		if c.ExitSignal == true {
			break
		}

		<-c.GbCpu.Trigger

		requestInterrupt := false
		requestInterruptVBL := false
		
		switch {

		////////////// HBLANK //////////////
		case c.GbPpu.Mode == 0:

			c.GbPpu.Line++           // go next line
			if c.GbPpu.Line == 144 { // if last line...

				c.GbPpu.Mode = 1 // ...go VBLANK
				requestInterrupt = (getbyte(c.GbMmu.Get(0xff41), 4) == 1) || (getbyte(c.GbMmu.Get(0xff41), 5) == 1) // STAT INTERRUPT
				// requestInterruptVBL = true

				c.GbPpu.Render <- true

			} else { // ... else ...

				c.GbPpu.Mode = 2 // go OAM SEARCH
				requestInterrupt = getbyte(c.GbMmu.Get(0xff41), 5) == 1
			}

		////////////// VBLANK //////////////
		case c.GbPpu.Mode == 1:

			if c.GbPpu.Line == 153 {
				c.GbPpu.Line = 0

				c.GbPpu.Mode = 2 // go OAM SEARCH
				requestInterrupt = getbyte(c.GbMmu.Get(0xff41), 5) == 1

			} else {

				c.GbPpu.Line++

			}

		////////////// OAM SEARCH //////////////
		case c.GbPpu.Mode == 2:

			c.CacheSprites()

			c.GbPpu.Mode = 3 // go PIXEL TRANSFER

		////////////// PIXEL TRANSFER //////////////
		case c.GbPpu.Mode == 3:

			c.UpdateBuffer()

			c.GbPpu.Mode = 0 // go HBLANK
			requestInterrupt = getbyte(c.GbMmu.Get(0xff41), 3) == 1
		}

		c.GbMmu.Set(0xff44, c.GbPpu.Line)
		c.GbMmu.Set(0xff41, c.GbMmu.Get(0xff41)&0xfc+c.GbPpu.Mode)

		if requestInterruptVBL {
			c.requestInterrupt(0) // VBLANK INTERRUPT
		}
		if requestInterrupt {
			c.requestInterrupt(1)
		}
		if c.GbMmu.Get(0xff44) == c.GbMmu.Get(0xff45) {
			setbyte(&c.GbMmu.Memory[0xff41], 2)
			if getbyte(c.GbMmu.Get(0xff41), 6) == 1 {
				c.requestInterrupt(1)
			}
		} else {
			clearbyte(&c.GbMmu.Memory[0xff41], 2)
		}


		c.GbPpu.Trigger <- true
	}
}

func (c *GbCore) setMap(mapline int, maptype string) {

	var mapbit byte

	switch maptype {
	case "bg":
		mapbit = 3
	case "win":
		mapbit = 6
	}

	switch getbyte(c.GbMmu.Get(0xff40), mapbit) {
	case 1:
		c.GbPpu.Map = c.GbMmu.Memory[(0x9c00 + mapline*32):(0x9c00 + mapline*32 + 32)]
	case 0:
		c.GbPpu.Map = c.GbMmu.Memory[(0x9800 + mapline*32):(0x9800 + mapline*32 + 32)]
	}

}

func (c *GbCore) CacheSprites() {

	i := 0
	for k, _ := range [40]byte{} {

		y := c.GbMmu.Memory[0xfe00+k*4]
		x := c.GbMmu.Memory[0xfe00+k*4+1]
		tile := c.GbMmu.Memory[0xfe00+k*4+2]
		attr := c.GbMmu.Memory[0xfe00+k*4+3]

		if x > 0 && (c.GbPpu.Line+16) >= y && (c.GbPpu.Line+16) < y+8 {
			c.GbPpu.Sprites[i][0] = y
			c.GbPpu.Sprites[i][1] = x
			c.GbPpu.Sprites[i][2] = tile
			c.GbPpu.Sprites[i][3] = attr

			i = i + 1
		}
	}
}

func (c *GbCore) setTile() {

	switch getbyte(c.GbMmu.Get(0xff40), 4) {
	case 1:
		c.GbPpu.Tile = c.GbMmu.Memory[0x8000:(0x8fff + 1)]
	case 0:
		c.GbPpu.Tile = c.GbMmu.Memory[0x8800:(0x97FF + 1)]
	}
}

func (c *GbCore) UpdateBuffer() {

	x := c.GbMmu.Get(0xff43)
	y := c.GbMmu.Get(0xff42)

	screen_line := c.GbMmu.Get(0xff44)

	l := (uint64(screen_line) + uint64(y)) % 256

	mapline := int(l / 8)

	// var bg_enable = getbyte(c.GbMmu.Get(0xff40),0) == 1
	// var obj_enable = getbyte(c.GbMmu.Get(0xff40),1) == 1
	// var win_enable = getbyte(c.GbMmu.Get(0xff40),5) == 1

	c.setTile()
	c.setMap(mapline, "bg")

	var b2i = map[byte]int{48: 0, 49: 1}

	// Array of current line 32 tile of 8*8 = 256 bytes
	var tile_line [256]byte

	// Background Loop
	for m, v := range c.GbPpu.Map {

		tile_row := uint64(l) % 8 * 2
		tile_index := uint64(v) * 16

		tile_down := fmt.Sprintf("%08b", c.GbPpu.Tile[tile_index+tile_row])
		tile_up := fmt.Sprintf("%08b", c.GbPpu.Tile[tile_index+tile_row+1])

		for i, _ := range tile_up {

			a := b2i[tile_down[i]]
			b := b2i[tile_up[i]]

			tile_line[m*8+i] = byte(a + b*2)

		}
	}

	pal_map := c.GbMmu.Get(0xff47)
	pal := map[byte]byte{0: pal_map & 0x3, 1: (pal_map & 0xc) >> 2, 2: (pal_map & 0x30) >> 4, 3: pal_map >> 6}

	for k, _ := range [160]uint8{} {

		v := tile_line[(k+int(x))%256]

		switch pal[v] {
		case 0:
			c.GbPpu.Buffer[screen_line][k][0] = 196 //255
			c.GbPpu.Buffer[screen_line][k][1] = 207 //255
			c.GbPpu.Buffer[screen_line][k][2] = 161 //255
		case 1:
			c.GbPpu.Buffer[screen_line][k][0] = 139 //192
			c.GbPpu.Buffer[screen_line][k][1] = 149 //192
			c.GbPpu.Buffer[screen_line][k][2] = 109 //192
		case 2:
			c.GbPpu.Buffer[screen_line][k][0] = 77 //96
			c.GbPpu.Buffer[screen_line][k][1] = 83 //96
			c.GbPpu.Buffer[screen_line][k][2] = 60 //96
		case 3:
			c.GbPpu.Buffer[screen_line][k][0] = 31 //0
			c.GbPpu.Buffer[screen_line][k][1] = 31 //0
			c.GbPpu.Buffer[screen_line][k][2] = 31 //0
		}
	}
}
