package core

import (
	// "fmt"
	"time"
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

func (c *GbCore) CpuThread(op chan bool) {

	ticker := time.NewTicker(time.Duration(230) * time.Nanosecond )
	for range ticker.C {

		a := c.GbMmu.Get(c.GbCpu.PC)
		// fmt.Println(strconv.FormatInt(int64(c.GbCpu.PC),16) + ":" + strconv.FormatInt(int64(a),16))
		c.GbCpu.PC ++ 
		t := c.Opcode(a)
		if c.GbCpu.PC == 0x68 {
			break
		}

		c.GbCpu.Timer += uint64(t)
	
	}

	op <- true
}
