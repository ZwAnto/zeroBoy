package core

import (
	// "fmt"
	// "time"
	// "strconv"
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

func (c *GbCore) CpuThread(op chan bool, step chan uint64) {

	i := 0
	for ;i <= 1; { 

		a := c.GbMmu.Get(c.GbCpu.PC)
		//fmt.Println(strconv.FormatInt(int64(c.GbCpu.PC),16) + ":" + strconv.FormatInt(int64(a),16))
		c.GbCpu.PC ++ 
		t := c.Opcode(a)
		if c.GbCpu.PC == 0xe9 {
			break
		}

		c.GbCpu.Timer += uint64(t)
	
		step <- c.GbCpu.Timer
	}

	op <- true
}
