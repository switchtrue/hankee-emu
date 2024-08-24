package main

type AddressingMode int

const (
	Immediate AddressingMode = iota
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	IndirectX
	IndirectY
	NoneAddressing
)

type OpCode struct {
	Opcode         uint8
	Name           string
	AddressingMode AddressingMode
	Bytes          int
	Cycles         int
}

var CPU_OPS_CODES = map[uint8]OpCode{
	// LDA
	0xA9: {0xA9, "LDA", Immediate, 2, 2},
	0xA5: {0xA5, "LDA", ZeroPage, 2, 3},
	0xB5: {0xB5, "LDA", ZeroPageX, 2, 4},
	0xAD: {0xAD, "LDA", Absolute, 3, 4},
	0xBD: {0xBD, "LDA", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0xB9: {0xB9, "LDA", AbsoluteY, 3, 4 /* +1 if page crossed */},
	0xA1: {0xA1, "LDA", IndirectX, 2, 6},
	0xB1: {0xB1, "LDA", IndirectY, 2, 5 /* +1 if page crossed */},
	// STA
	0x85: {0x85, "STA", ZeroPage, 2, 3},
	0x95: {0x95, "STA", ZeroPageX, 2, 4},
	0x8D: {0x8D, "STA", Absolute, 3, 4},
	0x9D: {0x9D, "STA", AbsoluteX, 3, 5},
	0x99: {0x99, "STA", AbsoluteY, 3, 5},
	0x81: {0x81, "STA", IndirectX, 2, 6},
	0x91: {0x91, "STA", IndirectY, 2, 6},
	// TAX
	0xAA: {0xAA, "TAX", NoneAddressing, 1, 2},
	// INX
	0xE8: {0xE8, "INX", NoneAddressing, 1, 2},
	// NOP
	0xEA: {0xEA, "NOP", NoneAddressing, 1, 2},
	// BRK
	0x00: {0x00, "BRK", NoneAddressing, 1, 7},
}
