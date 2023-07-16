package juxn

import (
	"fmt"
	"os"
)

type Device interface {
	Input(addr byte) byte
	Output(addr byte, val byte)
}

type VM struct {
	Pc            uint16
	Halted        bool
	HaltedBecause string
	RStack        *Stack
	WStack        *Stack
	Memory        []byte
	DevPage       []byte
	Dev           []Device
}

func NewVM() *VM {
	v := &VM{
		Pc:      0x100,
		Halted:  false,
		Memory:  make([]byte, 65536),
		RStack:  NewStack(),
		WStack:  NewStack(),
		DevPage: make([]byte, 256),
		Dev:     make([]Device, 16),
	}
	v.SetDevice(0x00, &SystemDevice{Vm: v})
	v.SetDevice(0x10, &ConsoleDevice{Vm: v})
	return v
}

func (vm *VM) SetDevice(addr byte, dev Device) {
	vm.Dev[addr>>8] = dev
}

func (vm *VM) Run(steps int) {
	for i := 0; i < steps && !vm.Halted; i++ {
		vm.ExecuteInstruction(vm.FetchInstruction())
	}
}

func (vm *VM) LoadROM(filename string) error {
	rom, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if len(rom) > 65536 {
		return fmt.Errorf("%s: too large to fit in memory (%d > 65536 bytes)", filename, len(rom))
	}
	copy(vm.Memory[0x100:], rom)
	return nil
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
			vm.Pc = uint16(int32(vm.Pc) + int32(int8(instr.Pop())))
		}
	case JCN:
		addr := instr.Pop()
		cond := instr.PopByte()
		if cond != 0 {
			if instr.Short {
				vm.Pc = addr
			} else {
				vm.Pc = uint16(int32(vm.Pc) + int32(int8(addr)))
			}
		}
	case JSR:
		vm.RStack.PushShort(vm.Pc)
		if instr.Short {
			vm.Pc = instr.Pop()
		} else {
			vm.Pc = uint16(int32(vm.Pc) + int32(int8(instr.Pop())))
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
		instr.Push(vm.ReadMemory(uint16(int32(vm.Pc)+int32(int8(reladdr))), instr.Short))
	case STR:
		reladdr := instr.PopByte()
		val := instr.Pop()
		vm.WriteMemory(uint16(int32(vm.Pc)+int32(int8(reladdr))), val, instr.Short)
	case LDA:
		addr := instr.PopShort()
		instr.Push(vm.ReadMemory(addr, instr.Short))
	case STA:
		addr := instr.PopShort()
		val := instr.Pop()
		vm.WriteMemory(addr, val, instr.Short)
	case DEI:
		devaddr := instr.PopByte()
		d := vm.Dev[devaddr>>8]
		if d != nil {
			instr.Push(uint16(d.Input(devaddr)))
		} else {
			instr.Push(uint16(vm.DevPage[devaddr]))
		}
	case DEO:
		devaddr := instr.PopByte()
		val := instr.Pop()
		vm.DevPage[devaddr] = byte(val)
		d := vm.Dev[devaddr>>8]
		if d != nil {
			d.Output(devaddr, byte(val))
		}
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
		if b == 0 {
			vm.Halt("divide by 0")
		} else {
			instr.Push(a / b)
		}
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
		bitsLeft := (shift & 0xf0) >> 4
		bitsRight := (shift & 0x0f)
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
	v := uint16(vm.Memory[vm.Pc])
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
