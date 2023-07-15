package juxn

import (
	"fmt"
	"os"
)

type VM struct {
	Pc            uint16
	Halted        bool
	HaltedBecause string
	RStack        *Stack
	WStack        *Stack
	Memory        []byte
}

func NewVM() *VM {
	return &VM{
		Pc:     0,
		Halted: false,
		Memory: make([]byte, 65536),
		RStack: NewStack(),
		WStack: NewStack(),
	}
}

func (vm *VM) Run(steps int) {
	for i := 0; i < steps && !vm.Halted; i++ {
		vm.ExecuteInstruction(vm.FetchInstruction())
	}
}

func (vm *VM) ExecuteInstruction(instr Instruction) {
	switch instr.Operator {
	case BRK:
		vm.Halted = true
	case LIT:
		instr.Push(vm.FetchOperand(instr.Short))
	case INC:
		instr.Push(instr.Pop() + 1)
	case POP:
		_ = instr.Pop()
	case NIP:
		v := instr.Pop()
		_ = instr.Pop()
		instr.Push(v)
	case SWP:
		a := instr.Pop()
		b := instr.Pop()
		instr.Push(a)
		instr.Push(b)
	case ROT:
		c := instr.Pop()
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(b)
		instr.Push(c)
		instr.Push(a)
	default:
		fmt.Fprintf(os.Stderr, "Not implemented operator %02x (opcode=%02x)\n", instr.Operator, instr.Opcode)
		os.Exit(1)
	}
}

func (vm *VM) FetchInstruction() Instruction {
	instr := DecodeInstruction(vm.Memory[vm.Pc])
	instr.Vm = vm
	vm.Pc += 1
	return instr
}

func (vm *VM) FetchOperand(short bool) uint16 {
	var v uint16
	v = uint16(vm.Memory[vm.Pc])
	vm.Pc += 1
	if short {
		v = v*256 + uint16(vm.Memory[vm.Pc])
		vm.Pc += 1
	}
	return v
}
