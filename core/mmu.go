package core

import (
	"os"
	"fmt"
	"bufio"
	"path/filepath"
)

type GbMmu struct {
	Memory [0xffff]int8
	Bios [0x100]int8
}

func (m *GbMmu) Init() {
	path, _ := filepath.Abs("../src/github.com/zwanto/gogb/DMG_ROM.bin")
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err) 
		fmt.Println("An Error Occured") 
	}
	stats, statsErr := file.Stat()
    if statsErr != nil {
        fmt.Println(statsErr) 
		fmt.Println("An Error Occured") 
    }

    var size int64 = stats.Size()
    bytes := make([]byte, size)

    bufr := bufio.NewReader(file)
	_,err = bufr.Read(bytes)
	
	fmt.Println(bytes)
	fmt.Println(err)
}

func (m *GbMmu) Get(addr int16) int8 {
	return m.Memory[addr]
}

func (m *GbMmu) Set(addr int16, value int8) {
	m.Memory[addr] = value
}