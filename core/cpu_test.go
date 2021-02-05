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

func TestIncDec16(t *testing.T) {
	assert := assert.New(t)
	c := new(Cpu)

	c.H=1
	c.L=1
	hl := uint16(c.H)<<8 + uint16(c.L)
	assert.Equal(uint16(257), hl)

	c.L=0
	c.Dec16(&c.H, &c.L)
	hl = uint16(c.H)<<8 + uint16(c.L)
	assert.Equal(uint16(255), hl)

	c.Inc16(&c.H, &c.L)
	hl = uint16(c.H)<<8 + uint16(c.L)
	assert.Equal(uint16(256), hl)
}

func TestIncDec8(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.Inc8(&c.A)
	assert.Equal(byte(1), c.A)
	assert.Equal(byte(0), c.GetfS())

	c.Dec8(&c.A)
	assert.Equal(byte(0), c.A)
	assert.Equal(byte(1), c.GetfS())
	assert.Equal(byte(1), c.GetfZ())

	c.Dec8(&c.A)
	assert.Equal(byte(255), c.A)
	assert.Equal(byte(1), c.GetfH())

	addr := uint16(c.H)<<8 + uint16(c.L)
	c.Inc8(c.Mmu.Rb(addr))
	assert.Equal(byte(1), *c.Mmu.Rb(0))

}

// Read next 8 bits
func TestRead8(t *testing.T){
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.PC=1
	m.Wb(1,8)

	assert.Equal(byte(8), c.read8())
	assert.Equal(uint16(2), c.PC)
}

// Read next 16 bits
func TestRead16(t *testing.T){
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.PC=1
	m.Wb(1,8)
	m.Wb(2,1)

	assert.Equal(uint16(264), c.read16())
	assert.Equal(uint16(3), c.PC)
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

	c.Mmu.Wb(0,3)

	c.jr(true)
	assert.Equal(uint16(4), c.PC)
}

// JP
func TestJp(t *testing.T) {
	assert := assert.New(t)
	m := new(Mmu)
	c := new(Cpu)
	c.Mmu = m

	c.Mmu.Wb(0,6)
	c.Mmu.Wb(1,1)

	c.jp(true)
	assert.Equal(uint16(262), c.PC)
}