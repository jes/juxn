package main

import (
	"fmt"
	"juxn/pkg/juxn"
)

func main() {
	vm := juxn.NewVM()
	vm.Memory[0] = 0x80
	vm.Memory[1] = 0x12
	vm.Run(100)
	fmt.Println("completed!")
	fmt.Println("Memory = %v", vm.WStack)
}
