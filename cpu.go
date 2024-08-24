package main

import "fmt"

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

// LDA - Load Accumulator
// Loads a byte of memory into the accumulator setting the zero and negative flags
// as appropriate.
func (cpu *CPU) lda(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerA = value
	cpu.updateZeroAndNegativeFlags(cpu.registerA)
}

// TAX - Transfer Accumulator to X
// Copies the current contents of the accumulator into the X register and sets
// the zero and negative flags as appropriate.
func (cpu *CPU) tax() {
	cpu.registerX = cpu.registerA
	cpu.updateZeroAndNegativeFlags(cpu.registerX)
}

// INX - Increment X Register
// Adds one to the X register setting the zero and negative flags as appropriate.
func (cpu *CPU) inx() {
	cpu.registerX += uint8(1)
	cpu.updateZeroAndNegativeFlags(cpu.registerX)
}

// STA - Store Accumulator
// Stores the contents of the accumulator into memory.
func (cpu *CPU) sta(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	cpu.memWrite(addr, cpu.registerA)
}

// NOP - No Operation
// The NOP instruction causes no changes to the processor other than the normal
// incrementing of the program counter to the next instruction.
func (cpu *CPU) nop() {}

// BRK - Force Interrupt
// The BRK instruction forces the generation of an interrupt request. The program
// counter and processor status are pushed on the stack then the IRQ interrupt
// vector at $FFFE/F is loaded into the PC and the break flag in the status set
// to one.
func (cpu *CPU) brk() {}

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

		opcode := CPU_OPS_CODES[code]

		switch code {
		// LDA
		case 0xA9, 0xA5, 0xB5, 0xAD, 0xBD, 0xB9, 0xA1, 0xB1:
			cpu.lda(opcode.AddressingMode)
		// STA
		case 0x85, 0x95, 0x8D, 0x9D, 0x99, 0x81, 0x91:
			cpu.sta(IndirectY)
		// TAX
		case 0xAA:
			cpu.tax()
		// INX
		case 0xE8:
			cpu.inx()
		// NOP
		case 0xEA:
			cpu.nop()
		// BRK
		case 0x00:
			cpu.brk()
			return
		default:
			panic(fmt.Sprintf("Unsupported opcode: 0x%x\n", opcode))
		}

		if programCounterState == cpu.programCounter {
			cpu.programCounter += uint16(opcode.Bytes) - 1
		}
	}
}
