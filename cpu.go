package main

import "fmt"

type CPU struct {
	registerA      uint8
	registerX      uint8
	registerY      uint8
	status         uint8
	programCounter uint16
	stackPointer   uint8
	memory         []uint8
}

func NewCPU() *CPU {
	return &CPU{
		registerA:      0,
		registerX:      0,
		registerY:      0,
		status:         0,
		programCounter: 0,
		stackPointer:   STACK_RESET,
		memory:         make([]uint8, 0xFFFF),
	}
}

// ADC - Add with Carry
// This instruction adds the contents of a memory location to the accumulator
// together with the carry bit. If overflow occurs the carry bit is set, this
// enables multiple byte addition to be performed.
func (cpu *CPU) adc(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.addToRegisterA(value)
}

// AND - Logical AND
// A logical AND is performed, bit by bit, on the accumulator contents using the
// contents of a byte of memory.
func (cpu *CPU) and(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)

	cpu.registerA &= value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerA)
}

// ASL - Arithmetic Shift Left
// This operation shifts all the bits of the accumulator or memory contents one
// bit left. Bit 0 is set to 0 and bit 7 is placed in the carry flag. The effect
// of this operation is to multiply the memory contents by 2 (ignoring 2's
// complement considerations), setting the carry if the result will not fit in
// 8 bits.
func (cpu *CPU) asl(mode AddressingMode) {
	cpu.setFlagCarry(cpu.registerA>>7 == 1)
	cpu.registerA = cpu.registerA << 1
}

// BCC - Branch if Carry Clear
// If the carry flag is clear then add the relative displacement to the program
// counter to cause a branch to a new location.
func (cpu *CPU) bcc() {
	cpu.branch(!cpu.getFlagCarry())
}

// BCS - Branch if Carry Set
// If the carry flag is set then add the relative displacement to the program
// counter to cause a branch to a new location.
func (cpu *CPU) bcs() {
	cpu.branch(cpu.getFlagCarry())
}

// BEQ - Branch if Equal
// If the zero flag is set then add the relative displacement to the program
// counter to cause a branch to a new location.
func (cpu *CPU) beq() {
	cpu.branch(!cpu.getFlagZero())
}

// BIT - Bit Test
// This instructions is used to test if one or more bits are set in a target
// memory location. The mask pattern in A is ANDed with the value in memory to
// set or clear the zero flag, but the result is not kept. Bits 7 and 6 of the
// value from memory are copied into the N and V flags.
func (cpu *CPU) bit(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	anded := cpu.registerA & value
	cpu.setFlagZero(anded == 0)
	cpu.setFlagNegative(value&0b10000000 > 0)
	cpu.setFlagOverflow(value&0b01000000 > 0)
}

// BMI - Branch if Minus
// If the negative flag is set then add the relative displacement to the program
// counter to cause a branch to a new location.
func (cpu *CPU) bmi() {
	cpu.branch(cpu.getFlagNegative())
}

// BNE - Branch if Not Equal
// If the zero flag is clear then add the relative displacement to the program
// counter to cause a branch to a new location.
func (cpu *CPU) bne() {
	cpu.branch(!cpu.getFlagZero())
}

// BPL - Branch if Positive
// If the negative flag is clear then add the relative displacement to the
// program counter to cause a branch to a new location.
func (cpu *CPU) bpl() {
	cpu.branch(!cpu.getFlagNegative())
}

// BRK - Force Interrupt
// The BRK instruction forces the generation of an interrupt request. The
// program counter and processor status are pushed on the stack then the IRQ
// interrupt vector at $FFFE/F is loaded into the PC and the break flag in the
// status set to one.
func (cpu *CPU) brk() {}

// BVC - Branch if Overflow Clear
// If the overflow flag is clear then add the relative displacement to the
// program counter to cause a branch to a new location.
func (cpu *CPU) bvc() {
	cpu.branch(!cpu.getFlagOverflow())
}

// BVS - Branch if Overflow Set
// If the overflow flag is set then add the relative displacement to the program
// counter to cause a branch to a new location.
func (cpu *CPU) bvs() {
	cpu.branch(cpu.getFlagOverflow())
}

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
// Clears the interrupt disable flag allowing normal interrupt requests to be
// serviced.
func (cpu *CPU) cli() {
	cpu.setFlagInterruptDisable(false)
}

