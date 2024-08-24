package main

import "fmt"

// 7 6 5 4 3 2 1 0
// N V _ B D I Z C
// | |   | | | | +--- Carry Flag
// | |   | | | +----- Zero Flag
// | |   | | +------- Interrupt Disable
// | |   | +--------- Decimal Mode (not used on NES)
// | |   +----------- Break Command
// | +--------------- Overflow Flag
// +----------------- Negative Flag
type CPU struct {
	registerA      uint8
	registerX      uint8
	registerY      uint8
	status         uint8
	programCounter uint16
	memory         []uint8
}

func NewCPU() *CPU {
	return &CPU{
		registerA:      0,
		registerX:      0,
		registerY:      0,
		status:         0,
		programCounter: 0,
		memory:         make([]uint8, 0xFFFF),
	}
}

// BRK - Force Interrupt
// The BRK instruction forces the generation of an interrupt request. The program
// counter and processor status are pushed on the stack then the IRQ interrupt
// vector at $FFFE/F is loaded into the PC and the break flag in the status set
// to one.
func (cpu *CPU) brk() {}

// CLC - Clear Carry Flag
// Set the carry flag to zero.
func (cpu *CPU) clc() {
	cpu.status &= uint8(0b0000_0001)
}

// CLD - Clear Decimal Mode
// Sets the decimal mode flag to zero.
func (cpu *CPU) cld() {
	cpu.status &= uint8(0b0000_1000)
}

// CLI - Clear Interrupt Disable
// Clears the interrupt disable flag allowing normal interrupt requests to be serviced.
func (cpu *CPU) cli() {
	cpu.status &= uint8(0b0000_0100)
}

// CLV - Clear Overflow Flag
// Clears the overflow flag.
func (cpu *CPU) clv() {
	cpu.status &= uint8(0b0100_0000)
}

// CMP - Compare
// This instruction compares the contents of the accumulator with another memory held value
// and sets the zero and carry flags as appropriate.
func (cpu *CPU) cmp(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	delta := cpu.registerA - value
	if cpu.registerA > value {
		cpu.status |= uint8(0b0000_0001)
	}
	cpu.updateZeroAndNegativeFlags(delta)
}

// CPX - Compare X Register
// This instruction compares the contents of the X register with another memory held value
// and sets the zero and carry flags as appropriate.
func (cpu *CPU) cpx(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	delta := cpu.registerX - value
	if cpu.registerX > value {
		cpu.status |= uint8(0b0000_0001)
	}
	cpu.updateZeroAndNegativeFlags(delta)
}

// CPY - Compare Y Register
// This instruction compares the contents of the Y register with another memory held value
// and sets the zero and carry flags as appropriate.
func (cpu *CPU) cpy(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	delta := cpu.registerY - value
	if cpu.registerY > value {
		cpu.status |= uint8(0b0000_0001)
	}
	cpu.updateZeroAndNegativeFlags(delta)
}

// INC - Increment Memory
// Adds one to the value held at a specified memory location setting the zero and negative
// flags as appropriate.
func (cpu *CPU) inc(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	value = value + uint8(1)
	cpu.memWrite(addr, value)
	cpu.updateZeroAndNegativeFlags(value)
}

// INX - Increment X Register
// Adds one to the X register setting the zero and negative flags as appropriate.
func (cpu *CPU) inx() {
	cpu.registerX += uint8(1)
	cpu.updateZeroAndNegativeFlags(cpu.registerX)
}

// INY - Increment Y Register
// Adds one to the Y register setting the zero and negative flags as appropriate.
func (cpu *CPU) iny() {
	cpu.registerY += uint8(1)
	cpu.updateZeroAndNegativeFlags(cpu.registerY)
}

// LDA - Load Accumulator
// Loads a byte of memory into the accumulator setting the zero and negative flags
// as appropriate.
func (cpu *CPU) lda(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerA = value
	cpu.updateZeroAndNegativeFlags(cpu.registerA)
}

