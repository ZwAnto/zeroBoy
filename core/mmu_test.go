package core

import (
	"testing"
)

func TestWb(t *testing.T) {
	m := new(Mmu)
	m.Wb(0xffff, 1)

    if m.Memory[0xffff] != 1 {
        t.Errorf("Write 1 to 0xffff but read %b", m.Memory[0xffff])
    }
}

func TestRb(t *testing.T) {
	m := new(Mmu)
	m.Wb(0xffff, 1)

    if *m.Rb(0xffff) != 1 {
        t.Errorf("Read %d from 0xffff; Must be 1", m.Rb(0xffff))
    }
}