// CLV - Clear Overflow Flag
// Clears the overflow flag.
func (cpu *CPU) clv() {
	cpu.setFlagOverflow(false)
}

// CMP - Compare
// This instruction compares the contents of the accumulator with another memory
// held value and sets the zero and carry flags as appropriate.
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
// This instruction compares the contents of the X register with another memory
// held value and sets the zero and carry flags as appropriate.
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
// This instruction compares the contents of the Y register with another memory
// held value and sets the zero and carry flags as appropriate.
func (cpu *CPU) cpy(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	delta := cpu.registerY - value
	if cpu.registerY > value {
		cpu.status |= uint8(0b0000_0001)
	}
	cpu.setFlagZeroAndNegativeForResult(delta)
}

// DEC - Decrement Memory
// Subtracts one from the value held at a specified memory location setting the
// zero and negative flags as appropriate.
func (cpu *CPU) dec(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	value = value - uint8(1)
	cpu.memWrite(addr, value)
	cpu.setFlagZeroAndNegativeForResult(value)
}

// DEX - Decrement X Register
// Subtracts one from the X register setting the zero and negative flags as
// appropriate.
func (cpu *CPU) dex(mode AddressingMode) {
	cpu.registerX -= uint8(1)
	cpu.setFlagZeroAndNegativeForResult(cpu.registerX)
}

// DEY - Decrement Y Register
// Subtracts one from the Y register setting the zero and negative flags as
// appropriate.
func (cpu *CPU) dey(mode AddressingMode) {
	cpu.registerY -= uint8(1)
	cpu.setFlagZeroAndNegativeForResult(cpu.registerY)
}

// EOR - Exclusive OR
// An exclusive OR is performed, bit by bit, on the accumulator contents using
// the contents of a byte of memory.
func (cpu *CPU) eor(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerA ^= value
}

// INC - Increment Memory
// Adds one to the value held at a specified memory location setting the zero
// and negative flags as appropriate.
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

// JMP - Jump
// Sets the program counter to the address specified by the operand.
func (cpu *CPU) jmp(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memReadUInt16(addr)
	cpu.programCounter = value
}

// JSR - Jump to Subroutine
// The JSR instruction pushes the address (minus one) of the return point on to
// the stack and then sets the program counter to the target memory address.
func (cpu *CPU) jsr() {
	cpu.stackPushUInt16(cpu.programCounter + 2 - 1)
	addr := cpu.memReadUInt16(cpu.programCounter)
	cpu.programCounter = addr
}

// LDA - Load Accumulator
// Loads a byte of memory into the accumulator setting the zero and negative
// flags as appropriate.
func (cpu *CPU) lda(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerA = value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerA)
}

// LDX - Load X Register
// Loads a byte of memory into the X register setting the zero and negative
// flags as appropriate.
func (cpu *CPU) ldx(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerX = value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerX)
}

// LDX - Load Y Register
// Loads a byte of memory into the Y register setting the zero and negative
// flags as appropriate.
func (cpu *CPU) ldy(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerY = value
	cpu.setFlagZeroAndNegativeForResult(cpu.registerY)
}

// LSR - Logical Shift Right
// Each of the bits in A or M is shift one place to the right. The bit that was
// in bit 0 is shifted into the carry flag. Bit 7 is set to zero.
func (cpu *CPU) lsr(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	result := value >> 1
	cpu.setFlagCarry(value&1 == 1)
	cpu.memWrite(addr, result)
	cpu.setFlagZeroAndNegativeForResult(result)
}

// NOP - No Operation
// The NOP instruction causes no changes to the processor other than the normal
// incrementing of the program counter to the next instruction.
func (cpu *CPU) nop() {}

// ORA - Logical Inclusive OR
// inclusive OR is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
func (cpu *CPU) ora(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.registerA |= value
}

// PHA - Push Accumulator
// Pushes a copy of the accumulator on to the stack.
func (cpu *CPU) pha() {
	cpu.stackPush(cpu.registerA)
}

// PHP - Push Processor Status
// Pushes a copy of the status flags on to the stack.
func (cpu *CPU) php() {
	cpu.stackPush(cpu.status)
}

