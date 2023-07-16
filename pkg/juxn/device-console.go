package juxn

import (
	"fmt"
)

type ConsoleDevice struct {
	Vm *VM
}

func (d *ConsoleDevice) Dei(addr byte) byte {
	fmt.Printf("DEI: %d\n", addr)
	return 0
}

func (d *ConsoleDevice) Deo(addr byte, val byte) {
	if addr&0xf == 0x8 {
		fmt.Printf("%c", val)
	}
}
