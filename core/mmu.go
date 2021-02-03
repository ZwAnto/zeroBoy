package core

import (
	log "github.com/sirupsen/logrus"
)

type Mmu struct {

	Memory [0xffff + 1]byte

}

func (m *Mmu) Rb(addr uint16) *byte {
	log.Debug("Reading byte at ", addr, " from memory and get ", m.Memory[addr])

	return &m.Memory[addr]
}

func (m *Mmu) Wb(addr uint16, value byte) {
	log.Debug("Writing byte at ", addr, " to memory with ", value)

	m.Memory[addr] = value

}