// PLA - Pull Accumulator
// Pulls an 8 bit value from the stack and into the accumulator. The zero and
// negative flags are set as appropriate.
func (cpu *CPU) pla() {
	cpu.registerA = cpu.stackPop()
	cpu.setFlagNegativeForResult(cpu.registerA)
}

// PLP - Pull Processor Status
// Pulls an 8 bit value from the stack and into the processor flags. The flags
// will take on new states as determined by the value pulled.
func (cpu *CPU) plp() {
	cpu.status = cpu.stackPop()
}

// ROL - Rotate Left
// Move each of the bits in either A or M one place to the left. Bit 0 is filled
// with the current value of the carry flag whilst the old bit 7 becomes the
// new carry flag value.
func (cpu *CPU) rol(mode AddressingMode) {
	var value uint8
	if mode == Accumulator {
		value = cpu.registerA
	} else {
		addr := cpu.getOperandAddress(mode)
		value = cpu.memRead(addr)
	}

	oldCarry := cpu.getFlagCarry()
	cpu.setFlagCarry(value>>7 == 1)
	result := value << 1
	if oldCarry {
		result = result | 1
	}

	if mode == Accumulator {
		cpu.registerA = result
	} else {
		addr := cpu.getOperandAddress(mode)
		cpu.memWrite(addr, result)
		cpu.setFlagNegativeForResult(result)
	}
}

// ROR - Rotate Right
// Move each of the bits in either A or M one place to the right. Bit 7 is
// filled with the current value of the carry flag whilst the old bit 0
// becomes the new carry flag value.
func (cpu *CPU) ror(mode AddressingMode) {
	var value uint8
	if mode == Accumulator {
		value = cpu.registerA
	} else {
		addr := cpu.getOperandAddress(mode)
		value = cpu.memRead(addr)
	}

	oldCarry := cpu.getFlagCarry()
	cpu.setFlagCarry(value&1 == 1)
	result := value >> 1
	if oldCarry {
		result = result | 0b10000000
	}

	if mode == Accumulator {
		cpu.registerA = result
	} else {
		addr := cpu.getOperandAddress(mode)
		cpu.memWrite(addr, result)
		cpu.setFlagNegativeForResult(result)
	}
}

// RTI - Return from Interrupt
// The RTI instruction is used at the end of an interrupt processing routine.
// It pulls the processor flags from the stack followed by the program counter.
func (cpu *CPU) rti() {
	cpu.status = cpu.stackPop()
	cpu.programCounter = cpu.stackPopUInt16()
}

// RTS - Return from Subroutine
// The RTS instruction is used at the end of a subroutine to return to the
// calling routine. It pulls the program counter (minus one) from the stack.
func (cpu *CPU) rts() {
	cpu.programCounter = cpu.stackPopUInt16() + 1
}

// SBC - Subtract with Carry
// This instruction subtracts the contents of a memory location to the
// accumulator together with the not of the carry bit. If overflow occurs the
// carry bit is clear, this enables multiple byte subtraction to be performed.
func (cpu *CPU) sbc(mode AddressingMode) {
	addr := cpu.getOperandAddress(mode)
	value := cpu.memRead(addr)
	cpu.addToRegisterA(-value - uint8(1))
}

// SEC - Set Carry Flag
// Set the carry flag to one.
func (cpu *CPU) sec() {
	cpu.setFlagCarry(true)
}

// SED - Set Decimal Flag
// Set the decimal mode flag to one.
func (cpu *CPU) sed() {
	cpu.setFlagDecimalMode(true)
}

// SEI - Set Interrupt Disable
// Set the interrupt disable flag to one.
func (cpu *CPU) sei() {
	cpu.setFlagInterruptDisable(true)
}

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

// TSX - Transfer Stack Pointer to X
// Copies the current contents of the stack register into the X register and
// sets the zero and negative flags as appropriate.
func (cpu *CPU) tsx() {
	cpu.registerX = cpu.stackPop()
	cpu.setFlagZeroAndNegativeForResult(cpu.registerX)
}

// TXA - Transfer X to Accumulator
// Copies the current contents of the X register into the accumulator and sets
// the zero and negative flags as appropriate.
func (cpu *CPU) txa() {
	cpu.registerA = cpu.registerX
	cpu.setFlagZeroAndNegativeForResult(cpu.registerA)
}

