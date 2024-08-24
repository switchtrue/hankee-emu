package main

import "fmt"

type CPU struct {
	register_a      uint8
	register_x      uint8
	status          uint8
	program_counter uint16
	memory          []uint8
}

func NewCPU() *CPU {
	return &CPU{
		register_a:      0,
		register_x:      0,
		status:          0,
		program_counter: 0,
		memory:          make([]uint8, 0xFFFF),
	}
}

func (cpu *CPU) lda(value uint8) {
	cpu.register_a = value
	cpu.updateZeroAndNegativeFlags(cpu.register_a)
}

func (cpu *CPU) tax() {
	cpu.register_x = cpu.register_a
	cpu.updateZeroAndNegativeFlags(cpu.register_x)
}

func (cpu *CPU) inx() {
	cpu.register_x += uint8(1)
	cpu.updateZeroAndNegativeFlags(cpu.register_x)
}

func (cpu *CPU) nop() {
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

func (cpu *CPU) memRead(addr uint16) uint8 {
	return cpu.memory[addr]
}

func (cpu *CPU) memWrite(addr uint16, data uint8) {
	cpu.memory[addr] = data
}

func (cpu *CPU) memReadUInt16(addr uint16) uint16 {
	lo := uint16(cpu.memRead(addr))
	hi := uint16(cpu.memRead(addr + 1))
	return (hi << 8) | (lo)
}

func (cpu *CPU) memWriteUInt16(addr uint16, data uint16) {
	hi := uint8(data >> 8)
	lo := uint8(data & 0xff)
	cpu.memWrite(addr, lo)
	cpu.memWrite(addr+1, hi)
}

func (cpu *CPU) reset() {
	cpu.register_a = 0
	cpu.register_x = 0
	cpu.status = 0
	cpu.program_counter = cpu.memReadUInt16(0xFFFC)
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
		opcode := cpu.memRead(uint16(cpu.program_counter))
		cpu.program_counter++
		switch opcode {
		case 0xA9:
			param := cpu.memRead(cpu.program_counter)
			cpu.program_counter++
			cpu.lda(param)
		case 0xAA:
			cpu.tax()
		case 0xE8:
			cpu.inx()
		case 0xEA:
			cpu.nop()
		case 0x00:
			return
		default:
			fmt.Printf("Unsupported opcode: 0x%x\n", opcode)
		}
	}
}
