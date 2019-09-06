package core

// RETRIEVING DATA 

func (c *GbCore) getbyte() byte {
	/*
	Retrieve next byte according to PC
	Increment PC
	*/
	val := c.GbMmu.Get(c.GbCpu.GetPC())
	c.GbCpu.PC++
	return val
}
func (c *GbCore) getuint16() uint16 {
	/*
	Retrieve next two bytes according to PC
	Increment PC two times
	*/
	var val uint16
	val = (uint16(c.GbMmu.Get(c.GbCpu.GetPC() + 1)) << 8) | uint16(c.GbMmu.Get(c.GbCpu.GetPC()))
	c.GbCpu.PC++
	c.GbCpu.PC++
	return val
}

// PUT VALUE INTO MEMORY ACCORDING TO SOMETHING

func (c *GbCore) msetA32(val uint16) {
	/*
	Put 16 bits data into memory according to a 16 bits address 
	*/
	addr := c.getuint16()
	c.GbMmu.Set(addr, byte(val & 0xff))
	c.GbMmu.Set(addr + 1, byte((val & 0xff00) >> 8))
}
func (c *GbCore) msetA16(val byte) {
	/*
	Put 8 bits data into memory according to a 16 bits address 
	*/
	c.GbMmu.Set(c.getuint16(), val)
}
func (c *GbCore) msetA8(val byte) {
	/*
	Put 8 bits data into memory according to a 8 bits address 
	*/
	c.GbMmu.Set( 0xff00 + uint16(c.getbyte()), val)
}
func (c *GbCore) msetC(val byte) {
	/*
	Put 8 bits data into memory according to C
	*/
	c.GbMmu.Set( 0xff00 + uint16(c.GbCpu.GetC()), val)
}
func (c *GbCore) msetHL(val byte) {
	/*
	Put 8 bits data into memory according to HL
	*/
	c.GbMmu.Set(c.GbCpu.GetHL(),val)
}
func (c *GbCore) msetBC(val byte) {
	/*
	Put 8 bits data into memory according to BC
	*/
	c.GbMmu.Set(c.GbCpu.GetBC(),val)
}
func (c *GbCore) msetDE(val byte) {
	/*
	Put 8 bits data into memory according to DE
	*/
	c.GbMmu.Set(c.GbCpu.GetDE(),val)
}