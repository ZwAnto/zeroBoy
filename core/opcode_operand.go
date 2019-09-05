package core

func (c *GbCore) getuint8() byte {
	val := c.GbMmu.Get(c.GbCpu.GetPC())
	c.GbCpu.PC++
	return val
}

func (c *GbCore) setHL(val byte) {
	c.GbMmu.Set(c.GbCpu.GetHL(),val)
}
func (c *GbCore) setBC(val byte) {
	c.GbMmu.Set(c.GbCpu.GetBC(),val)
}
func (c *GbCore) setDE(val byte) {
	c.GbMmu.Set(c.GbCpu.GetDE(),val)
}
func (c *GbCore) getHL() byte {
	return c.GbMmu.Get(c.GbCpu.GetHL())
}
func (c *GbCore) getBC() byte {
	return c.GbMmu.Get(c.GbCpu.GetBC())
}
func (c *GbCore) getDE() byte {
	return c.GbMmu.Get(c.GbCpu.GetDE())
}

// 8bit setter operator
func (c *GbCore) setter(val byte) func(byte) {

	var f func(byte) 

	switch val {
	// B
	case 0x04, 0x05, 0x06,
		0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 
		0x90, 
		0xa0, 0xa8,
		0xb0, 0xb8:
		f = c.GbCpu.SetB
	// C
	case 0x0c, 0x0d, 0x0e, 
		0x38,
		0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
		0x91,
		0xa1, 0xa9,
		0xb1, 0xb9,
		0xd8, 0xda, 0xdc:
		f = c.GbCpu.SetC
	// D
	case 0x14, 0x15, 0x16,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57,
		0x92,
		0xa2, 0xaa,
		0xb2, 0xba:
		f = c.GbCpu.SetD
	// E
	case 0x1c, 0x1d, 0x1e,
		0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
		0x93,
		0xa3, 0xab,
		0xb3, 0xbb:
		f = c.GbCpu.SetE
	// H
	case 0x24, 0x25, 0x26,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67,
		0x94,
		0xa4, 0xac,
		0xb4, 0xbc:
		f = c.GbCpu.SetH
	// L
	case 0x2c, 0x2d, 0x2e,
		0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
		0x95,
		0xa5, 0xad,
		0xb5, 0xbd:
		f = c.GbCpu.SetL
	// (HL)
	case 0x34, 0x35, 0x36,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75,       0x77,
		0x96,
		0xa6, 0xae,
		0xb6, 0xbe,
		0xe9:
		f = c.setHL
	// A
	case 0x0a,
		0x1a,
		0x2a,
		0x3a, 0x3c, 0x3d, 0x3e,
		0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f, 
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		0xa7, 0xaf,
		0xb7, 0xbf,
		0xc6, 0xce,
		0xde,
		0xf0, 0xf2, 0xfa:
		f = c.GbCpu.SetA
	//(BC)
	case 0x02:
		f = c.setBC
	//(DE)
	case 0x12:
		f = c.setDE
	}
	return f
}

// 8bit left operand

func (c *GbCore) operand1(val byte) func() byte{

	var f func() byte 

	switch val {
	// A
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
		0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
		0x48, 0x58, 0x68, 0x78, 0x88, 0x98:
		f = c.GbCpu.GetA
	}
	return f
}

// 8bit right operand

func (c *GbCore) operand2(val byte) func() byte{

	var f func() byte 

	switch val {
	// B
	case 0x40, 0x50, 0x60, 0x70, 0x80,
		0x48, 0x58, 0x68, 0x78, 0x88, 0x98:
		f = c.GbCpu.GetB
	// C
	case 0x41, 0x51, 0x61, 0x71, 0x81,
		0x49, 0x59, 0x69, 0x79, 0x89, 0x99:
		f = c.GbCpu.GetC
	// D
	case 0x42, 0x52, 0x62, 0x72, 0x82,
		0x4a, 0x5a, 0x6a, 0x7a, 0x8a, 0x9a:
		f = c.GbCpu.GetD
	// E
	case 0x43, 0x53, 0x63, 0x73, 0x83,
		0x4b, 0x5b, 0x6b, 0x7b, 0x8b, 0x9b:
		f = c.GbCpu.GetE
	// H
	case 0x44, 0x54, 0x64, 0x74, 0x84,
		0x4c, 0x5c, 0x6c, 0x7c, 0x8c, 0x9c:
		f = c.GbCpu.GetH
	// L
	case 0x45, 0x55, 0x65, 0x75, 0x85,
		0x4d, 0x5d, 0x6d, 0x7d, 0x8d, 0x9d:
		f = c.GbCpu.GetL
	// (HL)
	case 0x46, 0x56, 0x66,       0x86,
		0x4e, 0x5e, 0x6e, 0x7e, 0x8e, 0x9e:
		f = c.getHL
	// A
	case 0x02, 0x12, 0x22, 0x32, 
		0x47, 0x57, 0x67, 0x77, 0x87,
		0x4f, 0x5f, 0x6f, 0x7f, 0x8f, 0x9f,
		0xe0, 0xe2, 0xea:
		f = c.GbCpu.GetA
	// (BC)
	case 0x0a:
		f = c.getBC
	// (DE)
	case 0x1a:
		f = c.getDE
	// d8
	case 0x06, 0x16, 0x26, 0x36, 0x0e, 0x1e, 0x2e, 0x3e:
		f =  c.getuint8
	}

	return f
}
