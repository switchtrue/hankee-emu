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

// ADC - Add with Carry
// This instruction adds the contents of a memory location to the accumulator together
// with the carry bit. If overflow occurs the carry bit is set, this enables multiple
// byte addition to be performed.
func (cpu *CPU) adc(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)

	sum := cpu.registerA + value
	if cpu.getFlagCarry() {
		sum += 1
	}

	overflow := (cpu.registerA^value)&0x80 == 0 && (cpu.registerA^sum)&0x80 != 0

	cpu.registerA += sum

	cpu.setFlagOverflow(overflow)
	cpu.setFlagCarryForResult(cpu.registerA)
	cpu.setFlagZeroAndNegativeForResult(cpu.registerA)
}

// AND - Logical AND
// A logical AND is performed, bit by bit, on the accumulator contents using the contents
// of a byte of memory.
func (cpu *CPU) and(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)

	cpu.registerA &= value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerA)
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
	cpu.setFlagCarry(false)
}

// CLD - Clear Decimal Mode
// Sets the decimal mode flag to zero.
func (cpu *CPU) cld() {
	cpu.setFlagDecimalMode(false)
}

// CLI - Clear Interrupt Disable
// Clears the interrupt disable flag allowing normal interrupt requests to be serviced.
func (cpu *CPU) cli() {
	cpu.setFlagInterruptDisable(false)
}

// CLV - Clear Overflow Flag
// Clears the overflow flag.
func (cpu *CPU) clv() {
	cpu.setFlagOverflow(false)
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
	cpu.setFlagZeroAndNegativeForResult(delta)
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
	cpu.setFlagZeroAndNegativeForResult(delta)
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
	cpu.setFlagZeroAndNegativeForResult(delta)
}

// INC - Increment Memory
// Adds one to the value held at a specified memory location setting the zero and negative
// flags as appropriate.
func (cpu *CPU) inc(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	value = value + uint8(1)
	cpu.memWrite(addr, value)
	cpu.setFlagZeroAndNegativeForResult(value)
}

// INX - Increment X Register
// Adds one to the X register setting the zero and negative flags as appropriate.
func (cpu *CPU) inx() {
	cpu.registerX += uint8(1)
	cpu.setFlagZeroAndNegativeForResult(cpu.registerX)
}

// INY - Increment Y Register
// Adds one to the Y register setting the zero and negative flags as appropriate.
func (cpu *CPU) iny() {
	cpu.registerY += uint8(1)
	cpu.setFlagZeroAndNegativeForResult(cpu.registerY)
}

// LDA - Load Accumulator
// Loads a byte of memory into the accumulator setting the zero and negative flags
// as appropriate.
func (cpu *CPU) lda(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerA = value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerA)
}

// LDX - Load X Register
// Loads a byte of memory into the X register setting the zero and negative flags as appropriate.
func (cpu *CPU) ldx(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerX = value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerX)
}

// LDX - Load Y Register
// Loads a byte of memory into the Y register setting the zero and negative flags as appropriate.
func (cpu *CPU) ldy(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerY = value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerY)
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
	cpu.setFlagZeroAndNegativeForResult(cpu.registerX)
}

// TAY - Transfer Accumulator to Y
// Copies the current contents of the accumulator into the Y register and sets
// the zero and negative flags as appropriate.
func (cpu *CPU) tay() {
	cpu.registerY = cpu.registerA
	cpu.setFlagZeroAndNegativeForResult(cpu.registerY)
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

		switch opcode.Name {
		// ADC
		case "ADC":
			cpu.adc(opcode.AddressingMode)
		// AND
		case "AND":
			cpu.and(opcode.AddressingMode)
		// BRK
		case "BRK":
			cpu.brk()
			return
		// CLC
		case "CLC":
			cpu.clc()
		// CLD
		case "CLD":
			cpu.cld()
		// CLI
		case "CLI":
			cpu.cli()
		// CLV
		case "CLV":
			cpu.clv()
		// CMP
		case "CMP":
			cpu.cmp(opcode.AddressingMode)
		// CPX
		case "CPX":
			cpu.cpx(opcode.AddressingMode)
		// CPY
		case "CPY":
			cpu.cpy(opcode.AddressingMode)
		// INC
		case "INC":
			cpu.inc(opcode.AddressingMode)
		// INX
		case "INX":
			cpu.inx()
		// INY
		case "INY":
			cpu.iny()
		// LDA
		case "LDA":
			cpu.lda(opcode.AddressingMode)
		// LDX
		case "LDX":
			cpu.ldx(opcode.AddressingMode)
		// LDY
		case "LDY":
			cpu.ldy(opcode.AddressingMode)
		// NOP
		case "NOP":
			cpu.nop()
		// STA
		case "STA":
			cpu.sta(opcode.AddressingMode)
		// STX
		case "STX":
			cpu.stx(opcode.AddressingMode)
		// STY
		case "STY":
			cpu.sty(opcode.AddressingMode)
		// TAX
		case "TAX":
			cpu.tax()
		// TAY
		case "TAY":
			cpu.tay()
		default:
			panic(fmt.Sprintf("Unsupported opcode: 0x%x\n", opcode))
		}

		if programCounterState == cpu.programCounter {
			cpu.programCounter += uint16(opcode.Bytes) - 1
		}
	}
}
