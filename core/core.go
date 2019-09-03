package core

import "fmt"

type GbCore struct {
	GbMmu GbMmu
	GbCpu GbCpu
}

func (core *GbCore) Init() {
	fmt.Println("Test")
}
