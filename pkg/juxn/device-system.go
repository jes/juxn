package juxn

import (
	"fmt"
)

type SystemDevice struct {
	Vm *VM
}

func (d *SystemDevice) Dei(addr byte) byte {
	fmt.Printf("DEI: %d\n", addr)
	return 0
}

func (d *SystemDevice) Deo(addr byte, val byte) {
	fmt.Printf("DEO: %d, %d\n", addr, val)
}
