package core

type Cpu struct {
	A  byte
	B  byte
	C  byte
	D  byte
	E  byte
	F  byte
	H  byte
	L  byte
	SP uint16
	PC uint16

	ClockSpeed float64
	Time       uint64

	Mmu *Mmu
}

// FLAGS
// Get
func (c *Cpu) GetfZ() byte {
	return c.F >> 7
}
func (c *Cpu) GetfS() byte {
	return (c.F >> 6) & 1
}
func (c *Cpu) GetfH() byte {
	return (c.F >> 5) & 1
}
func (c *Cpu) GetfC() byte {
	return (c.F >> 4) & 1
}

// Set
func (c *Cpu) SetfZ(val bool) {
	if val {
		c.F = c.F | (1 << 7)
	} else {
		c.F = c.F &^ (1 << 7)
	}
}
func (c *Cpu) SetfS(val bool) {
	if val {
		c.F = c.F | (1 << 6)
	} else {
		c.F = c.F &^ (1 << 6)
	}
}
func (c *Cpu) SetfH(val bool) {
	if val {
		c.F = c.F | (1 << 5)
	} else {
		c.F = c.F &^ (1 << 5)
	}
}
func (c *Cpu) SetfC(val bool) {
	if val {
		c.F = c.F | (1 << 4)
	} else {
		c.F = c.F &^ (1 << 4)
	}
}

func (c *Cpu) Inc16(l *byte, r *byte) {
	r16 := uint16(*l)<<8 + uint16(*r)
	r16 = r16 + 1
	*l = byte(r16 >> 8)
	*r = byte(r16 & 0xff)
}
func (c *Cpu) Dec16(l *byte, r *byte) {
	r16 := uint16(*l)<<8 + uint16(*r)
	r16 = r16 - 1
	*l = byte(r16 >> 8)
	*r = byte(r16 & 0xff)
}

func (c *Cpu) Inc8(r *byte) {
	val := *r + 1
	val4 := uint16(*r&0xf) + 1

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(val4 > 0xf)

	*r = val
}

func (c *Cpu) Dec8(r *byte) {
	val := *r - 1

	c.SetfZ(val == 0)
	c.SetfS(true)
	c.SetfH(1 > *r&0xf)

	*r = val
}

func (c *Cpu) read8() byte {
	val := *c.Mmu.Rb(c.PC)
	c.PC = c.PC + 1
	return val
}

func (c *Cpu) read16() uint16 {
	val := uint16(*c.Mmu.Rb(c.PC))
	c.PC = c.PC + 1
	val = val<<8 + uint16(*c.Mmu.Rb(c.PC))
	c.PC = c.PC + 1

	return val
}

func (c *Cpu) add8(y byte) {
	add := uint16(c.A) + uint16(y)
	add4 := uint16(c.A&0xf) + uint16(y&0xf)

	c.SetfZ(byte(add) == 0)
	c.SetfS(false)
	c.SetfH(add4 > 0xf)
	c.SetfC(add > 0xff)

	c.A = byte(add)
}

func (c *Cpu) adc8(y byte) {
	carry := uint16(c.GetfC())

	add := uint16(c.A) + uint16(y) + carry
	add4 := uint16(c.A&0xf) + uint16(y&0xf) + carry&0xf

	c.SetfZ(byte(add) == 0)
	c.SetfS(false)
	c.SetfH(add4 > 0xf)
	c.SetfC(add > 0xff)

	c.A = byte(add)
}

func (c *Cpu) sub8(y byte) {
	sub := uint16(c.A) - uint16(y)

	c.SetfZ(sub == 0)
	c.SetfS(true)
	c.SetfH(y&0xf > c.A&0xf)
	c.SetfC(y > c.A)

	c.A = byte(sub)
}

func (c *Cpu) sbc8(y byte) {

	carry := uint16(c.GetfC())

	sub := uint16(c.A) - uint16(y) - carry

	c.SetfZ(sub == 0)
	c.SetfS(true)
	c.SetfH(uint16(c.A&0xf) < uint16(y&0xf)+carry&0xf)
	c.SetfC(uint16(c.A) < uint16(y)+carry)

	c.A = byte(sub)
}

