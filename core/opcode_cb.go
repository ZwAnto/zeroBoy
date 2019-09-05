package core

func (c *GbCore) setHL(val byte) {
	c.GbMmu.Set(c.GbCpu.GetHL(),val)
}
func (c *GbCore) getHL() byte {
	return c.GbMmu.Get(c.GbCpu.GetHL())
}

func (c *GbCore) getter(i byte) func() byte {

	var f func() byte

	switch i{
	case 1: c.GbCpu.GetB
	case 2: c.GbCpu.GetC
	case 3: c.GbCpu.GetD
	case 4: c.GbCpu.GetE
	case 5: c.GbCpu.GetH
	case 6: c.GbCpu.GetL
	case 7: c.getHL
	case 8: c.GbCpu.GetA
	}

	return f
}
func (c *GbCore) setter(i byte) func(byte) {

	var f func(byte)

	switch i{
		case 1: c.GbCpu.SetB
		case 2: c.GbCpu.SetC
		case 3: c.GbCpu.SetD
		case 4: c.GbCpu.SetE
		case 5: c.GbCpu.SetH
		case 6: c.GbCpu.SetL
		case 7: c.setHL
		case 8: c.GbCpu.SetA
		}
	
		return f
}

func (c *GbCore) OpcodeCB(a byte) {

	in = c.getter((a & 0xf) % 8)()
	out = c.setter((a & 0xf) % 8)

	switch a {
	// RLC
	case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
		carry := in >> 7
		rot := (in << 1) & 0xff | carry
		out(rot)

		c.GbCpu.SetfZ(rot == 0)
        c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)
		c.GbCpu.SetfC(carry == 1)
	// RRC
	case 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		carry := in & 7
		rot := (in >> 1) | (carry << 7 ) 
		out(rot)

		c.GbCpu.SetfZ(rot == 0)
        c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)
		c.GbCpu.SetfC(carry == 1)
	// BIT
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f:
		pos = (a-0x40) / 8
		c.GbCpu.SetfZ((in >> pos) & 1 == 0)
        c.GbCpu.SetfS(false)
        c.GbCpu.SetfH(true)

	// RES
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
		0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
		pos = (a-0x80) / 8
		out(in &^ 1<<pos)
	// SET
	case 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
		0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf,
		0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef,
		0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff:
		pos = (a-0xc0) / 8
		out(in | 1<<pos)
	}
}