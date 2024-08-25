package main

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
