package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMmu(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.Mmu.Wb(0x0, 1)
	assert.Equal(*m.Rb(0x0), byte(1))

}

// Read next 8 bits
func TestRead8(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	m.Wb(1, 8)

	assert.Equal(byte(8), c.read8())
	assert.Equal(uint16(1), c.PC)
}

// Read next 16 bits
func TestRead16(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	m.Wb(1, 8)
	m.Wb(2, 1)

	assert.Equal(uint16(264), c.read16())
	assert.Equal(uint16(2), c.PC)
}

// INC DEC 16 bits
func TestIncDec16(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.H = 1
	c.L = 0

	c.Mmu.Wb(0x1, 0x2b)
	c.Mmu.Wb(0x2, 0x23)

	c.execute_next_op()
	assert.Equal(byte(0), c.H)
	assert.Equal(byte(0xff), c.L)
	assert.Equal(uint64(8), c.Time)

	c.execute_next_op()
	assert.Equal(byte(1), c.H)
	assert.Equal(byte(0), c.L)
	assert.Equal(uint64(16), c.Time)
}

// INC DEC 8 bits
func TestIncDec8(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.Mmu.Wb(1, 0x3c)
	c.Mmu.Wb(2, 0x3d)
	c.Mmu.Wb(3, 0x3d)
	c.Mmu.Wb(4, 0x34)

	c.execute_next_op()
	assert.Equal(byte(1), c.A)
	assert.Equal(byte(0), c.GetfS())
	assert.Equal(uint64(4), c.Time)

	c.execute_next_op()
	assert.Equal(byte(0), c.A)
	assert.Equal(byte(1), c.GetfS())
	assert.Equal(byte(1), c.GetfZ())
	assert.Equal(uint64(8), c.Time)

	c.execute_next_op()
	assert.Equal(byte(255), c.A)
	assert.Equal(byte(1), c.GetfH())
	assert.Equal(uint64(12), c.Time)

	c.H = 0
	c.L = 5
	c.Mmu.Wb(5, 5)
	c.execute_next_op()
	assert.Equal(byte(6), *c.Mmu.Rb(5))
	assert.Equal(uint64(24), c.Time)
}

// ADD 8 bit
func TestAdd8(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.A = 3
	c.add8(13)
	assert.Equal(c.A, byte(16))
	assert.Equal(c.F, byte(0b10<<4))

	c.add8(240)
	assert.Equal(c.A, byte(0))
	assert.Equal(c.F, byte(0b1001<<4))
}

// ADC 8 bit
func TestAdc8(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.F = 0b0001 << 4
	c.A = 3
	c.adc8(13)
	assert.Equal(c.A, byte(17))
	assert.Equal(c.F, byte(0b10<<4))

	c.F = 0b0001 << 4
	c.adc8(240)
	assert.Equal(c.A, byte(2))
	assert.Equal(c.F, byte(0b1<<4))
}

// SUB 8 bit
func TestSub8(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.A = 2
	c.sub8(2)
	assert.Equal(c.A, byte(0))
	assert.Equal(c.F, byte(0b1100<<4))

	c.sub8(2)
	assert.Equal(c.A, byte(254))
	assert.Equal(c.F, byte(0b0111<<4))
}

// SBC 8 bit
func TestSbc8(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.F = 0b0001 << 4
	c.A = 2
	c.sbc8(1)
	assert.Equal(c.A, byte(0))
	assert.Equal(c.F, byte(0b1100<<4))

	c.F = 0b0001 << 4
	c.sbc8(1)
	assert.Equal(c.A, byte(254))
	assert.Equal(c.F, byte(0b0111<<4))
}

// AND
func TestAnd(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.A = 0b1011
	c.and8(0b1100)
	assert.Equal(byte(0b1000), c.A)
	assert.Equal(byte(0b0010<<4), c.F)
}

// XOR
func TestXor(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.A = 0b1011
	c.xor8(0b1100)
	assert.Equal(byte(0b0111), c.A)
	assert.Equal(byte(0b0000<<4), c.F)
}

// OR
func TestOr(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.A = 0b1011
	c.or8(0b1100)
	assert.Equal(byte(0b1111), c.A)
	assert.Equal(byte(0b0000<<4), c.F)
}

// CP
func TestCp(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.A = 0b10000
	c.cp8(0b110000)
	assert.Equal(byte(0b0101<<4), c.F)
}

// JR
func TestJr(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.Mmu.Wb(1, 0x18)
	c.Mmu.Wb(2, 5)
	c.execute_next_op()

	assert.Equal(uint16(7), c.PC)
}

// JP
func TestJp(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.Mmu.Wb(1, 0xc3)
	c.Mmu.Wb(2, 6)
	c.Mmu.Wb(3, 1)

	c.execute_next_op()
	assert.Equal(uint16(262-1), c.PC)
	assert.Equal(uint64(16), c.Time)
}

// ADD 16 HL
func TestAdd16HL(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.H = 1
	c.L = 0

	c.B = 1
	c.C = 0

	c.Mmu.Wb(1, 0x09)
	c.execute_next_op()
	assert.Equal(byte(0b10), c.H)
	assert.Equal(byte(0), c.L)
	assert.Equal(uint64(8), c.Time)
}

// ADD 16 SP
func TestAdd16SP(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.Mmu.Wb(1, 0xe8)
	c.Mmu.Wb(2, 1)

	c.SP = 240
	c.execute_next_op()
	assert.Equal(uint16(241), c.SP)
	assert.Equal(uint64(16), c.Time)
}

// SP PUSH POP
func TestSPPUSHPOP(t *testing.T) {

	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.SP = 10
	c.Mmu.Wb(10, 10)
	c.Mmu.Wb(11, 1)

	pop := c.popstack()
	assert.Equal(uint16(266), pop)
	assert.Equal(uint16(12), c.SP)

	c.pushstack(uint16(12))
	assert.Equal(byte(0), *c.Mmu.Rb(11))
	assert.Equal(byte(12), *c.Mmu.Rb(10))
	assert.Equal(uint16(10), c.SP)
}

// CALL
func TestCall(t *testing.T) {

	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.SP = 100

	c.Mmu.Wb(1, 0xcd)
	c.Mmu.Wb(2, 10)
	c.Mmu.Wb(3, 0)

	c.execute_next_op()

	t.Log(c.Time)

	assert.Equal(uint16(10), c.PC)
	assert.Equal(uint16(98), c.SP)
	assert.Equal(uint64(24), c.Time)
	assert.Equal(uint16(4), c.popstack())
}

// RST
func TestRst(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.SP = 100

	c.Mmu.Wb(1, 0xf7)

	c.execute_next_op()

	assert.Equal(uint16(30), c.PC)
	assert.Equal(uint16(98), c.SP)
	assert.Equal(uint64(16), c.Time)
	assert.Equal(uint16(1), c.popstack())

}
