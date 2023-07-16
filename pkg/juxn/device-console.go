package juxn

import (
	"fmt"
	"os"
)

type ConsoleDevice struct {
	Vm *VM
}

func (d *ConsoleDevice) Input(addr byte) byte {
	fmt.Printf("DEI: %d\n", addr)
	return 0
}

func (d *ConsoleDevice) Output(addr byte, val byte) {
	if addr&0xf == 0x8 {
		os.Stdout.Write([]byte{val})
	}
}
