package core

import (
	//"fmt"
	"time"
)
type GbCore struct {
	GbMmu GbMmu
	GbCpu GbCpu
}

func (c *GbCore) Init() {
	c.GbMmu.Init()
	c.GbCpu.Init()
}

func (c *GbCore) CpuThread(op chan bool) {

	freq := 1000 * 1 / c.GbCpu.ClockSpeed

	i := 0
	for ;i <= 1; { 
		start := time.Now()
		a := c.GbMmu.Get(c.GbCpu.PC)
		c.GbCpu.PC ++ 
		t := c.Opcode(a)
		if c.GbCpu.PC == 0x68 {
			break
		}
		c.GbCpu.Timer += uint64(t)
		elapsed := 1-time.Since(start)
		b := time.Duration(freq*float64(t)) * time.Nanosecond - elapsed
		time.Sleep(b)
	}

	op <- true
}