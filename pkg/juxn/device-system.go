package juxn

import (
	"fmt"
)

type SystemDevice struct {
	Vm *VM
}

func (d *SystemDevice) Input(addr byte) uint16 {
	fmt.Printf("DEI: %d\n", addr)
	return 0
}

func (d *SystemDevice) Output(addr byte, val uint16) {
	fmt.Printf("DEO: %d, %d\n", addr, val)
}
