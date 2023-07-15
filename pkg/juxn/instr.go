package juxn

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
		Opcode: opcode,
	}
	i.Operator = opcode & 0x1f
	i.Short = (opcode & 0x20) != 0
	i.Return = (opcode & 0x40) != 0
	i.Keep = (opcode & 0x80) != 0
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
		i.Vm.Halted = true
		i.Vm.HaltedBecause = "push to full stack"
	}
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
		i.Vm.Halted = true
		i.Vm.HaltedBecause = "pop from empty stack"
	}
	return v
}
