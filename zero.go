package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/zwanto/zeroBoy/core"
)

func main() {

	log.SetLevel(log.DebugLevel)

	mmu := new(core.Mmu)
	cpu := new(core.Cpu)

	cpu.Mmu = mmu

	mmu.Wb(0xffff, 1)
	test := mmu.Rb(0xffff)
	log.Debug(test)

}