func (c *Cpu) and8(y byte) {

	val := c.A & y

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(true)
	c.SetfC(false)

	c.A = val
}

func (c *Cpu) xor8(y byte) {

	val := c.A ^ y

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(false)

	c.A = val
}

func (c *Cpu) or8(y byte) {

	val := c.A | y

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(false)

	c.A = val
}

func (c *Cpu) cp8(y byte) {
	c.SetfZ(c.A == y)
	c.SetfS(true)
	c.SetfH(c.A&0xf < y&0xf)
	c.SetfC(c.A < y)
}

func (c *Cpu) Instruction(op uint16) {

	switch op {
	case 0x00: // NOP
		c.Time += 4
	case 0x01: // LD BC, d16
		c.Time += 12
		val := c.read16()
		c.B = byte(val >> 8)
		c.C = byte(val & 0xff)
	case 0x02: // LD (BC), A
		c.Time += 8
		addr := uint16(c.B)<<8 + uint16(c.C)
		c.Mmu.Wb(addr, c.A) // put A into (BC)
	case 0x03: // INC BC
		c.Time += 8
		c.Inc16(&c.B, &c.C)
	case 0x04: // INC B
		c.Time += 4
		c.Inc8(&c.B)
	case 0x05: // DEC B
		c.Time += 4
		c.Dec8(&c.B)
	case 0x06: // LD B, d8
		c.Time += 8
		c.B = c.read8()
	case 0x07:

	case 0x08: // LD (a8), SP
		// NOT SURE
		c.Time += 20
		w_addr := c.read16()
		c.Mmu.Wb(w_addr, byte(c.SP&0xff))
		c.Mmu.Wb(w_addr+1, byte(c.SP>>8))

	case 0x09:
	case 0x0A: // LD A, (BC)
		c.Time += 8
		addr := uint16(c.B)<<8 + uint16(c.C)
		c.A = *c.Mmu.Rb(addr)
	case 0x0B: // DEC BC
		c.Time += 8
		c.Dec16(&c.B, &c.C)
	case 0x0C: // INC C
		c.Time += 4
		c.Inc8(&c.C)
	case 0x0D: // DEC C
		c.Time += 4
		c.Dec8(&c.C)
	case 0x0E: // LD C, d8
		c.Time += 8
		c.C = c.read8()
	case 0x0F:

	case 0x10: // STOP TODO
		c.Time += 4 // TODO Check what it does
	case 0x11: // LD DE, d16
		c.Time += 12
		val := c.read16()
		c.D = byte(val >> 8)
		c.E = byte(val & 0xff)
	case 0x12: // LD (DE), A
		c.Time += 8
		addr := uint16(c.D)<<8 + uint16(c.E)
		c.Mmu.Wb(addr, c.A) // put A into (BC)
	case 0x13: // INC DE
		c.Time += 8
		c.Inc16(&c.D, &c.E)
	case 0x14: // INC D
		c.Time += 4
		c.Inc8(&c.D)
	case 0x15: // DEC D
		c.Time += 4
		c.Dec8(&c.D)
	case 0x16: // LD D, d8
		c.Time += 8
		c.D = c.read8()
	case 0x17:

	case 0x18:
	case 0x19:
	case 0x1A: // LD A, (DE)
		c.Time += 8
		addr := uint16(c.D)<<8 + uint16(c.E)
		c.A = *c.Mmu.Rb(addr)
	case 0x1B: // DEC DE
		c.Time += 8
		c.Dec16(&c.D, &c.E)
	case 0x1C: // INC E
		c.Time += 4
		c.Inc8(&c.E)
	case 0x1D: // DEC E
		c.Time += 4
		c.Dec8(&c.E)
	case 0x1E: // LD E, d8
		c.Time += 8
		c.E = c.read8()
	case 0x1F:

	case 0x20:
	case 0x21: // LD HL, d16
		c.Time += 12
		val := c.read16()
		c.H = byte(val >> 8)
		c.L = byte(val & 0xff)
	case 0x22: // LD (HL+), A
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.A)
		c.Inc16(&c.H, &c.L)
	case 0x23: // INC HL
		c.Time += 8
		c.Inc16(&c.H, &c.L)
	case 0x24: // INC H
		c.Time += 4
		c.Inc8(&c.H)
	case 0x25: // DEC H
		c.Time += 4
		c.Dec8(&c.H)
	case 0x26: // LD H, d8
		c.Time += 8
		c.H = c.read8()
	case 0x27:

	case 0x28:
	case 0x29:
	case 0x2A: // LD A, (HL+)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.A = *c.Mmu.Rb(addr)
		c.Inc16(&c.H, &c.L)
	case 0x2B: // DEC HL
		c.Time += 8
		c.Dec16(&c.H, &c.L)
	case 0x2C: // INC L
		c.Time += 4
		c.Inc8(&c.L)
	case 0x2D: // DEC L
		c.Time += 4
		c.Dec8(&c.L)
	case 0x2E: // LD L, d8
		c.Time += 8
		c.L = c.read8()
	case 0x2F:

	case 0x30:
	case 0x31: // LD SP, d16
		c.Time += 12
		c.SP = c.read16()
	case 0x32: // LD (HL-), A
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.A)
		c.Dec16(&c.H, &c.L)
	case 0x33: // INC SP
		c.Time += 8
		c.SP += 1
	case 0x34: // INC (HL)
		c.Time += 12
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Inc8(c.Mmu.Rb(addr))
	case 0x35: // DEC (HL)
		c.Time += 12
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Dec8(c.Mmu.Rb(addr))
	case 0x36: // LD (HL), d8
		c.Time += 12
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.read8())
	case 0x37:

	case 0x38:
	case 0x39:
	case 0x3A: // LD A, (HL-)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.A = *c.Mmu.Rb(addr)
		c.Dec16(&c.H, &c.L)
	case 0x3B: // DEC SP
		c.Time += 8
		c.SP -= 1
	case 0x3C: // INC A
		c.Time += 4
		c.Inc8(&c.A)
	case 0x3D: // DEC A
		c.Time += 4
		c.Dec8(&c.A)
	case 0x3E: // LD A, d8
		c.Time += 8
		c.A = c.read8()
	case 0x3F:

	case 0x40: // LD B, B
		c.Time += 4
	case 0x41: // LD B, C
		c.Time += 4
		c.B = c.C
	case 0x42: // LD B, D
		c.Time += 4
		c.B = c.D
	case 0x43: // LD B, E
		c.Time += 4
		c.B = c.E
	case 0x44: // LD B, H
		c.Time += 4
		c.B = c.H
	case 0x45: // LD B, L
		c.Time += 4
		c.B = c.L
	case 0x46: // LD B, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.B = *c.Mmu.Rb(addr)
	case 0x47: // LD B, A
		c.Time += 4
		c.B = c.A

	case 0x48: // LD C, B
		c.Time += 4
		c.C = c.B
	case 0x49: // LD C, C
		c.Time += 4
	case 0x4A: // LD C, D
		c.Time += 4
		c.C = c.D
	case 0x4B: // LD C, E
		c.Time += 4
		c.C = c.E
	case 0x4C: // LD C, H
		c.Time += 4
		c.C = c.H
	case 0x4D: // LD C, L
		c.Time += 4
		c.C = c.L
	case 0x4E: // LD C, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.C = *c.Mmu.Rb(addr)
	case 0x4F: // LD C, A
		c.Time += 4
		c.C = c.A

	case 0x50: // LD D, B
		c.Time += 4
		c.D = c.B
	case 0x51: // LD D, C
		c.Time += 4
		c.D = c.C
	case 0x52: // LD D, D
		c.Time += 4
	case 0x53: // LD D, E
		c.Time += 4
		c.D = c.E
	case 0x54: // LD D, H
		c.Time += 4
		c.D = c.H
	case 0x55: // LD D, L
		c.Time += 4
		c.D = c.L
	case 0x56: // LD D, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.D = *c.Mmu.Rb(addr)
	case 0x57: // LD D, A
		c.Time += 4
		c.D = c.A

	case 0x58: // LD E, B
		c.Time += 4
		c.E = c.B
	case 0x59: // LD E, C
		c.Time += 4
		c.E = c.C
	case 0x5A: // LD E, D
		c.Time += 4
		c.E = c.D
	case 0x5B: // LD E, E
		c.Time += 4
	case 0x5C: // LD E, H
		c.Time += 4
		c.E = c.H
	case 0x5D: // LD E, L
		c.Time += 4
		c.E = c.L
	case 0x5E: // LD E, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.E = *c.Mmu.Rb(addr)
	case 0x5F: // LD E, A
		c.Time += 4
		c.E = c.A

	case 0x60: // LD H, B
		c.Time += 4
		c.H = c.B
	case 0x61: // LD H, C
		c.Time += 4
		c.H = c.C
	case 0x62: // LD H, D
		c.Time += 4
		c.H = c.D
	case 0x63: // LD H, E
		c.Time += 4
		c.H = c.E
	case 0x64: // LD H, H
		c.Time += 4
	case 0x65: // LD H, L
		c.Time += 4
		c.H = c.L
	case 0x66: // LD H, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.H = *c.Mmu.Rb(addr)
	case 0x67: // LD H, A
		c.Time += 4
		c.H = c.A

	case 0x68: // LD L, B
		c.Time += 4
		c.L = c.B
	case 0x69: // LD L, C
		c.Time += 4
		c.L = c.C
	case 0x6A: // LD L, D
		c.Time += 4
		c.L = c.D
	case 0x6B: // LD L, E
		c.Time += 4
		c.L = c.E
	case 0x6C: // LD L, H
		c.Time += 4
		c.L = c.H
	case 0x6D: // LD L, L
		c.Time += 4
	case 0x6E: // LD L, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.L = *c.Mmu.Rb(addr)
	case 0x6F: // LD L, A
		c.Time += 4
		c.L = c.A

	case 0x70: // LD (HL), B
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.B)
	case 0x71: // LD (HL), C
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.C)
	case 0x72: // LD (HL), D
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.D)
	case 0x73: // LD (HL), E
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.E)
	case 0x74: // LD (HL), H
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.H)
	case 0x75: // LD (HL), L
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.L)
	case 0x76: // HALT TODO
		c.Time += 4
		// TODO check what it does
	case 0x77: // LD (HL), A
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.A)

	case 0x78: // LD A, B
		c.Time += 4
		c.A = c.B
	case 0x79: // LD A, C
		c.Time += 4
		c.A = c.C
	case 0x7A: // LD A, D
		c.Time += 4
		c.A = c.D
	case 0x7B: // LD A, E
		c.Time += 4
		c.A = c.E
	case 0x7C: // LD A, H
		c.Time += 4
		c.A = c.H
	case 0x7D: // LD A, L
		c.Time += 4
		c.A = c.L
	case 0x7E: // LD A, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.A = *c.Mmu.Rb(addr)
	case 0x7F: // LD A, A
		c.Time += 4

	case 0x80: // ADD A, B
		c.Time += 4
		c.add8(c.B)
	case 0x81: // ADD A, C
		c.Time += 4
		c.add8(c.C)
	case 0x82: // ADD A, D
		c.Time += 4
		c.add8(c.D)
	case 0x83: // ADD A, E
		c.Time += 4
		c.add8(c.E)
	case 0x84: // ADD A, H
		c.Time += 4
		c.add8(c.H)
	case 0x85: // ADD A, L
		c.Time += 4
		c.add8(c.L)
	case 0x86: // ADD A, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.add8(*c.Mmu.Rb(addr))
	case 0x87: // ADD A, A
		c.Time += 4
		c.add8(c.A)

	case 0x88: // ADC A, B
		c.Time += 4
		c.adc8(c.B)
	case 0x89: // ADC A, C
		c.Time += 4
		c.adc8(c.C)
	case 0x8A: // ADC A, D
		c.Time += 4
		c.adc8(c.D)
	case 0x8B: // ADC A, E
		c.Time += 4
		c.adc8(c.E)
	case 0x8C: // ADC A, H
		c.Time += 4
		c.adc8(c.H)
	case 0x8D: // ADC A, L
		c.Time += 4
		c.adc8(c.L)
	case 0x8E: // ADC A, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.adc8(*c.Mmu.Rb(addr))
	case 0x8F: // ADC A, A
		c.Time += 4
		c.adc8(c.A)

	case 0x90: // SUB A, B
		c.Time += 4
		c.sub8(c.B)
	case 0x91: // SUB A, C
		c.Time += 4
		c.sub8(c.C)
	case 0x92: // SUB A, D
		c.Time += 4
		c.sub8(c.D)
	case 0x93: // SUB A, E
		c.Time += 4
		c.sub8(c.E)
	case 0x94: // SUB A, H
		c.Time += 4
		c.sub8(c.H)
	case 0x95: // SUB A, L
		c.Time += 4
		c.sub8(c.L)
	case 0x96: // SUB A, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.sub8(*c.Mmu.Rb(addr))
	case 0x97: // SUB A, A
		c.Time += 4
		c.sub8(c.A)

	case 0x98: // SBC A, B
		c.Time += 4
		c.sbc8(c.B)
	case 0x99: // SBC A, C
		c.Time += 4
		c.sbc8(c.C)
	case 0x9A: // SBC A, D
		c.Time += 4
		c.sbc8(c.D)
	case 0x9B: // SBC A, E
		c.Time += 4
		c.sbc8(c.E)
	case 0x9C: // SBC A, H
		c.Time += 4
		c.sbc8(c.H)
	case 0x9D: // SBC A, L
		c.Time += 4
		c.sbc8(c.L)
	case 0x9E: // SBC A, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.sbc8(*c.Mmu.Rb(addr))
	case 0x9F: // SBC A, A
		c.Time += 4
		c.sbc8(c.A)

	case 0xA0: // AND B
		c.Time += 4
		c.and8(c.B)
	case 0xA1: // AND C
		c.Time += 4
		c.and8(c.C)
	case 0xA2: // AND D
		c.Time += 4
		c.and8(c.D)
	case 0xA3: // AND E
		c.Time += 4
		c.and8(c.E)
	case 0xA4: // AND H
		c.Time += 4
		c.and8(c.H)
	case 0xA5: // AND L
		c.Time += 4
		c.and8(c.L)
	case 0xA6: // AND (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.and8(*c.Mmu.Rb(addr))
	case 0xA7: // AND A
		c.Time += 4
		c.and8(c.A)

	case 0xA8: // XOR B
		c.Time += 4
		c.xor8(c.B)
	case 0xA9: // XOR C
		c.Time += 4
		c.xor8(c.C)
	case 0xAA: // XOR D
		c.Time += 4
		c.xor8(c.D)
	case 0xAB: // XOR E
		c.Time += 4
		c.xor8(c.E)
	case 0xAC: // XOR H
		c.Time += 4
		c.xor8(c.H)
	case 0xAD: // XOR L
		c.Time += 4
		c.xor8(c.L)
	case 0xAE: // XOR (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.xor8(*c.Mmu.Rb(addr))
	case 0xAF: // XOR A
		c.Time += 4
		c.xor8(c.A)

	case 0xB0: // OR B
		c.Time += 4
		c.or8(c.B)
	case 0xB1: // OR C
		c.Time += 4
		c.or8(c.B)
	case 0xB2: // OR D
		c.Time += 4
		c.or8(c.B)
	case 0xB3: // OR E
		c.Time += 4
		c.or8(c.B)
	case 0xB4: // OR H
		c.Time += 4
		c.or8(c.B)
	case 0xB5: // OR L
		c.Time += 4
		c.or8(c.B)
	case 0xB6: // OR (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.or8(*c.Mmu.Rb(addr))
	case 0xB7: // OR A
		c.Time += 4
		c.or8(c.B)

	case 0xB8: // CP B
		c.Time += 4
		c.cp8(c.B)
	case 0xB9: // CP C
		c.Time += 4
		c.cp8(c.C)
	case 0xBA: // CP D
		c.Time += 4
		c.cp8(c.D)
	case 0xBB: // CP E
		c.Time += 4
		c.cp8(c.E)
	case 0xBC: // CP H
		c.Time += 4
		c.cp8(c.H)
	case 0xBD: // CP L
		c.Time += 4
		c.cp8(c.L)
	case 0xBE: // CP (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.cp8(*c.Mmu.Rb(addr))
	case 0xBF: // CP A
		c.Time += 4
		c.cp8(c.A)

	case 0xC0:
	case 0xC1:
	case 0xC2:
	case 0xC3:
	case 0xC4:
	case 0xC5:
	case 0xC6: // ADD A, d8
		c.Time += 8
		c.add8(c.read8())
	case 0xC7:

	case 0xC8:
	case 0xC9:
	case 0xCA:
	case 0xCB:
	case 0xCC:
	case 0xCD:
	case 0xCE: // ADC A, d8
		c.Time += 8
		c.adc8(c.read8())
	case 0xCF:

	case 0xD0:
	case 0xD1:
	case 0xD2:
	case 0xD3:
	case 0xD4:
	case 0xD5:
	case 0xD6: // SUB d8
		c.Time += 8
		c.sub8(c.read8())
	case 0xD7:

	case 0xD8:
	case 0xD9:
	case 0xDA:
	case 0xDB:
	case 0xDC:
	case 0xDD:
	case 0xDE: // SBC A, d8
		c.Time += 8
		c.sbc8(c.read8())
	case 0xDF:

	case 0xE0: // LDH (a8), A
		c.Time += 12
		addr := uint16(0xff00) + uint16(c.read8())
		c.Mmu.Wb(addr, c.A)
	case 0xE1:
	case 0xE2: // LDH (C), A
		c.Time += 8
		addr := uint16(0xff00) + uint16(c.C)
		c.Mmu.Wb(addr, c.A)
	case 0xE3:
	case 0xE4:
	case 0xE5:
	case 0xE6: // AND d8
		c.Time += 8
		c.and8(c.read8())
	case 0xE7:

	case 0xE8:
	case 0xE9:
	case 0xEA: // LD (a16), A
		c.Time += 16
		addr := c.read16()
		c.Mmu.Wb(addr, c.A)
	case 0xEB:
	case 0xEC:
	case 0xED:
	case 0xEE: // XOR d8
		c.Time += 8
		c.xor8(c.read8())
	case 0xEF:

	case 0xF0: // LDH A, (a8)
		c.Time += 12
		addr := uint16(0xff00) + uint16(c.read8())
		c.A = *c.Mmu.Rb(addr)
	case 0xF1:
	case 0xF2: // LDH A, (C)
		c.Time += 8
		addr := uint16(0xff00) + uint16(c.C)
		c.A = *c.Mmu.Rb(addr)
	case 0xF3:
	case 0xF4:
	case 0xF5:
	case 0xF6: // OR d8
		c.Time += 8
		c.or8(c.read8())
	case 0xF7:

	case 0xF8: // LD HL,SP+r8
		// NOT SURE
		c.Time += 12
		r8 := c.read8()
		val := uint16(int16(c.SP) + int16(r8))

		tmpVal := c.SP ^ uint16(r8) ^ val

		c.SetfZ(false)
		c.SetfH(false)
		c.SetfH((tmpVal & 0x10) == 0x10)
		c.SetfC((tmpVal & 0x100) == 0x100)

		c.H = byte(val >> 8)
		c.L = byte(val & 0xff)
	case 0xF9: // LD SP, HL
		c.Time += 8
		val := uint16(c.H)<<8 + uint16(c.L)
		c.SP = val
	case 0xFA: // LD A, (a16)
		c.Time += 16
		addr := c.read16()
		c.A = *c.Mmu.Rb(addr)
	case 0xFB:
	case 0xFC:
	case 0xFD:
	case 0xFE: // CP d8
		c.Time += 8
		c.cp8(c.read8())
	case 0xFF:
	}

}
