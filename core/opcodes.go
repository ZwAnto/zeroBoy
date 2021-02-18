package core

func (c *Cpu) InstrTime(op byte) {

	times := []uint64{
		4, 12, 8, 8, 4, 4, 8, 4, 20, 8, 8, 8, 4, 4, 8, 4,
		4, 12, 8, 8, 4, 4, 8, 4, 12, 8, 8, 8, 4, 4, 8, 4,
		8, 12, 8, 8, 4, 4, 8, 4, 8, 8, 8, 8, 4, 4, 8, 4,
		8, 12, 8, 8, 12, 12, 12, 4, 8, 8, 8, 8, 4, 4, 8, 4,
		4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
		4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
		4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
		8, 8, 8, 8, 8, 8, 4, 8, 4, 4, 4, 4, 4, 4, 8, 4,
		4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
		4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
		4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
		4, 4, 4, 4, 4, 4, 8, 4, 4, 4, 4, 4, 4, 4, 8, 4,
		8, 12, 12, 16, 12, 16, 8, 16, 8, 16, 12, 4, 12, 24, 8, 16,
		8, 12, 12, 99, 12, 16, 8, 16, 8, 16, 12, 99, 12, 99, 8, 16,
		12, 12, 8, 99, 99, 16, 8, 16, 16, 4, 16, 99, 99, 99, 8, 16,
		12, 12, 8, 4, 99, 16, 8, 16, 12, 8, 16, 4, 99, 99, 8, 16,
	}
	c.Time += times[op]
}

