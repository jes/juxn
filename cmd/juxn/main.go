package main

import (
	"fmt"
	"juxn/pkg/juxn"
	"os"
)

func main() {
	vm := juxn.NewVM()
	err := vm.LoadROM("tests.rom")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	vm.Run(1000000)
	if vm.Halted {
		fmt.Printf("halted: %s\n", vm.HaltedBecause)
	}
	fmt.Println("completed!")
}
