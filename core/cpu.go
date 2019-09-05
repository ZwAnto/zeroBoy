package core

type GbCpu struct {
	A byte
	B byte
	C byte
	D byte
	E byte
	F byte
	H byte
	L byte
	SP uint16
	PC uint16
	IME bool
}

// BYTE REGISTER

// set 

func (c *GbCpu) SetA(val byte) {
	c.A = val
}
func (c *GbCpu) SetB(val byte) {
	c.B = val
}
func (c *GbCpu) SetC(val byte) {
	c.C = val
}
func (c *GbCpu) SetD(val byte) {
	c.D = val
}
func (c *GbCpu) SetE(val byte) {
	c.E = val
}
func (c *GbCpu) SetF(val byte) {
	c.F = val
}
func (c *GbCpu) SetH(val byte) {
	c.H = val
}
func (c *GbCpu) SetL(val byte) {
	c.L = val
}

// get

func (c *GbCpu) GetA() byte {
	return c.A
}
func (c *GbCpu) GetB() byte {
	return c.B
}
func (c *GbCpu) GetC() byte {
	return c.C
}
func (c *GbCpu) GetD() byte {
	return c.D
}
func (c *GbCpu) GetE() byte {
	return c.E
}
func (c *GbCpu) GetF() byte {
	return c.F
}
func (c *GbCpu) GetH() byte {
	return c.H
}
func (c *GbCpu) GetL() byte {
	return c.L
}

// UINT16 REGISTER

// set

func (c *GbCpu) SetBC(val uint16) {
	c.B = byte( val >> 8 )
	c.C = byte( val & 0xff )
}
func (c *GbCpu) SetDE(val uint16) {
	c.D = byte( val >> 8 )
	c.E = byte( val & 0xff )
}
func (c *GbCpu) SetHL(val uint16) {
	c.H = byte( val >> 8 )
	c.L = byte( val & 0xff )
}
func (c *GbCpu) SetAF(val uint16) {
	c.A = byte( val >> 8 )
	c.F = byte( val & 0xff )
}
func (c *GbCpu) SetSP(val uint16) {
	c.SP = val
}
func (c *GbCpu) SetPC(val uint16) {
	c.PC = val
}

// get

func (c *GbCpu) GetBC() uint16 {
	return uint16(( c.B & 0xff ) + c.C)
}
func (c *GbCpu) GetDE() uint16 {
	return uint16(( c.D & 0xff ) + c.E)
}
func (c *GbCpu) GetHL() uint16 {
	return uint16(( c.H & 0xff ) + c.L)
}
func (c *GbCpu) GetAF() uint16 {
	return uint16(( c.A & 0xff ) + c.F)
}
func (c *GbCpu) GetSP() uint16 {
	return c.SP
}
func (c *GbCpu) GetPC() uint16 {
	return c.PC
}

// FLAGS
// Get
func (c *GbCpu) GetfZ() byte {
	return c.GetF() >> 7
}
func (c *GbCpu) GetfS() byte {
	return (c.GetF() >> 6) & 1
}
func (c *GbCpu) GetfH() byte {
	return (c.GetF() >> 5) & 1
}
func (c *GbCpu) GetfC() byte {
	return (c.GetF() >> 4) & 1
}
// Set
func (c *GbCpu) SetfZ(val bool)  {
	if val {
		c.SetF(c.GetF() | (1<<7))
	} else {
		c.SetF(c.GetF() &^ (1<<7))
	}
}
func (c *GbCpu) SetfS(val bool) {
	if val {
		c.SetF(c.GetF() | (1<<6))
	} else {
		c.SetF(c.GetF() &^ (1<<6))
	}
}
func (c *GbCpu) SetfH(val bool) {
	if val {
		c.SetF(c.GetF() | (1<<5))
	} else {
		c.SetF(c.GetF() &^ (1<<5))
	}
}
func (c *GbCpu) SetfC(val bool) {
	if val {
		c.SetF(c.GetF() | (1<<4))
	} else {
		c.SetF(c.GetF() &^ (1<<4))
	}
}