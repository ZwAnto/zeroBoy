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
	assert.Equal(m.Rb(0x0), byte(1))

}

//  8 bit
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
