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

func (c *Cpu) Instruction(op uint16) {

	switch op {
	case 0x00: // NOP
		c.Time += 4
	case 0x10: // STOP
		c.Time += 4 // TODO Check what it does

	case 0x02: // LD (BC), A
		c.Time += 8
		addr := uint16(c.B)<<8 + uint16(c.C)
		c.Mmu.Wb(addr, c.A) // put A into (BC)
	case 0x12: // LD (DE), A
		c.Time += 8
		addr := uint16(c.D)<<8 + uint16(c.E)
		c.Mmu.Wb(addr, c.A) // put A into (BC)

	case 0x40: // LD B, B
		c.Time += 4
		c.B = c.B
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
		c.B = c.Mmu.Rb(addr)
	case 0x47: // LD B, A
		c.Time += 4
		c.B = c.A

	case 0x48: // LD C, B
		c.Time += 4
		c.C = c.B
	case 0x49: // LD C, C
		c.Time += 4
		c.C = c.C
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
		c.C = c.Mmu.Rb(addr)
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
		c.D = c.D
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
		c.D = c.Mmu.Rb(addr)
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
		c.E = c.E
	case 0x5C: // LD E, H
		c.Time += 4
		c.E = c.H
	case 0x5D: // LD E, L
		c.Time += 4
		c.E = c.L
	case 0x5E: // LD E, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.E = c.Mmu.Rb(addr)
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
		c.H = c.H
	case 0x65: // LD H, L
		c.Time += 4
		c.H = c.L
	case 0x66: // LD H, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.H = c.Mmu.Rb(addr)
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
		c.L = c.L
	case 0x6E: // LD L, (HL)
		c.Time += 8
		addr := uint16(c.H)<<8 + uint16(c.L)
		c.L = c.Mmu.Rb(addr)
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
	case 0x76: // HALT
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
		c.A = c.Mmu.Rb(addr)
	case 0x7F: // LD A, A
		c.Time += 4
		c.A = c.A
	}

}
