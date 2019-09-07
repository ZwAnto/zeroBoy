package core

import (
	//"fmt"
	"time"
	//"strconv"
)
type GbCore struct {
	GbMmu GbMmu
	GbCpu GbCpu
	GbPpu GbPpu
}

func (c *GbCore) Init() {
	c.GbMmu.Init()
	c.GbCpu.Init()
	c.GbPpu.Init()
}

func (c *GbCore) CpuThread(op chan bool) {

	ticker := time.NewTicker(time.Duration(1/c.GbPpu.FPS * 1e6) * time.Microsecond )
	stepByFrame := uint64(c.GbCpu.ClockSpeed / c.GbPpu.FPS * 1e6)
	
	for range ticker.C {
		c.GbCpu.Timer = 0
		for ;c.GbCpu.Timer <= stepByFrame; {

			a := c.GbMmu.Get(c.GbCpu.PC)
			//fmt.Println(strconv.FormatInt(int64(c.GbCpu.PC),16) + ":" + strconv.FormatInt(int64(a),16))
			c.GbCpu.PC ++ 
			t := c.Opcode(a)
			if c.GbCpu.PC == 0xe9 {
				op <- true
			}
			c.GbPpu.Clock += uint64(t)
			c.GbCpu.Timer += uint64(t)
			c.PpuThread()


		}
	}

	op <- true
}