func (c *Cpu) Instruction(op byte) {

	switch op {
	case 0x00: // NOP
	case 0x01: // LD BC, d16
		val := c.read16()
		c.B = byte(val >> 8)
		c.C = byte(val & 0xff)
	case 0x02: // LD (BC), A
		addr := uint16(c.B)<<8 + uint16(c.C)
		c.Mmu.Wb(addr, c.A) // put A into (BC)
	case 0x03: // INC BC
		c.Inc16(&c.B, &c.C)
	case 0x04: // INC B
		c.Inc8(&c.B)
	case 0x05: // DEC B
		c.Dec8(&c.B)
	case 0x06: // LD B, d8
		c.B = c.read8()
	case 0x07: // RLCA
		c.SetfZ(false)
		c.SetfS(false)
		c.SetfH(false)
		c.SetfC(c.A>>7 == 1)
		c.A = c.A<<1 | c.A>>7
	case 0x08: // LD (a16), SP
		// NOT SURE
		w_addr := c.read16()
		c.Mmu.Wb(w_addr, byte(c.SP&0xff))
		c.Mmu.Wb(w_addr+1, byte(c.SP>>8))
	case 0x09: // ADD HL, BC
		val := uint16(c.B)<<8 + uint16(c.C)
		c.add16HL(val)
	case 0x0A: // LD A, (BC)
		addr := uint16(c.B)<<8 + uint16(c.C)
		c.A = *c.Mmu.Rb(addr)
	case 0x0B: // DEC BC
		c.Dec16(&c.B, &c.C)
	case 0x0C: // INC C
		c.Inc8(&c.C)
	case 0x0D: // DEC C
		c.Dec8(&c.C)
	case 0x0E: // LD C, d8
		c.C = c.read8()
	case 0x0F: // RRCA
		c.SetfZ(false)
		c.SetfS(false)
		c.SetfH(false)
		c.SetfC(c.A&1 == 1)
		c.A = c.A>>1 | ((c.A & 1) << 7)
	case 0x10: // STOP TODO // TODO Check what it does
	case 0x11: // LD DE, d16
		val := c.read16()
		c.D = byte(val >> 8)
		c.E = byte(val & 0xff)
	case 0x12: // LD (DE), A
		addr := uint16(c.D)<<8 + uint16(c.E)
		c.Mmu.Wb(addr, c.A) // put A into (BC)
	case 0x13: // INC DE
		c.Inc16(&c.D, &c.E)
	case 0x14: // INC D
		c.Inc8(&c.D)
	case 0x15: // DEC D
		c.Dec8(&c.D)
	case 0x16: // LD D, d8
		c.D = c.read8()
	case 0x17: // RLA
		carry := c.GetfC()

		c.SetfZ(false)
		c.SetfS(false)
		c.SetfH(false)
		c.SetfC(c.A>>7 == 1)
		c.A = c.A<<1 + carry
	case 0x18: // JR a8
		c.jr(true)
	case 0x19: // ADD HL, DE
		val := uint16(c.D)<<8 + uint16(c.E)
		c.add16HL(val)
	case 0x1A: // LD A, (DE)
		addr := uint16(c.D)<<8 + uint16(c.E)
		c.A = *c.Mmu.Rb(addr)
	case 0x1B: // DEC DE
		c.Dec16(&c.D, &c.E)
	case 0x1C: // INC E
		c.Inc8(&c.E)
	case 0x1D: // DEC E
		c.Dec8(&c.E)
	case 0x1E: // LD E, d8
		c.E = c.read8()
	case 0x1F: // RRA
		carry := c.GetfC()

		c.SetfZ(false)
		c.SetfS(false)
		c.SetfH(false)
		c.SetfC(c.A&1 == 1)
		c.A = (c.A >> 1) + (carry << 7)
	case 0x20: // JR NZ, r8
		c.jr(c.GetfZ() == 0)
	case 0x21: // LD HL, d16
		val := c.read16()
		c.H = byte(val >> 8)
		c.L = byte(val & 0xff)
	case 0x22: // LD (HL+), A
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.A)
		c.Inc16(&c.H, &c.L)
	case 0x23: // INC HL
		c.Inc16(&c.H, &c.L)
	case 0x24: // INC H
		c.Inc8(&c.H)
	case 0x25: // DEC H
		c.Dec8(&c.H)
	case 0x26: // LD H, d8
		c.H = c.read8()
	case 0x27: // DAA
		if c.GetfS() == 0 {
			if c.GetfC() == 1 || c.A > 0x99 {
				c.A += 0x60
				c.SetfC(true)
			}
			if c.GetfH() == 1 || c.A&0xf > 0x9 {
				c.A += 0x06
				c.SetfH(true)
			}
		} else if c.GetfC() == 1 && c.GetfH() == 1 {
			c.A += 0x9a
			c.SetfH(false)
		} else if c.GetfC() == 1 {
			c.A += 0xa0
		} else if c.GetfH() == 1 {
			c.A += 0xfa
			c.SetfH(false)
		}
		c.SetfZ(c.A == 0)
	case 0x28: // JR Z, r8
		c.jr(c.GetfZ() == 1)
	case 0x29: // ADD HL, HL
		val := uint16(c.H)<<8 + uint16(c.L)
		c.add16HL(val)
	case 0x2A: // LD A, (HL+)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.A = *c.Mmu.Rb(addr)
		c.Inc16(&c.H, &c.L)
	case 0x2B: // DEC HL
		c.Dec16(&c.H, &c.L)
	case 0x2C: // INC L
		c.Inc8(&c.L)
	case 0x2D: // DEC L
		c.Dec8(&c.L)
	case 0x2E: // LD L, d8
		c.L = c.read8()
	case 0x2F: // CPL
		c.A = 0xFF ^ c.A
		c.SetfS(true)
		c.SetfH(true)
	case 0x30: // JR NC, r8
		c.jr(c.GetfC() == 0)
	case 0x31: // LD SP, d16
		c.SP = c.read16()
	case 0x32: // LD (HL-), A
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.A)
		c.Dec16(&c.H, &c.L)
	case 0x33: // INC SP

		c.SP += 1
	case 0x34: // INC (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Inc8(c.Mmu.Rb(addr))
	case 0x35: // DEC (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Dec8(c.Mmu.Rb(addr))
	case 0x36: // LD (HL), d8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.read8())
	case 0x37: //SCF
		c.SetfS(false)
		c.SetfH(false)
		c.SetfC(true)
	case 0x38: // JR C, r8
		c.jr(c.GetfC() == 1)
	case 0x39: // ADD HL, SP
		c.add16HL(c.SP)
	case 0x3A: // LD A, (HL-)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.A = *c.Mmu.Rb(addr)
		c.Dec16(&c.H, &c.L)
	case 0x3B: // DEC SP
		c.SP -= 1
	case 0x3C: // INC A
		c.Inc8(&c.A)
	case 0x3D: // DEC A
		c.Dec8(&c.A)
	case 0x3E: // LD A, d8
		c.A = c.read8()
	case 0x3F: // CCF
		c.SetfS(false)
		c.SetfH(false)
		c.SetfC(c.GetfC() == 0)
	case 0x40: // LD B, B
	case 0x41: // LD B, C
		c.B = c.C
	case 0x42: // LD B, D
		c.B = c.D
	case 0x43: // LD B, E
		c.B = c.E
	case 0x44: // LD B, H
		c.B = c.H
	case 0x45: // LD B, L
		c.B = c.L
	case 0x46: // LD B, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.B = *c.Mmu.Rb(addr)
	case 0x47: // LD B, A
		c.B = c.A
	case 0x48: // LD C, B
		c.C = c.B
	case 0x49: // LD C, C
	case 0x4A: // LD C, D
		c.C = c.D
	case 0x4B: // LD C, E
		c.C = c.E
	case 0x4C: // LD C, H
		c.C = c.H
	case 0x4D: // LD C, L
		c.C = c.L
	case 0x4E: // LD C, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.C = *c.Mmu.Rb(addr)
	case 0x4F: // LD C, A
		c.C = c.A
	case 0x50: // LD D, B
		c.D = c.B
	case 0x51: // LD D, C
		c.D = c.C
	case 0x52: // LD D, D
	case 0x53: // LD D, E
		c.D = c.E
	case 0x54: // LD D, H
		c.D = c.H
	case 0x55: // LD D, L
		c.D = c.L
	case 0x56: // LD D, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.D = *c.Mmu.Rb(addr)
	case 0x57: // LD D, A
		c.D = c.A
	case 0x58: // LD E, B
		c.E = c.B
	case 0x59: // LD E, C
		c.E = c.C
	case 0x5A: // LD E, D
		c.E = c.D
	case 0x5B: // LD E, E
	case 0x5C: // LD E, H
		c.E = c.H
	case 0x5D: // LD E, L
		c.E = c.L
	case 0x5E: // LD E, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.E = *c.Mmu.Rb(addr)
	case 0x5F: // LD E, A
		c.E = c.A
	case 0x60: // LD H, B
		c.H = c.B
	case 0x61: // LD H, C
		c.H = c.C
	case 0x62: // LD H, D
		c.H = c.D
	case 0x63: // LD H, E
		c.H = c.E
	case 0x64: // LD H, H
	case 0x65: // LD H, L
		c.H = c.L
	case 0x66: // LD H, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.H = *c.Mmu.Rb(addr)
	case 0x67: // LD H, A
		c.H = c.A
	case 0x68: // LD L, B
		c.L = c.B
	case 0x69: // LD L, C
		c.L = c.C
	case 0x6A: // LD L, D
		c.L = c.D
	case 0x6B: // LD L, E
		c.L = c.E
	case 0x6C: // LD L, H
		c.L = c.H
	case 0x6D: // LD L, L
	case 0x6E: // LD L, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.L = *c.Mmu.Rb(addr)
	case 0x6F: // LD L, A
		c.L = c.A
	case 0x70: // LD (HL), B
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.B)
	case 0x71: // LD (HL), C
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.C)
	case 0x72: // LD (HL), D
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.D)
	case 0x73: // LD (HL), E
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.E)
	case 0x74: // LD (HL), H
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.H)
	case 0x75: // LD (HL), L
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.L)
	case 0x76: // HALT TODO
		// TODO check what it does
	case 0x77: // LD (HL), A
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.Mmu.Wb(addr, c.A)
	case 0x78: // LD A, B
		c.A = c.B
	case 0x79: // LD A, C
		c.A = c.C
	case 0x7A: // LD A, D
		c.A = c.D
	case 0x7B: // LD A, E
		c.A = c.E
	case 0x7C: // LD A, H
		c.A = c.H
	case 0x7D: // LD A, L
		c.A = c.L
	case 0x7E: // LD A, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.A = *c.Mmu.Rb(addr)
	case 0x7F: // LD A, A
	case 0x80: // ADD A, B
		c.add8(c.B)
	case 0x81: // ADD A, C
		c.add8(c.C)
	case 0x82: // ADD A, D
		c.add8(c.D)
	case 0x83: // ADD A, E
		c.add8(c.E)
	case 0x84: // ADD A, H
		c.add8(c.H)
	case 0x85: // ADD A, L
		c.add8(c.L)
	case 0x86: // ADD A, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.add8(*c.Mmu.Rb(addr))
	case 0x87: // ADD A, A
		c.add8(c.A)
	case 0x88: // ADC A, B
		c.adc8(c.B)
	case 0x89: // ADC A, C
		c.adc8(c.C)
	case 0x8A: // ADC A, D
		c.adc8(c.D)
	case 0x8B: // ADC A, E
		c.adc8(c.E)
	case 0x8C: // ADC A, H
		c.adc8(c.H)
	case 0x8D: // ADC A, L
		c.adc8(c.L)
	case 0x8E: // ADC A, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.adc8(*c.Mmu.Rb(addr))
	case 0x8F: // ADC A, A
		c.adc8(c.A)
	case 0x90: // SUB A, B
		c.sub8(c.B)
	case 0x91: // SUB A, C
		c.sub8(c.C)
	case 0x92: // SUB A, D
		c.sub8(c.D)
	case 0x93: // SUB A, E
		c.sub8(c.E)
	case 0x94: // SUB A, H
		c.sub8(c.H)
	case 0x95: // SUB A, L
		c.sub8(c.L)
	case 0x96: // SUB A, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.sub8(*c.Mmu.Rb(addr))
	case 0x97: // SUB A, A
		c.sub8(c.A)
	case 0x98: // SBC A, B
		c.sbc8(c.B)
	case 0x99: // SBC A, C
		c.sbc8(c.C)
	case 0x9A: // SBC A, D
		c.sbc8(c.D)
	case 0x9B: // SBC A, E
		c.sbc8(c.E)
	case 0x9C: // SBC A, H
		c.sbc8(c.H)
	case 0x9D: // SBC A, L
		c.sbc8(c.L)
	case 0x9E: // SBC A, (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.sbc8(*c.Mmu.Rb(addr))
	case 0x9F: // SBC A, A
		c.sbc8(c.A)
	case 0xA0: // AND B
		c.and8(c.B)
	case 0xA1: // AND C
		c.and8(c.C)
	case 0xA2: // AND D
		c.and8(c.D)
	case 0xA3: // AND E
		c.and8(c.E)
	case 0xA4: // AND H
		c.and8(c.H)
	case 0xA5: // AND L
		c.and8(c.L)
	case 0xA6: // AND (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.and8(*c.Mmu.Rb(addr))
	case 0xA7: // AND A
		c.and8(c.A)
	case 0xA8: // XOR B
		c.xor8(c.B)
	case 0xA9: // XOR C
		c.xor8(c.C)
	case 0xAA: // XOR D
		c.xor8(c.D)
	case 0xAB: // XOR E
		c.xor8(c.E)
	case 0xAC: // XOR H
		c.xor8(c.H)
	case 0xAD: // XOR L
		c.xor8(c.L)
	case 0xAE: // XOR (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.xor8(*c.Mmu.Rb(addr))
	case 0xAF: // XOR A
		c.xor8(c.A)
	case 0xB0: // OR B
		c.or8(c.B)
	case 0xB1: // OR C
		c.or8(c.B)
	case 0xB2: // OR D
		c.or8(c.B)
	case 0xB3: // OR E
		c.or8(c.B)
	case 0xB4: // OR H
		c.or8(c.B)
	case 0xB5: // OR L
		c.or8(c.B)
	case 0xB6: // OR (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.or8(*c.Mmu.Rb(addr))
	case 0xB7: // OR A
		c.or8(c.B)
	case 0xB8: // CP B
		c.cp8(c.B)
	case 0xB9: // CP C
		c.cp8(c.C)
	case 0xBA: // CP D
		c.cp8(c.D)
	case 0xBB: // CP E
		c.cp8(c.E)
	case 0xBC: // CP H
		c.cp8(c.H)
	case 0xBD: // CP L
		c.cp8(c.L)
	case 0xBE: // CP (HL)
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.cp8(*c.Mmu.Rb(addr))
	case 0xBF: // CP A
		c.cp8(c.A)
	case 0xC0: // RET NZ
		c.ret(c.GetfZ() == 0)
	case 0xC1: // POP BC
		val := c.popstack()
		c.B = byte(val >> 8)
		c.C = byte(val & 0xff)
	case 0xC2: // JPNZ, a16
		c.jp(c.GetfZ() == 0)
	case 0xC3: // JP a16
		c.jp(true)
	case 0xC4: // CALL NZ, a8
		c.call(c.GetfZ() == 0)
	case 0xC5: // PUSH BC
		val := uint16(c.B)<<8 + uint16(c.C)
		c.pushstack(val)
	case 0xC6: // ADD A, d8
		c.add8(c.read8())
	case 0xC7: // RST 00H
		c.rst(0)
	case 0xC8: // RET Z
		c.ret(c.GetfZ() == 1)
	case 0xC9: // RET
		c.ret(true)
	case 0xCA: // JPZ, a16
		c.jp(c.GetfZ() == 1)
	case 0xCB: // PREFIX CB
		c.InstructionCB(c.read8())
	case 0xCC: // CALL Z, a8
		c.call(c.GetfZ() == 1)
	case 0xCD: // CALL a8
		c.call(true)
	case 0xCE: // ADC A, d8
		c.adc8(c.read8())
	case 0xCF: // RST 08H
		c.rst(8)
	case 0xD0: // RET NC
		c.ret(c.GetfC() == 0)
	case 0xD1: // POP DE
		val := c.popstack()
		c.D = byte(val >> 8)
		c.E = byte(val & 0xff)
	case 0xD2: // JPNC, a16
		c.jp(c.GetfC() == 0)
	case 0xD4: // CALL NC, a8
		c.call(c.GetfC() == 0)
	case 0xD5: // PUSH DE
		val := uint16(c.D)<<8 + uint16(c.E)
		c.pushstack(val)
	case 0xD6: // SUB d8
		c.sub8(c.read8())
	case 0xD7: // RST 10H
		c.rst(10)
	case 0xD8: // RET C
		c.ret(c.GetfC() == 1)
	case 0xD9: // RETI
		c.ret(true)
		c.IME_enabling = true
	case 0xDA: // JP C, a16
		c.jp(c.GetfC() == 1)
	case 0xDC: // CALL C, a8
		c.call(c.GetfC() == 1)
	case 0xDE: // SBC A, d8
		c.sbc8(c.read8())
	case 0xDF: // RST 18H
		c.rst(18)
	case 0xE0: // LDH (a8), A
		addr := uint16(0xff00) + uint16(c.read8())
		c.Mmu.Wb(addr, c.A)
	case 0xE1: // POP HL
		val := c.popstack()
		c.H = byte(val >> 8)
		c.L = byte(val & 0xff)
	case 0xE2: // LDH (C), A
		addr := uint16(0xff00) + uint16(c.C)
		c.Mmu.Wb(addr, c.A)
	case 0xE5: // PUSH HL
		val := uint16(c.H)<<8 + uint16(c.L)
		c.pushstack(val)
	case 0xE6: // AND d8
		c.and8(c.read8())
	case 0xE7: // RST 20H
		c.rst(20)
	case 0xE8: // ADD SP, r8
		c.add16SP()
	case 0xE9: // JP (HL)
		c.PC = uint16(c.H)<<8 + uint16(c.L) - uint16(1)
	case 0xEA: // LD (a16), A
		addr := c.read16()
		c.Mmu.Wb(addr, c.A)
	case 0xEE: // XOR d8
		c.xor8(c.read8())
	case 0xEF: // RST 28H
		c.rst(28)
	case 0xF0: // LDH A, (a8)
		addr := uint16(0xff00) + uint16(c.read8())
		c.A = *c.Mmu.Rb(addr)
	case 0xF1: // POP AF
		val := c.popstack()
		c.A = byte(val >> 8)
		c.F = byte(val & 0xff)
	case 0xF2: // LDH A, (C)
		addr := uint16(0xff00) + uint16(c.C)
		c.A = *c.Mmu.Rb(addr)
	case 0xF3: // DI
		c.IME = false
	case 0xF5: // PUSH AF
		val := uint16(c.A)<<8 + uint16(c.F)
		c.pushstack(val)
	case 0xF6: // OR d8
		c.or8(c.read8())
	case 0xF7: // RST 30H
		c.rst(30)
	case 0xF8: // LD HL,SP+r8
		// NOT SURE
		r8 := c.read8()
		val := uint16(int32(c.SP) + int32(r8))

		tmpVal := c.SP ^ uint16(r8) ^ val

		c.SetfZ(false)
		c.SetfH(false)
		c.SetfH((tmpVal & 0x10) == 0x10)
		c.SetfC((tmpVal & 0x100) == 0x100)

		c.H = byte(val >> 8)
		c.L = byte(val & 0xff)
	case 0xF9: // LD SP, HL
		val := uint16(c.H)<<8 + uint16(c.L)
		c.SP = val
	case 0xFA: // LD A, (a16)
		addr := c.read16()
		c.A = *c.Mmu.Rb(addr)
	case 0xFB: // EI
		c.IME_enabling = true
	case 0xFE: // CP d8
		c.cp8(c.read8())
	case 0xFF: // RST 38H
		c.rst(38)
	}

	c.InstrTime(op)
}