// LDX - Load X Register
// Loads a byte of memory into the X register setting the zero and negative flags as appropriate.
func (cpu *CPU) ldx(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerX = value
	cpu.updateZeroAndNegativeFlags(cpu.registerX)
}

// LDX - Load Y Register
// Loads a byte of memory into the Y register setting the zero and negative flags as appropriate.
func (cpu *CPU) ldy(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerY = value
	cpu.updateZeroAndNegativeFlags(cpu.registerY)
}

// NOP - No Operation
// The NOP instruction causes no changes to the processor other than the normal
// incrementing of the program counter to the next instruction.
func (cpu *CPU) nop() {}

// STA - Store Accumulator
// Stores the contents of the accumulator into memory.
func (cpu *CPU) sta(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	cpu.memWrite(addr, cpu.registerA)
}

// STX - Store X Register
// Stores the contents of the X register into memory.
func (cpu *CPU) stx(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	cpu.memWrite(addr, cpu.registerX)
}

// STY - Store Y Register
// Stores the contents of the Y register into memory.
func (cpu *CPU) sty(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	cpu.memWrite(addr, cpu.registerY)
}

// TAX - Transfer Accumulator to X
// Copies the current contents of the accumulator into the X register and sets
// the zero and negative flags as appropriate.
func (cpu *CPU) tax() {
	cpu.registerX = cpu.registerA
	cpu.updateZeroAndNegativeFlags(cpu.registerX)
}

// TAY - Transfer Accumulator to Y
// Copies the current contents of the accumulator into the Y register and sets
// the zero and negative flags as appropriate.
func (cpu *CPU) tay() {
	cpu.registerY = cpu.registerA
	cpu.updateZeroAndNegativeFlags(cpu.registerY)
}

func (cpu *CPU) updateZeroAndNegativeFlags(result uint8) {
	// If the register is zero set the zero flag
	if result == 0 {
		cpu.status = cpu.status | 0b0000_0010
	} else {
		cpu.status = cpu.status | 0b1111_1101
	}

	// If the register is negative set the negative flag
	if result&0b1000_0000 != 0 {
		cpu.status = cpu.status | 0b1000_0000
	} else {
		cpu.status = cpu.status & 0b0111_1111
	}
}

func (cpu *CPU) getOperandAddress(mode AddressingMode) uint16 {
	switch mode {
	case Immediate:
		return cpu.programCounter
	case ZeroPage:
		return uint16(cpu.memRead(cpu.programCounter))
	case Absolute:
		return cpu.memReadUInt16(cpu.programCounter)
	case ZeroPageX:
		pos := cpu.memRead(cpu.programCounter)
		addr := uint16(pos + cpu.registerX)
		return addr
	case ZeroPageY:
		pos := cpu.memRead(cpu.programCounter)
		addr := uint16(pos + cpu.registerY)
		return addr
	case AbsoluteX:
		base := cpu.memReadUInt16(cpu.programCounter)
		addr := base + uint16(cpu.registerX)
		return addr
	case AbsoluteY:
		base := cpu.memReadUInt16(cpu.programCounter)
		addr := base + uint16(cpu.registerY)
		return addr
	case IndirectX:
		base := cpu.memRead(cpu.programCounter)
		ptr := uint16(base + cpu.registerX)
		lo := uint16(cpu.memRead(ptr))
		hi := uint16(cpu.memRead(ptr + 1))
		return hi<<8 | lo
	case IndirectY:
		base := uint16(cpu.memRead(cpu.programCounter))
		lo := uint16(cpu.memRead(base))
		hi := uint16(cpu.memRead(base + 1))
		derefBase := hi<<8 | lo
		deref := derefBase + uint16(cpu.registerY)
		return deref
	default:
		panic(fmt.Sprintf("AddressingMode %x is not supported", mode))

	}
}

