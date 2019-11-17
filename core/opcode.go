package core

func (c *GbCore) Opcode(a byte) byte {
	
	var t byte

	// fallback
	t = 4

	switch a {
	//CB
	case 0xcb:
		next := c.getuint8()
		t = 4 + c.OpcodeCB(next)
	//8bits LD
	case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
		0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
		0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
		0x70, 0x71, 0x72, 0x73, 0x74, 0x75,       0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,

		0x02, 0x12, 0x22, 0x32,
		0x06, 0x16, 0x26, 0x36,
		0x0a, 0x1a, 0x2a, 0x3a,
		0x0e, 0x1e, 0x2e, 0x3e,

		0xe0, 0xf0,
		0xe2, 0xf2,
		0xea, 0xfa:

		c.setter_8(a)(c.operand_r8(a))
		switch a {
		case 0x22, 0x2a:
			c.GbCpu.SetHL(c.GbCpu.GetHL() + 1)
		case 0x32, 0x3a:
			c.GbCpu.SetHL(c.GbCpu.GetHL() - 1)
		} 

		switch {
			case (a >> 4) < 0x4 || (a & 0xf) == 0x6 || (a & 0xf) == 0xe || a >= 0x70 && a <= 0x77 || a == 0xe2 || a == 0xf2: 
				t = 8
			case a == 0xe0 || a == 0xf0:
				t = 12
			case a == 0xea || a == 0xfa:
				t = 16
			default:
				t = 4
		}
	
	// 16bits LD
	case 0x01, 0x11, 0x21, 0x31, 0x08, 0xf9:
		c.setter_16(a)(c.operand_r16(a))
		switch a {
		case 0x08:
			t = 20
		case 0xf9:
			t = 8
		default :
			t = 12
		}
	
	case 0xf8:
    	val := c.getuint8()
		tmp := c.GbCpu.GetSP() ^ uint16(val) ^ (uint16(int32(c.GbCpu.GetSP()) + int32(int8(val))))
		
		c.GbCpu.SetfZ( false )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( (tmp & 0x10) == 0x10 )
		c.GbCpu.SetfC( (tmp & 0x100) == 0x100 )

		c.GbCpu.SetHL(c.GbCpu.GetSP() + uint16(val))

		t = 12
	// ADD
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0xc6:
		
		val := c.operand_l8(a)
		
		add := c.GbCpu.A + val

		c.GbCpu.SetfZ( add == 0 )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( (c.GbCpu.A & 0xf) + (val & 0xf) > 0xf )
		c.GbCpu.SetfC( add > 0xff)

		c.GbCpu.SetA(add)

		switch {
		case (a & 0xf) == 0x6:	
			t = 8
		default :
			t = 4
		}

	case 0xe8:
    	val := c.getuint8()
		tmp := c.GbCpu.GetSP() ^ uint16(val) ^ (uint16(int32(c.GbCpu.GetSP()) + int32(int8(val))))
		
		c.GbCpu.SetfZ( false )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( (tmp & 0x10) == 0x10 )
		c.GbCpu.SetfC( (tmp & 0x100) == 0x100 )

		c.GbCpu.SetSP(c.GbCpu.GetSP() + uint16(val))

		t = 16
	// ADC
	case 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f, 0xce:

		val := c.operand_l8(a)
		add := c.GbCpu.A + val + c.GbCpu.GetfC()

		c.GbCpu.SetfZ( add == 0 )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( (c.GbCpu.A & 0xf) + (val & 0xf) + c.GbCpu.GetfC() > 0xf )
		c.GbCpu.SetfC( add > 0xff)
		
		c.GbCpu.SetA(add)

		switch {
		case (a & 0xf) == 0xe:	
			t = 8
		default :
			t = 4
		}
	//SUB
	case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0xd6:

		val := c.operand_l8(a)
		sub := c.GbCpu.A - val
		
		c.GbCpu.SetfZ( sub == 0 )
		c.GbCpu.SetfS( true )
		c.GbCpu.SetfH( (c.GbCpu.A & 0xf) - (val & 0xf) < 0)
		c.GbCpu.SetfC( sub < 0)
		
		c.GbCpu.SetA(sub)

		switch {
		case (a & 0xf) == 0x6:	
			t = 8
		default :
			t = 4
		}
	//SBC
	case 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f, 0xde:

		val := c.operand_l8(a)
		sub := c.GbCpu.A - val - c.GbCpu.GetfC()

		c.GbCpu.SetfZ( sub == 0 )
		c.GbCpu.SetfS( true )
		c.GbCpu.SetfH( (c.GbCpu.A & 0xf) - (val & 0xf) - c.GbCpu.GetfC() < 0 )
		c.GbCpu.SetfC( sub < 0)
		
		c.GbCpu.SetA(sub)

		switch {
		case (a & 0xf) == 0xe:	
			t = 8
		default :
			t = 4
		}
	//AND
	case 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xe6:

		val := c.GbCpu.A & c.operand_l8(a)

		c.GbCpu.SetfZ( val == 0 )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( true )
		c.GbCpu.SetfC( false )
		
		c.GbCpu.SetA(val)

		switch {
		case (a & 0xf) == 0x6:	
			t = 8
		default :
			t = 4
		}
	//XOR
	case 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xee:

		val := c.GbCpu.A ^ c.operand_l8(a)

		c.GbCpu.SetfZ( val == 0 )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( false )
		c.GbCpu.SetfC( false )
		
		c.GbCpu.SetA(val)

		switch {
		case (a & 0xf) == 0xe:	
			t = 8
		default :
			t = 4
		}
	//OR
	case 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,0xf6:
		
		val := c.GbCpu.A | c.operand_l8(a)

		c.GbCpu.SetfZ( val == 0 )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( false )
		c.GbCpu.SetfC( false )

		c.GbCpu.SetA(val)

		switch {
		case (a & 0xf) == 0x6:	
			t = 8
		default :
			t = 4
		}
	//CP
	case 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf, 0xfe:

		val := c.operand_l8(a)

		c.GbCpu.SetfZ( c.GbCpu.A == val )
		c.GbCpu.SetfS( true )
		c.GbCpu.SetfH( (val & 0x0f) > (c.GbCpu.A & 0x0f) )
		c.GbCpu.SetfC( val > c.GbCpu.A )

		switch {
		case (a & 0xf) == 0xe:	
			t = 8
		default :
			t = 4
		}
	//INC
	case 0x04, 0x14, 0x24, 0x34, 0x0c, 0x1c, 0x2c, 0x3c:

		val := c.operand_l8(a)

		c.GbCpu.SetfZ( val + 1 == 0 )
		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( (val & 0x0f) + (1 & 0x0f) > 0x0f )

		c.setter_8(a)(val + 1)

		switch a {
		case 0x34:	
			t = 12
		default :
			t = 4
		}
	//DEC
	case 0x05, 0x15, 0x25, 0x35, 0x0d, 0x1d, 0x2d, 0x3d:

		val := c.operand_l8(a)

		c.GbCpu.SetfZ( val - 1 == 0 )
		c.GbCpu.SetfS( true )
		c.GbCpu.SetfH( (val & 0x0f) == 0 )

		c.setter_8(a)(val - 1)
		switch a {
		case 0x35:	
			t = 12
		default :
			t = 4
		}
	//16bit ADD
	case 0x09, 0x19, 0x29, 0x39:

		val := c.GbCpu.GetHL() + c.operand_r16(a)

		c.GbCpu.SetfS( false )
		c.GbCpu.SetfH( (c.GbCpu.GetHL() & 0xfff) > ( val & 0xfff) )
		c.GbCpu.SetfC( val > 0xffff )

		c.GbCpu.SetHL(val)
		t = 8

	//16bit INC
	case 0x03, 0x13, 0x23, 0x33:
		c.setter_16(a)(c.operand_l16(a) + 1)
		t = 8
	//16bit DEC
	case 0x0b, 0x1b, 0x2b, 0x3b:
		c.setter_16(a)(c.operand_l16(a) - 1)
		t = 8
	// POP
	case 0xc1, 0xd1, 0xe1, 0xf1:
		val := (uint16(c.GbMmu.Get(c.GbCpu.GetSP() + 1)) << 8) + uint16(c.GbMmu.Get(c.GbCpu.GetSP()))
		c.GbCpu.SetSP(c.GbCpu.GetSP() + 2)
		c.setter_16(a)(val)
		t = 12
	// PUSH
	case 0xc5, 0xd5, 0xe5, 0xf5:
		data := c.operand_l16(a)
		c.GbMmu.Set(c.GbCpu.GetSP() - 1 ,byte(data >> 8))
		c.GbMmu.Set(c.GbCpu.GetSP() - 2 ,byte(data & 0xff))
		c.GbCpu.SetSP(c.GbCpu.GetSP() - 2)
		t = 16
	// JR
	case 0x18, 0x20, 0x28, 0x30, 0x38:
		val := c.getuint8()
		if c.tester(a){
			c.GbCpu.SetPC(uint16(int32(c.GbCpu.GetPC()) + int32(int8(val))))
			t = 12
		}
		t = 8
	// JP
	case 0xc2, 0xc3, 0xca, 0xd2, 0xda :
		val := c.getuint8()
		if c.tester(a){
			c.GbCpu.SetPC(uint16(val))
			t = 16
		}
		t = 12
	case 0xe9:
		c.GbCpu.SetPC(uint16(c.GbMmu.Get(c.GbCpu.GetHL())))
		t = 4
	// CALL
	case 0xc4, 0xd4, 0xcc, 0xdc, 0xcd:
		val := c.getuint8()
		if c.tester(a){
			// push
			data := c.GbCpu.GetPC()
			c.GbMmu.Set(c.GbCpu.GetSP() - 1 ,byte(data >> 8))
			c.GbMmu.Set(c.GbCpu.GetSP() - 2 ,byte(data & 0xff))
			c.GbCpu.SetSP(c.GbCpu.GetSP() - 2)
			// call
			c.GbCpu.SetPC(uint16(val))
			t = 24
		}
		t = 12
	// RET
	case 0xc0, 0xd0, 0xc8, 0xd8, 0xc9, 0xd9:
		if c.tester(a){
			// pop
			val := (uint16(c.GbMmu.Get(c.GbCpu.GetSP() + 1)) << 8) + uint16(c.GbMmu.Get(c.GbCpu.GetSP()))
			c.GbCpu.SetSP(c.GbCpu.GetSP() + 2)
			// call
			c.GbCpu.SetPC(val)
			if a == 0xd9{
				c.GbCpu.IME = 1
			}
			t = 20
		}
		t = 8
	// RST 
	case 0xc7, 0xd7, 0xe7, 0xf7, 0xcf, 0xdf, 0xef, 0xff:
		// push
		data := c.GbCpu.GetPC()
		c.GbMmu.Set(c.GbCpu.GetSP() - 1 ,byte(data >> 8))
		c.GbMmu.Set(c.GbCpu.GetSP() - 2 ,byte(data & 0xff))
		c.GbCpu.SetSP(c.GbCpu.GetSP() - 2)
		// jump
		switch a {
		case 0xc7:
			c.GbCpu.SetPC(0x0000)
		case 0xd7:
			c.GbCpu.SetPC(0x0010)
		case 0xe7:
			c.GbCpu.SetPC(0x0020)
		case 0xf7:
			c.GbCpu.SetPC(0x0030)
		case 0xcf:
			c.GbCpu.SetPC(0x0008)
		case 0xdf:
			c.GbCpu.SetPC(0x0018)
		case 0xef:
			c.GbCpu.SetPC(0x0028)
		case 0xff:
			c.GbCpu.SetPC(0x0038)
		}
		t = 16
	// DAA
	case 0x27:
		switch {
		case c.GbCpu.GetfS() == 0:
			if c.GbCpu.GetfC() == 1 || c.GbCpu.GetA() > 0x99{
				c.GbCpu.SetA(c.GbCpu.GetA() + 0x60)
				c.GbCpu.SetfC(true)
			}
			if c.GbCpu.GetfH() == 1 || c.GbCpu.GetA() & 0x0f > 0x09{
				c.GbCpu.SetA(c.GbCpu.GetA() + 0x06)
				c.GbCpu.SetfH(false)
			}
		case c.GbCpu.GetfC() == 1 && c.GbCpu.GetfH() == 1 :
			c.GbCpu.SetA(c.GbCpu.GetA() + 0x9a)
			c.GbCpu.SetfC(false)
		case c.GbCpu.GetfC() == 1:
			c.GbCpu.SetA(c.GbCpu.GetA() + 0xa0)
		case c.GbCpu.GetfH() == 1:
			c.GbCpu.SetA(c.GbCpu.GetA() + 0xfa)
			c.GbCpu.SetfH(false)
		}
		c.GbCpu.SetfZ(c.GbCpu.GetA() == 0)

		t = 4
	// SCF
	case 0x37:
		c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)
		c.GbCpu.SetfC(true)

		t = 4
	// CPL
	case 0x2f:
		c.GbCpu.SetA(0xff ^ c.GbCpu.GetA())
		c.GbCpu.SetfS(true)
		c.GbCpu.SetfH(true)

		t = 4
	// CCF
	case 0x3f:
		c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)
		c.GbCpu.SetfC(!(c.GbCpu.GetfC()==1))

		t = 4
	// RLCA
	case 0x07:
		c.GbCpu.SetfZ(false)
		c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)
		c.GbCpu.SetfC(c.GbCpu.GetA() > 0x7f)
		c.GbCpu.SetA( ( c.GbCpu.GetA() << 1 ) | ( c.GbCpu.GetA() >> 7 ) )

		t = 4
	// RLA
	case 0x17:
		c.GbCpu.SetfZ(false)
		c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)

		tmp := c.GbCpu.GetfC()
	
		c.GbCpu.SetfC((c.GbCpu.GetA() & 0x80) == 0x80)
		c.GbCpu.SetA(((c.GbCpu.GetA() << 1) & 0xff) | tmp)

		t = 4
	// RRCA
	case 0x0f:
		c.GbCpu.SetfZ(false)
		c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)

		c.GbCpu.SetA( (c.GbCpu.GetA() >> 1) | ((c.GbCpu.GetA() & 1) << 7) )
		c.GbCpu.SetfC( c.GbCpu.GetA() > 0x7F )

		t = 4
	// RRA
	case 0x1f:
		c.GbCpu.SetfZ(false)
		c.GbCpu.SetfS(false)
		c.GbCpu.SetfH(false)
	
		c.GbCpu.SetfC( (1 & c.GbCpu.GetA()) == 1 )
		c.GbCpu.SetA( (c.GbCpu.GetA() >> 1) | c.GbCpu.GetfC() * 0x80 )

		t = 4	
	// DI
	case 0xf3:
		c.GbCpu.IME = 0
	// EI
	case 0xfb:
		c.GbCpu.IME = 1
	}
	return t
} 
