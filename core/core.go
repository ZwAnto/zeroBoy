package core

import "fmt"

type GbCore struct {
	GbMmu GbMmu
}

func (core *GbCore) Init() {
	fmt.Println("Test")
}
