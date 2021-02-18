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

	IME          bool
	IME_enabling bool
}

func (c *Cpu) execute_next_op() {
	opcode := c.read8()
	c.Instruction(opcode)
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

func (c *Cpu) read8() byte {
	c.PC++
	val := *c.Mmu.Rb(c.PC)
	return val
}
func (c *Cpu) read16() uint16 {
	val := uint16(c.read8())
	val = val + uint16(c.read8())<<8
	return val
}
func (c *Cpu) Inc16(l *byte, r *byte) {

	r16 := uint16(*l)<<8 + uint16(*r)
	r16++
	*l = byte(r16 >> 8)
	*r = byte(r16 & 0xff)
}
func (c *Cpu) Dec16(l *byte, r *byte) {

	r16 := uint16(*l)<<8 + uint16(*r)
	r16--
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
func (c *Cpu) jr(cc bool) {
	d8 := c.read8()
	if cc {

		c.PC = uint16(int32(c.PC) + int32(d8))
	}
}
func (c *Cpu) jp(cc bool) {
	d16 := c.read16()
	if cc {

		c.PC = d16 - 1
	}
}
func (c *Cpu) add16HL(r uint16) {
	hl := uint16(c.H)<<8 + uint16(c.L)
	val := int32(hl) + int32(r)

	c.SetfS(false)
	c.SetfH(int32(hl&0xFFF) > (val & 0xFFF))
	c.SetfC(val > 0xffff)

	c.H = byte(uint16(val) >> 8)
	c.L = byte(uint16(val) & 0xff)

}
func (c *Cpu) add16SP() {
	val := c.read8()
	total := uint16(int32(c.SP) + int32(val))
	tmpVal := c.SP ^ uint16(val) ^ total
	c.SetfZ(false)
	c.SetfS(false)
	c.SetfH((tmpVal & 0x10) == 0x10)
	c.SetfC((tmpVal & 0x100) == 0x100)
	c.SP = total

}

func (c *Cpu) popstack() uint16 {

	low := uint16(*c.Mmu.Rb(c.SP))
	high := uint16(*c.Mmu.Rb(c.SP + 1)) << 8
	c.SP += 2
	return high + low

}

func (c *Cpu) pushstack(val uint16) {

	high := byte(val >> 8)
	low := byte(val & 0xff)

	c.Mmu.Wb(c.SP-1, high)
	c.Mmu.Wb(c.SP-2, low)

	c.SP -= 2

}

func (c *Cpu) call(cc bool) {
	addr := c.read16()
	if cc {

		c.pushstack(c.PC + 1)
		c.PC = addr
	}
}

func (c *Cpu) rst(addr uint16) {

	c.pushstack(c.PC)
	c.PC = addr

}

func (c *Cpu) ret(cc bool) {
	addr := c.popstack()
	if cc {
		c.PC = addr
	}
}

func (c *Cpu) rlc(reg *byte) {
	carry := *reg >> 7
	val := *reg<<1 | carry
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(carry == 1)
}
func (c *Cpu) rl(reg *byte) {
	newCarry := *reg >> 7
	oldCarry := c.GetfC()

	val := *reg<<1 | oldCarry
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(newCarry == 1)
}
func (c *Cpu) rrc(reg *byte) {
	carry := *reg & 1
	val := (*reg >> 1) | (carry << 7)
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(carry == 1)
}

func (c *Cpu) rr(reg *byte) {
	newCarry := *reg & 1
	oldCarry := c.GetfC()
	val := (*reg >> 1) | (oldCarry << 7)
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(newCarry == 1)
}
func (c *Cpu) sla(reg *byte) {

	carry := *reg >> 7
	val := *reg << 1
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(carry == 1)

}
func (c *Cpu) sra(reg *byte) {

	carry := *reg & 1
	val := (*reg & 128) | *reg>>1
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(carry == 1)

}
func (c *Cpu) swap(reg *byte) {
	val := (*reg << 4) | (*reg >> 4)
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(false)
}
func (c *Cpu) srl(reg *byte) {

	carry := *reg & 1
	val := *reg >> 1
	*reg = val

	c.SetfZ(val == 0)
	c.SetfS(false)
	c.SetfH(false)
	c.SetfC(carry == 1)

}
func (c *Cpu) bit(reg *byte, b byte) {
	c.SetfZ((*reg>>b)&1 == 0)
	c.SetfS(false)
	c.SetfH(true)
}
func (c *Cpu) set(reg *byte, b byte) {
	*reg = *reg | (1 << b)
}
func (c *Cpu) res(reg *byte, b byte) {
	*reg = *reg &^ (1 << b)
}
