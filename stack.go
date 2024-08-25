package main

const (
	STACK       uint16 = 0x0100
	STACK_RESET uint8  = 0xfd
)

func (cpu *CPU) stackPush(value uint8) {
	cpu.memWrite(STACK+uint16(cpu.stackPointer), value)
	cpu.stackPointer = cpu.stackPointer - 1
}

func (cpu *CPU) stackPop() uint8 {
	cpu.stackPointer += 1
	return cpu.memRead(STACK + uint16(cpu.stackPointer))
}

func (cpu *CPU) stackPushUInt16(value uint16) {
	cpu.memWriteUInt16(STACK+uint16(cpu.stackPointer), value)
	cpu.stackPointer = cpu.stackPointer - 1
}

func (cpu *CPU) stackPopUInt16() uint16 {
	cpu.stackPointer += 1
	return cpu.memReadUInt16(STACK + uint16(cpu.stackPointer))
}