func (c *Cpu) InstructionCB(op byte) {

	col := op & 0xff
	col_mod := col % 8
	row := op >> 4

	var reg *uint8

	switch col_mod {
	case 0:
		reg = &c.B
	case 1:
		reg = &c.C
	case 2:
		reg = &c.D
	case 3:
		reg = &c.E
	case 4:
		reg = &c.H
	case 5:
		reg = &c.L
	case 6:
		addr := uint16(c.H)<<8 + uint16(c.L)
		reg = c.Mmu.Rb(addr)
		c.Time += 8
	case 7:
		reg = &c.A
	}

	switch {
	case row == 0x0:
		if col < 0x8 {
			c.rlc(reg)
		} else {
			c.rrc(reg)
		}
	case row == 0x1:
		if col < 0x8 {
			c.rl(reg)
		} else {
			c.rr(reg)
		}
	case row == 0x2:
		if col < 0x8 {
			c.sla(reg)
		} else {
			c.sra(reg)
		}
	case row == 0x3:
		if col < 0x8 {
			c.swap(reg)
		} else {
			c.srl(reg)
		}
	case row <= 0x7:
		c.bit(reg, row-0x4+col/8)
	case row <= 0xb:
		c.res(reg, row-0x8+col/8)
	case row <= 0xf:
		c.set(reg, row-0xc+col/8)
	}

	c.Time += 8
}
