package core

import (
	"os"
	"fmt"
	"bufio"
	"path/filepath"
)

type GbMmu struct {
	Memory [0xffff]byte
	Bios []byte
}

func (m *GbMmu) Init() {
	path, _ := filepath.Abs("./DMG_ROM.bin")
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

	for k, v := range bytes {
		m.Memory[k] = v
	}
}

func (m *GbMmu) Get(addr uint16) byte {
	return m.Memory[addr]
}

func (m *GbMmu) Set(addr uint16, value byte) {
	m.Memory[addr] = value
}