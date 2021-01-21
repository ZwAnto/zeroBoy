package core

import (
	// "fmt"

	"time"
)

type GbCore struct {
	GbMmu      GbMmu
	GbCpu      GbCpu
	GbPpu      GbPpu
	ExitSignal bool
}

func (c *GbCore) Init() {
	c.GbMmu.Init()
	c.GbCpu.Init()
	c.GbPpu.Init()
	c.GbMmu.LoadROM()
	c.ExitSignal = false
}

func (c *GbCore) requestInterrupt(interrupt byte) {
	setbyte(&c.GbMmu.Memory[0xFF0F], interrupt)
}

func (c *GbCore) CpuThread() {

	ticker := time.NewTicker(time.Duration(1/c.GbPpu.FPS*1e6) * time.Microsecond)
	stepByFrame := uint64(c.GbCpu.ClockSpeed / c.GbPpu.FPS * 1e6)

	for range ticker.C {

		// Exit when window closed
		if c.ExitSignal == true {
			break
		}

		c.GbCpu.Timer = 0
		for c.GbCpu.Timer <= stepByFrame {

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

			for curClock = 0; curClock <= stepPpu; {
				a := c.GbMmu.Get(c.GbCpu.PC)
				//fmt.Println(strconv.FormatInt(int64(c.GbCpu.PC),16) + ":" + strconv.FormatInt(int64(a),16))

				// if c.GbMmu.FlagBios == false {
				// 	logger.Log.Printf("PC: " + strconv.FormatInt(int64(c.GbCpu.PC), 16))
				// 	logger.Log.Printf("Opcode: " + strconv.FormatInt(int64(a), 16))
				// }

				c.GbCpu.PC++

				// Interrupts
				if c.GbCpu.IME*c.GbMmu.Get(0xffff)&c.GbMmu.Get(0xff0f) > 0 {

					data := c.GbCpu.GetPC()
					c.GbMmu.Set(c.GbCpu.GetSP()-1, byte(data>>8))
					c.GbMmu.Set(c.GbCpu.GetSP()-2, byte(data&0xff))
					c.GbCpu.SetSP(c.GbCpu.GetSP() - 2)

					// VBLANK INTERRUPT
					if c.GbMmu.Get(0xffff)&c.GbMmu.Get(0xff0f)&0x01 > 0 {
						c.GbCpu.SetPC(0x0040)
						c.GbMmu.Set(0xff0f, c.GbMmu.Get(0xff0f)&(255-0x1))
					// LCD-STAT INTERRUPT
					} else if c.GbMmu.Get(0xffff)&c.GbMmu.Get(0xff0f)&0x02 > 0 {
						c.GbCpu.SetPC(0x0048)
						c.GbMmu.Set(0xff0f, c.GbMmu.Get(0xff0f)&(255-0x2))
					// TIMER INTERRUPT
					} else if c.GbMmu.Get(0xffff)&c.GbMmu.Get(0xff0f)&0x04 > 0 {
						c.GbCpu.SetPC(0x0050)
						c.GbMmu.Set(0xff0f, c.GbMmu.Get(0xff0f)&(255-0x4))
					// SERIAL INTERRUPT
					} else if c.GbMmu.Get(0xffff)&c.GbMmu.Get(0xff0f)&0x08 > 0 {
						c.GbCpu.SetPC(0x0058)
						c.GbMmu.Set(0xff0f, c.GbMmu.Get(0xff0f)&(255-0x8))
					// JOYPAD INTERRUPT
					} else if c.GbMmu.Get(0xffff)&c.GbMmu.Get(0xff0f)&0x10 > 0 {
						c.GbCpu.SetPC(0x0060)
						c.GbMmu.Set(0xff0f, c.GbMmu.Get(0xff0f)&(255-0x10))
					}
				}

				if c.GbCpu.PC > 0xff {
					if c.GbMmu.FlagBios == true {

						c.GbCpu.SetPC(0x100)

						c.GbCpu.SetAF(0x01b0)
						c.GbCpu.SetBC(0x0013)
						c.GbCpu.SetDE(0x00d8)
						c.GbCpu.SetHL(0x014d)
						c.GbCpu.SetSP(0xfffe)

						c.GbMmu.Set(0xFF05, 0x00)
						c.GbMmu.Set(0xFF06, 0x00)
						c.GbMmu.Set(0xFF07, 0x00)
						c.GbMmu.Set(0xFF10, 0x80)
						c.GbMmu.Set(0xFF11, 0xBF)
						c.GbMmu.Set(0xFF12, 0xF3)
						c.GbMmu.Set(0xFF14, 0xBF)
						c.GbMmu.Set(0xFF16, 0x3F)
						c.GbMmu.Set(0xFF17, 0x00)
						c.GbMmu.Set(0xFF19, 0xBF)
						c.GbMmu.Set(0xFF1A, 0x7F)
						c.GbMmu.Set(0xFF1B, 0xFF)
						c.GbMmu.Set(0xFF1C, 0x9F)
						c.GbMmu.Set(0xFF1E, 0xBF)
						c.GbMmu.Set(0xFF20, 0xFF)
						c.GbMmu.Set(0xFF21, 0x00)
						c.GbMmu.Set(0xFF22, 0x00)
						c.GbMmu.Set(0xFF23, 0xBF)
						c.GbMmu.Set(0xFF24, 0x77)
						c.GbMmu.Set(0xFF25, 0xF3)
						c.GbMmu.Set(0xFF26, 0xF1)
						c.GbMmu.Set(0xFF40, 0x91)
						c.GbMmu.Set(0xFF42, 0x00)
						c.GbMmu.Set(0xFF43, 0x00)
						c.GbMmu.Set(0xFF45, 0x00)
						c.GbMmu.Set(0xFF47, 0xFC)
						c.GbMmu.Set(0xFF48, 0xFF)
						c.GbMmu.Set(0xFF49, 0xFF)
						c.GbMmu.Set(0xFF4A, 0x00)
						c.GbMmu.Set(0xFF4B, 0x00)
						c.GbMmu.Set(0xFFFF, 0x00)

						c.GbMmu.FlagBios = false

					}
				}

				t := c.Opcode(a)

				// if c.GbCpu.PC >= 0x105 {
				// 	fmt.Println(c.GbCpu.SP)
				// 	op <- true
				// }
				curClock += uint64(t)
			}

			<-c.GbPpu.Trigger

			c.GbCpu.Timer += curClock
		}
	}
}
