package core

import (
	"testing"
)

func TestMmu(t *testing.T) {
	m := new(Mmu)
	c := new(Cpu)

	c.Mmu = m

	c.Mmu.Wb(0x0, 1)

	if m.Rb(0x0) != 1 {
		t.Errorf("Cpu Mmu not link to real Mmu", m.Rb(0xffff))
	}

}
