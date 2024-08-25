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
	Relative
	NoneAddressing
)

type OpCode struct {
	Opcode         uint8
	Name           string
	AddressingMode AddressingMode
	Bytes          int
	Cycles         int
}

var CPU_OP_CODE_TABLE = map[uint8]OpCode{
	// ADC
	0x69: {0x69, "ADC", Immediate, 2, 2},
	0x65: {0x65, "ADC", ZeroPage, 2, 3},
	0x75: {0x75, "ADC", ZeroPageX, 2, 4},
	0x6D: {0x6D, "ADC", Absolute, 3, 4},
	0x7D: {0x7D, "ADC", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0x79: {0x79, "ADC", AbsoluteY, 3, 4 /* +1 if page crossed */},
	0x61: {0x61, "ADC", IndirectX, 2, 6},
	0x71: {0x71, "ADC", IndirectY, 2, 5 /* +1 if page crossed */},
	// AND
	0x29: {0x29, "AND", Immediate, 2, 2},
	0x25: {0x25, "AND", ZeroPage, 2, 3},
	0x35: {0x35, "AND", ZeroPageX, 2, 4},
	0x2D: {0x2D, "AND", Absolute, 2, 4},
	0x3D: {0x3D, "AND", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0x39: {0x39, "AND", AbsoluteY, 3, 4 /* +1 if page crossed */},
	0x21: {0x21, "AND", IndirectX, 2, 6},
	0x31: {0x31, "AND", IndirectY, 2, 5 /* +1 if page crossed */},
	// BRK
	0x00: {0x00, "BRK", NoneAddressing, 1, 7},
	// CLC
	0x18: {0x18, "CLC", NoneAddressing, 1, 2},
	// CLD
	0xD8: {0xD8, "CLD", NoneAddressing, 1, 2},
	// CLI
	0x58: {0x58, "CLI", NoneAddressing, 1, 2},
	// CLV
	0xB8: {0xB8, "CLV", NoneAddressing, 1, 2},
	// CMP
	0xC9: {0xC9, "CMP", Immediate, 2, 2},
	0xC5: {0xC5, "CMP", ZeroPage, 2, 3},
	0xD5: {0xD5, "CMP", ZeroPageX, 2, 4},
	0xCD: {0xCD, "CMP", ZeroPageX, 3, 4},
	0xDD: {0xDD, "CMP", Absolute, 3, 4 /* +1 if page crossed */},
	0xD9: {0xD9, "CMP", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0xC1: {0xC1, "CMP", IndirectX, 2, 6},
	0xD1: {0xD1, "CMP", IndirectY, 2, 5 /* +1 if page crossed */},
	// CPX
	0xE0: {0xE0, "CPX", Immediate, 2, 2},
	0xE4: {0xE4, "CPX", ZeroPage, 2, 3},
	0xEC: {0xEC, "CPX", Absolute, 3, 4},
	// CPY
	0xC0: {0xC0, "CPY", Immediate, 2, 2},
	0xC4: {0xC4, "CPY", ZeroPage, 2, 3},
	0xCC: {0xCC, "CPY", Absolute, 3, 4},
	// INC
	0xE6: {0xE6, "INC", ZeroPage, 2, 5},
	0xF6: {0xF6, "INC", ZeroPageX, 2, 6},
	0xEE: {0xEE, "INC", Absolute, 3, 6},
	0xFE: {0xFE, "INC", AbsoluteX, 3, 7},
	// INX
	0xE8: {0xE8, "INX", NoneAddressing, 1, 2},
	// INY
	0xC8: {0xC8, "INY", NoneAddressing, 1, 2},
	// LDA
	0xA9: {0xA9, "LDA", Immediate, 2, 2},
	0xA5: {0xA5, "LDA", ZeroPage, 2, 3},
	0xB5: {0xB5, "LDA", ZeroPageX, 2, 4},
	0xAD: {0xAD, "LDA", Absolute, 3, 4},
	0xBD: {0xBD, "LDA", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0xB9: {0xB9, "LDA", AbsoluteY, 3, 4 /* +1 if page crossed */},
	0xA1: {0xA1, "LDA", IndirectX, 2, 6},
	0xB1: {0xB1, "LDA", IndirectY, 2, 5 /* +1 if page crossed */},
	// LDX
	0xA2: {0xA2, "LDX", Immediate, 2, 2},
	0xA6: {0xA6, "LDX", ZeroPage, 2, 3},
	0xB6: {0xB6, "LDX", ZeroPageY, 2, 4},
	0xAE: {0xAE, "LDX", Absolute, 3, 4},
	0xBE: {0xBE, "LDX", AbsoluteY, 3, 4 /* +1 if page crossed */},
	// LDY
	0xA0: {0xA0, "LDY", Immediate, 2, 2},
	0xA4: {0xA4, "LDY", ZeroPage, 2, 3},
	0xB4: {0xB4, "LDY", ZeroPageX, 2, 4},
	0xAC: {0xAC, "LDY", ZeroPageX, 3, 4},
	0xBC: {0xBC, "LDY", AbsoluteX, 3, 4 /* +1 if page crossed */},
	// NOP
	0xEA: {0xEA, "NOP", NoneAddressing, 1, 2},
	// STA
	0x85: {0x85, "STA", ZeroPage, 2, 3},
	0x95: {0x95, "STA", ZeroPageX, 2, 4},
	0x8D: {0x8D, "STA", Absolute, 3, 4},
	0x9D: {0x9D, "STA", AbsoluteX, 3, 5},
	0x99: {0x99, "STA", AbsoluteY, 3, 5},
	0x81: {0x81, "STA", IndirectX, 2, 6},
	0x91: {0x91, "STA", IndirectY, 2, 6},
	// STX
	0x86: {0x86, "STX", ZeroPage, 2, 3},
	0x96: {0x96, "STX", ZeroPageY, 2, 4},
	0x8E: {0x8E, "STX", Absolute, 3, 4},
	// STY
	0x84: {0x84, "STX", ZeroPage, 2, 3},
	0x94: {0x94, "STX", ZeroPageX, 2, 4},
	0x8C: {0x8C, "STX", Absolute, 3, 4},
	// TAX
	0xAA: {0xAA, "TAX", NoneAddressing, 1, 2},
	// TAY
	0xA8: {0xA8, "TAX", NoneAddressing, 1, 2},
}
