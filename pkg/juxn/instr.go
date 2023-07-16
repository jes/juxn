package juxn

import "fmt"

const (
	BRK = 0x00
	INC = 0x01
	POP = 0x02
	NIP = 0x03
	SWP = 0x04
	ROT = 0x05
	DUP = 0x06
	OVR = 0x07
	EQU = 0x08
	NEQ = 0x09
	GTH = 0x0a
	LTH = 0x0b
	JMP = 0x0c
	JCN = 0x0d
	JSR = 0x0e
	STH = 0x0f
	LDZ = 0x10
	STZ = 0x11
	LDR = 0x12
	STR = 0x13
	LDA = 0x14
	STA = 0x15
	DEI = 0x16
	DEO = 0x17
	ADD = 0x18
	SUB = 0x19
	MUL = 0x1a
	DIV = 0x1b
	AND = 0x1c
	ORA = 0x1d
	EOR = 0x1e
	SFT = 0x1f

	JCI = 0x20
	JMI = 0x40
	JSI = 0x60
	LIT = 0x80
)

type Instruction struct {
	Opcode   byte
	Operator byte
	Short    bool // operate on shorts instead of bytes
	Return   bool // operate on the return stack
	Keep     bool // operate without consuming items
	Vm       *VM
}

func DecodeInstruction(opcode byte) Instruction {
	i := Instruction{
		Opcode:   opcode,
		Operator: opcode & 0x1f,
		Short:    (opcode & 0x20) != 0,
		Return:   (opcode & 0x40) != 0,
		Keep:     (opcode & 0x80) != 0,
	}
	if i.Operator == BRK && i.Opcode != BRK {
		i.Keep = false
		switch i.Opcode {
		case JCI, JMI, JSI:
			i.Operator = i.Opcode
			i.Short = false
			i.Return = false
		default:
			i.Operator = LIT
		}
	}
	return i
}

func (i Instruction) Push(val uint16) {
	stk := i.Vm.WStack
	if i.Return {
		stk = i.Vm.RStack
	}
	ok := stk.Push(val, i.Short)
	if !ok {
		i.Vm.Halt("stack overflow")
	}
}

func (i Instruction) PushByte(val byte) {
	wasShort := i.Short
	i.Short = false
	i.Push(uint16(val))
	i.Short = wasShort
}

func (i Instruction) PushShort(val uint16) {
	wasShort := i.Short
	i.Short = true
	i.Push(val)
	i.Short = wasShort
}

func (i Instruction) Pop() uint16 {
	stk := i.Vm.WStack
	if i.Return {
		stk = i.Vm.RStack
	}
	var v uint16
	var ok bool
	if i.Keep {
		v, ok = stk.Peek(i.Short)
	} else {
		v, ok = stk.Pop(i.Short)
	}
	if !ok {
		i.Vm.Halt("stack underflow")
	}
	return v
}

func (i Instruction) PopByte() byte {
	wasShort := i.Short
	i.Short = false
	v := byte(i.Pop())
	i.Short = wasShort
	return v
}

func (i Instruction) PopShort() uint16 {
	wasShort := i.Short
	i.Short = true
	v := i.Pop()
	i.Short = wasShort
	return v
}

func (i Instruction) String() string {
	var op string
	switch i.Operator {
	case BRK:
		op = "BRK"
	case INC:
		op = "INC"
	case POP:
		op = "POP"
	case NIP:
		op = "NIP"
	case SWP:
		op = "SWP"
	case ROT:
		op = "ROT"
	case DUP:
		op = "DUP"
	case OVR:
		op = "OVR"
	case EQU:
		op = "EQU"
	case NEQ:
		op = "NEQ"
	case GTH:
		op = "GTH"
	case LTH:
		op = "LTH"
	case JMP:
		op = "JMP"
	case JCN:
		op = "JCN"
	case JSR:
		op = "JSR"
	case STH:
		op = "STH"
	case LDZ:
		op = "LDZ"
	case STZ:
		op = "STZ"
	case LDR:
		op = "LDR"
	case STR:
		op = "STR"
	case LDA:
		op = "LDA"
	case STA:
		op = "STA"
	case DEI:
		op = "DEI"
	case DEO:
		op = "DEO"
	case ADD:
		op = "ADD"
	case SUB:
		op = "SUB"
	case MUL:
		op = "MUL"
	case DIV:
		op = "DIV"
	case AND:
		op = "AND"
	case ORA:
		op = "ORA"
	case EOR:
		op = "EOR"
	case SFT:
		op = "SFT"
	case LIT:
		op = "LIT"
	case JCI:
		op = "JCI"
	case JMI:
		op = "JMI"
	case JSI:
		op = "JSI"
	default:
		op = fmt.Sprintf("<0x%02x>", i.Opcode)
	}
	if i.Short {
		op += "2"
	}
	if i.Keep {
		op += "k"
	}
	if i.Return {
		op += "r"
	}
	return op
}