func (cpu *CPU) memRead(addr uint16) uint8 {
	return cpu.memory[addr]
}

func (cpu *CPU) memWrite(addr uint16, data uint8) {
	cpu.memory[addr] = data
}

func (cpu *CPU) memReadUInt16(pos uint16) uint16 {
	lo := uint16(cpu.memRead(pos))
	hi := uint16(cpu.memRead(pos + 1))
	return (hi << 8) | (lo)
}

func (cpu *CPU) memWriteUInt16(pos uint16, data uint16) {
	hi := uint8(data >> 8)
	lo := uint8(data & 0xff)
	cpu.memWrite(pos, lo)
	cpu.memWrite(pos+1, hi)
}

func (cpu *CPU) reset() {
	cpu.registerA = 0
	cpu.registerX = 0
	cpu.registerY = 0
	cpu.status = 0
	cpu.programCounter = cpu.memReadUInt16(0xFFFC)
}

func (cpu *CPU) loadAndRun(program []uint8) {
	cpu.load(program)
	cpu.reset()
	cpu.run()
}

func (cpu *CPU) load(program []uint8) {
	copy(cpu.memory[0x8000:0x8000+len(program)], program[:])
	cpu.memWriteUInt16(0xFFFC, 0x8000)
}

func (cpu *CPU) run() {
	for {
		code := cpu.memRead(uint16(cpu.programCounter))
		cpu.programCounter++
		programCounterState := cpu.programCounter

		opcode, ok := CPU_OP_CODE_TABLE[code]
		if !ok {
			panic(fmt.Sprintf("Could not locate opcode in opcode table: 0x%x\n", code))
		}

		switch code {
		// BRK
		case 0x00:
			cpu.brk()
			return
		// CLC
		case 0x18:
			cpu.clc()
		// CLD
		case 0xD8:
			cpu.cld()
		// CLI
		case 0x58:
			cpu.cli()
		// CLV
		case 0xB8:
			cpu.clv()
		// CMP
		case 0xC9, 0xC5, 0xD5, 0xCD, 0xDD, 0xD9, 0xC1, 0xD1:
			cpu.cmp(opcode.AddressingMode)
		// CPX
		case 0xE0, 0xE4, 0xEC:
			cpu.cpx(opcode.AddressingMode)
		// CPY
		case 0xC0, 0xC4, 0xCC:
			cpu.cpy(opcode.AddressingMode)
		// INC
		case 0xE6, 0xF6, 0xEE, 0xFE:
			cpu.inc(opcode.AddressingMode)
		// INX
		case 0xE8:
			cpu.inx()
		// INY
		case 0xC8:
			cpu.iny()
		// LDA
		case 0xA9, 0xA5, 0xB5, 0xAD, 0xBD, 0xB9, 0xA1, 0xB1:
			cpu.lda(opcode.AddressingMode)
		// LDX
		case 0xA2, 0xA6, 0xB6, 0xAE, 0xBE:
			cpu.ldx(opcode.AddressingMode)
		// LDY
		case 0xA0, 0xA4, 0xB4, 0xAC, 0xBC:
			cpu.ldy(opcode.AddressingMode)
		// NOP
		case 0xEA:
			cpu.nop()
		// STA
		case 0x85, 0x95, 0x8D, 0x9D, 0x99, 0x81, 0x91:
			cpu.sta(opcode.AddressingMode)
		// STX
		case 0x86, 0x96, 0x8E:
			cpu.stx(opcode.AddressingMode)
		// STY
		case 0x84, 0x94, 0x8C:
			cpu.sty(opcode.AddressingMode)
		// TAX
		case 0xAA:
			cpu.tax()
		// TAY
		case 0xA8:
			cpu.tay()
		default:
			panic(fmt.Sprintf("Unsupported opcode: 0x%x\n", opcode))
		}

		if programCounterState == cpu.programCounter {
			cpu.programCounter += uint16(opcode.Bytes) - 1
		}
	}
}
