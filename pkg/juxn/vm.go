package juxn

import (
	"fmt"
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

func (vm *VM) Halt(reason string) {
	vm.Halted = true
	vm.HaltedBecause = reason
}

func (vm *VM) ExecuteInstruction(instr Instruction) {
	switch instr.Operator {
	case BRK:
		vm.Halt("BRK instruction")
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
	case DUP:
		v := instr.Pop()
		instr.Push(v)
		instr.Push(v)
	case OVR:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a)
		instr.Push(b)
		instr.Push(a)
	case EQU, NEQ:
		b := instr.Pop()
		a := instr.Pop()
		if (a == b) == (instr.Operator == EQU) {
			instr.PushByte(1)
		} else {
			instr.PushByte(0)
		}
	case GTH, LTH:
		b := instr.Pop()
		a := instr.Pop()
		if a != b && ((a > b) == (instr.Operator == GTH)) {
			instr.PushByte(1)
		} else {
			instr.PushByte(0)
		}
	case JMP:
		if instr.Short {
			vm.Pc = instr.Pop()
		} else {
			// TODO: signed
			vm.Pc += instr.Pop()
		}
	case JCN:
		addr := instr.Pop()
		cond := instr.PopByte()
		if cond != 0 {
			if instr.Short {
				vm.Pc = addr
			} else {
				// TODO: signed
				vm.Pc += addr
			}
		}
	case JSR:
		vm.RStack.PushShort(vm.Pc)
		if instr.Short {
			vm.Pc = instr.Pop()
		} else {
			// TODO: signed
			vm.Pc += instr.Pop()
		}
	case STH:
		vm.RStack.Push(instr.Pop(), instr.Short)
	case LDZ:
		instr.Push(vm.ReadMemory(uint16(instr.PopByte()), instr.Short))
	case STZ:
		addr := instr.PopByte()
		val := instr.Pop()
		vm.WriteMemory(uint16(addr), val, instr.Short)
	case LDR:
		reladdr := instr.PopByte()
		// TODO: signed
		instr.Push(vm.ReadMemory(vm.Pc+uint16(reladdr), instr.Short))
	case STR:
		reladdr := instr.PopByte()
		val := instr.Pop()
		vm.WriteMemory(vm.Pc+uint16(reladdr), val, instr.Short)
	case LDA:
		addr := instr.PopShort()
		instr.Push(vm.ReadMemory(addr, instr.Short))
	case STA:
		addr := instr.PopShort()
		val := instr.Pop()
		vm.WriteMemory(addr, val, instr.Short)
	case DEI:
		dev := instr.PopByte()
		_ = dev
		vm.Halt("DEI not implemented")
	case DEO:
		dev := instr.PopByte()
		val := instr.Pop()
		_ = dev
		_ = val
		vm.Halt("DEO not implemented")
	case ADD:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a + b)
	case SUB:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a - b)
	case MUL:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a * b)
	case DIV:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a / b)
	case AND:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a & b)
	case ORA:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a | b)
	case EOR:
		b := instr.Pop()
		a := instr.Pop()
		instr.Push(a ^ b)
	case SFT:
		shift := instr.PopByte()
		a := instr.Pop()
		bitsLeft := (shift & 0x80) >> 4
		bitsRight := (shift & 0x08)
		instr.Push((a >> bitsRight) << bitsLeft)
	default:
		vm.Halt(fmt.Sprintf("Not implemented operator %02x (opcode=%02x)\n", instr.Operator, instr.Opcode))
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

func (vm *VM) ReadMemory(addr uint16, short bool) uint16 {
	if short {
		return 256*uint16(vm.Memory[addr]) + uint16(vm.Memory[addr+1])
	} else {
		return uint16(vm.Memory[addr])
	}
}

func (vm *VM) WriteMemory(addr uint16, val uint16, short bool) {
	if short {
		vm.Memory[addr] = byte(val >> 8)
		vm.Memory[addr+1] = byte(val & 0xff)
	} else {
		vm.Memory[addr] = byte(val)
	}
}