// TXS - Transfer X to Stack Pointer
// Copies the current contents of the X register into the stack register.
func (cpu *CPU) txs() {
	cpu.stackPush(cpu.registerX)
}

// TYA - Transfer Y to Accumulator
// Copies the current contents of the Y register into the accumulator and sets
// the zero and negative flags as appropriate.
func (cpu *CPU) tya() {
	cpu.registerA = cpu.registerY
	cpu.setFlagZeroAndNegativeForResult(cpu.registerA)
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

func (cpu *CPU) branch(shouldBranch bool) {
	if shouldBranch {
		jump := cpu.memRead(cpu.programCounter)
		jump_addr := cpu.programCounter + uint16(1) + uint16(jump)
		cpu.programCounter = jump_addr
	}
}

func (cpu *CPU) addToRegisterA(value uint8) {
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
		case "ADC":
			cpu.adc(opcode.AddressingMode)
		case "AND":
			cpu.and(opcode.AddressingMode)
		case "ASL":
			cpu.asl(opcode.AddressingMode)
		case "BCC":
			cpu.bcc()
		case "BCS":
			cpu.bcs()
		case "BEQ":
			cpu.beq()
		case "BIT":
			cpu.bit(opcode.AddressingMode)
		case "BMI":
			cpu.bmi()
		case "BNE":
			cpu.bne()
		case "BPL":
			cpu.bpl()
		case "BRK":
			cpu.brk()
			return
		case "BVC":
			cpu.bvc()
		case "BVS":
			cpu.bvs()
		case "CLC":
			cpu.clc()
		case "CLD":
			cpu.cld()
		case "CLI":
			cpu.cli()
		case "CLV":
			cpu.clv()
		case "CMP":
			cpu.cmp(opcode.AddressingMode)
		case "CPX":
			cpu.cpx(opcode.AddressingMode)
		case "CPY":
			cpu.cpy(opcode.AddressingMode)
		case "DEC":
			cpu.dec(opcode.AddressingMode)
		case "DEX":
			cpu.dex(opcode.AddressingMode)
		case "DEY":
			cpu.dey(opcode.AddressingMode)
		case "EOR":
			cpu.eor(opcode.AddressingMode)
		case "INC":
			cpu.inc(opcode.AddressingMode)
		case "INX":
			cpu.inx()
		case "INY":
			cpu.iny()
		case "JMP":
			cpu.jmp(opcode.AddressingMode)
		case "JSR":
			cpu.jsr()
		case "LDA":
			cpu.lda(opcode.AddressingMode)
		case "LDX":
			cpu.ldx(opcode.AddressingMode)
		case "LDY":
			cpu.ldy(opcode.AddressingMode)
		case "NOP":
			cpu.nop()
		case "ORA":
			cpu.ora(opcode.AddressingMode)
		case "PHA":
			cpu.pha()
		case "PHP":
			cpu.php()
		case "PLA":
			cpu.pla()
		case "PLP":
			cpu.plp()
		case "ROL":
			cpu.rol(opcode.AddressingMode)
		case "ROR":
			cpu.ror(opcode.AddressingMode)
		case "RTI":
			cpu.rti()
		case "RTS":
			cpu.rts()
		case "SBC":
			cpu.sbc(opcode.AddressingMode)
		case "SEC":
			cpu.sec()
		case "SED":
			cpu.sed()
		case "SEI":
			cpu.sei()
		case "STA":
			cpu.sta(opcode.AddressingMode)
		case "STX":
			cpu.stx(opcode.AddressingMode)
		case "STY":
			cpu.sty(opcode.AddressingMode)
		case "TAX":
			cpu.tax()
		case "TAY":
			cpu.tay()
		case "TSX":
			cpu.tsx()
		case "TXA":
			cpu.txa()
		case "TXS":
			cpu.txs()
		case "TYA":
			cpu.tya()
		default:
			panic(fmt.Sprintf("Unsupported opcode: 0x%x\n", opcode))
		}

		if programCounterState == cpu.programCounter {
			cpu.programCounter += uint16(opcode.Bytes) - 1
		}
	}
}
