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
	Indirect
	IndirectX
	IndirectY
	Accumulator
	Relative
	Implied
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
	// ASL
	0x0A: {0x0A, "ASL", Accumulator, 1, 2},
	0x06: {0x06, "ASL", ZeroPage, 2, 5},
	0x16: {0x16, "ASL", ZeroPageX, 2, 6},
	0x0E: {0x0E, "ASL", Absolute, 3, 6},
	0x1E: {0x1E, "ASL", AbsoluteX, 3, 7},
	// BCC
	0x90: {0x90, "BCC", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// BCS
	0xB0: {0xB0, "BCS", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// BEQ
	0xF0: {0xF0, "BEQ", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// BIT
	0x24: {0x24, "BIT", ZeroPage, 2, 3},
	0x2C: {0x2C, "BIT", Absolute, 3, 4},
	// BMI
	0x30: {0x30, "BMI", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// BNE
	0xD0: {0xD0, "BNE", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// BPL
	0x10: {0x10, "BPL", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// BRK
	0x00: {0x00, "BRK", Implied, 1, 7},
	// BVC
	0x50: {0x50, "BVC", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// BVS
	0x70: {0x70, "BVS", Relative, 2, 2 /* +1 if branch succeeds, +2 if to a new page */},
	// CLC
	0x18: {0x18, "CLC", Implied, 1, 2},
	// CLD
	0xD8: {0xD8, "CLD", Implied, 1, 2},
	// CLI
	0x58: {0x58, "CLI", Implied, 1, 2},
	// CLV
	0xB8: {0xB8, "CLV", Implied, 1, 2},
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
	// DEC
	0xC6: {0xC6, "DEC", ZeroPage, 2, 5},
	0xD6: {0xD6, "DEC", ZeroPageX, 2, 6},
	0xCE: {0xCE, "DEC", Absolute, 3, 6},
	0xDE: {0xDE, "DEC", AbsoluteX, 3, 7},
	// DEX
	0xCA: {0xCA, "DEX", Implied, 1, 2},
	// DEY
	0x88: {0x88, "DEY", Implied, 1, 2},
	// EOR
	0x49: {0x49, "EOR", Immediate, 2, 2},
	0x45: {0x45, "EOR", ZeroPage, 2, 3},
	0x55: {0x55, "EOR", ZeroPageX, 2, 4},
	0x4D: {0x4D, "EOR", Absolute, 3, 4},
	0x5D: {0x5D, "EOR", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0x59: {0x59, "EOR", AbsoluteY, 3, 4 /* +1 if page crossed */},
	0x41: {0x41, "EOR", IndirectX, 2, 6},
	0x51: {0x51, "EOR", IndirectY, 2, 5 /* +1 if page crossed */},
	// INC
	0xE6: {0xE6, "INC", ZeroPage, 2, 5},
	0xF6: {0xF6, "INC", ZeroPageX, 2, 6},
	0xEE: {0xEE, "INC", Absolute, 3, 6},
	0xFE: {0xFE, "INC", AbsoluteX, 3, 7},
	// INX
	0xE8: {0xE8, "INX", Implied, 1, 2},
	// INY
	0xC8: {0xC8, "INY", Implied, 1, 2},
	// JMP
	0x4C: {0x4C, "JMP", Absolute, 3, 3},
	0x6C: {0x6C, "JMP", Indirect, 3, 5},
	// JSR
	0x20: {0x20, "JSR", Absolute, 3, 6},
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
	// LSR
	0x4A: {0x4A, "LDY", Accumulator, 1, 2},
	0x46: {0x46, "LDY", ZeroPage, 2, 5},
	0x56: {0x56, "LDY", ZeroPageX, 2, 6},
	0x4E: {0x4E, "LDY", Absolute, 3, 6},
	0x5E: {0x5E, "LDY", AbsoluteX, 3, 7},
	// NOP
	0xEA: {0xEA, "NOP", Implied, 1, 2},
	// ORA
	0x09: {0x09, "ORA", Immediate, 2, 2},
	0x05: {0x05, "ORA", ZeroPage, 2, 3},
	0x15: {0x15, "ORA", ZeroPageX, 2, 4},
	0x0D: {0x0D, "ORA", Absolute, 3, 4},
	0x1D: {0x1D, "ORA", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0x19: {0x19, "ORA", AbsoluteY, 3, 4 /* +1 if page crossed */},
	0x01: {0x01, "ORA", IndirectX, 2, 6},
	0x11: {0x11, "ORA", IndirectY, 2, 5 /* +1 if page crossed */},
	// PHA
	0x48: {0x48, "PHA", Implied, 1, 3},
	// PHP
	0x08: {0x08, "PHP", Implied, 1, 3},
	// PLA
	0x68: {0x68, "PLA", Implied, 1, 4},
	// PLP
	0x28: {0x28, "PLP", Implied, 1, 4},
	// ROL
	0x2A: {0x2A, "ROL", Accumulator, 1, 2},
	0x26: {0x26, "ROL", ZeroPage, 2, 5},
	0x36: {0x36, "ROL", ZeroPageX, 2, 6},
	0x2E: {0x2E, "ROL", Absolute, 3, 6},
	0x3E: {0x3E, "ROL", AbsoluteX, 3, 7},
	// ROR
	0x6A: {0x6A, "ROR", Accumulator, 1, 2},
	0x66: {0x66, "ROR", ZeroPage, 2, 5},
	0x76: {0x76, "ROR", ZeroPageX, 2, 6},
	0x6E: {0x6E, "ROR", Absolute, 3, 6},
	0x7E: {0x7E, "ROR", AbsoluteX, 3, 7},
	// RTI
	0x40: {0x40, "RTI", Implied, 1, 6},
	// RTS
	0x60: {0x60, "RTS", Implied, 1, 6},
	// SBC
	0xE9: {0xE9, "SBC", Immediate, 2, 2},
	0xE5: {0xE5, "SBC", ZeroPage, 2, 3},
	0xF5: {0xF5, "SBC", ZeroPageX, 2, 4},
	0xED: {0xED, "SBC", Absolute, 3, 4},
	0xFD: {0xFd, "SBC", AbsoluteX, 3, 4 /* +1 if page crossed */},
	0xF9: {0xF9, "SBC", AbsoluteY, 3, 4 /* +1 if page crossed */},
	0xE1: {0xE1, "SBC", IndirectX, 2, 6},
	0xF1: {0xF1, "SBC", IndirectY, 2, 5 /* +1 if page crossed */},
	// SEC
	0x38: {0x38, "SEC", Implied, 1, 2},
	// SED
	0xF8: {0xF8, "SED", Implied, 1, 2},
	// SEI
	0x78: {0x78, "SEI", Implied, 1, 2},
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
	0x84: {0x84, "STY", ZeroPage, 2, 3},
	0x94: {0x94, "STY", ZeroPageX, 2, 4},
	0x8C: {0x8C, "STY", Absolute, 3, 4},
	// TAX
	0xAA: {0xAA, "TAX", Implied, 1, 2},
	// TAY
	0xA8: {0xA8, "TAY", Implied, 1, 2},
	// TSX
	0xBA: {0xBA, "TSX", Implied, 1, 2},
	// TXA
	0x8A: {0x8A, "TXA", Implied, 1, 2},
	// TXS
	0x9A: {0x9A, "TXS", Implied, 1, 2},
	// TYA
	0x98: {0x98, "TYA", Implied, 1, 2},
}
