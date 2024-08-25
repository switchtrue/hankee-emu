package main

type Flag = uint8

// 7 6 5 4 3 2 1 0
// N V _ B D I Z C
// | |   | | | | +--- Carry Flag
// | |   | | | +----- Zero Flag
// | |   | | +------- Interrupt Disable
// | |   | +--------- Decimal Mode (not used on NES)
// | |   +----------- Break Command
// | +--------------- Overflow Flag
// +----------------- Negative Flag
const (
	FlagCarry           Flag = 1 << 0
	FlagZero                 = 1 << 1
	FlagInterruptDiable      = 1 << 2
	FlagDecimalMode          = 1 << 3
	FlagBreakCommand         = 1 << 4
	FlagOverflow             = 1 << 6
	FlagNegative             = 1 << 7
)

// Generic method to set any flag to 1 or 0
func (cpu *CPU) setFlag(flag Flag, isSet bool) {
	if isSet {
		cpu.status |= flag
		return
	}

	cpu.status &= 0xFF - flag
}

// Generic method to get any flag balue as a boolean
func (cpu *CPU) getFlag(flag Flag) bool {
	return cpu.status&flag != 0
}

func (cpu *CPU) setFlagZero(isSet bool) {
	cpu.setFlag(FlagZero, isSet)
}

// Set the Zero Flag based on the result of an operation
func (cpu *CPU) setFlagZeroForResult(result uint8) {
	cpu.setFlag(FlagZero, result == 0)
}

// Gets the Zero flag as a boolean
func (cpu *CPU) getFlagZero() bool {
	return cpu.getFlag(FlagZero)
}

// Sets the Interrupt Disable flag to 1 or 0
func (cpu *CPU) setFlagInterruptDisable(isSet bool) {
	cpu.setFlag(FlagInterruptDiable, isSet)
}

// Sets the Decimal Mode flag to 1 or 0
func (cpu *CPU) setFlagDecimalMode(isSet bool) {
	cpu.setFlag(FlagDecimalMode, isSet)
}

// Sets the carry flag to 1 or 0
func (cpu *CPU) setFlagCarry(isSet bool) {
	cpu.setFlag(FlagCarry, isSet)
}

// Gets the carry flag as a boolean
func (cpu *CPU) getFlagCarry() bool {
	return cpu.getFlag(FlagCarry)
}

// Sets the Carry flag based on the result of an operation
func (cpu *CPU) setFlagCarryForResult(result uint8) {
	cpu.setFlag(FlagCarry, result > 0xff)
}

// Sets the Break Command flag to 1 or 0
func (cpu *CPU) setFlagBreakCommand(isSet bool) {
	cpu.setFlag(FlagBreakCommand, isSet)
}

// Sets the Overflow flag to 1 or 0
func (cpu *CPU) setFlagOverflow(isSet bool) {
	cpu.setFlag(FlagOverflow, isSet)
}

// Gets the Overflow flag as a boolean
func (cpu *CPU) getFlagOverflow() bool {
	return cpu.getFlag(FlagOverflow)
}

func (cpu *CPU) setFlagNegative(isSet bool) {
	cpu.setFlag(FlagNegative, isSet)
}

// Sets the Negative flag base on the result of an operation
func (cpu *CPU) setFlagNegativeForResult(result uint8) {
	cpu.setFlag(FlagNegative, result&0b1000_0000 != 0)
}

// Gets the Negative flag as a boolean
func (cpu *CPU) getFlagNegative() bool {
	return cpu.getFlag(FlagNegative)
}

// Sets both the Zero and Negative flags based on the result of an operation.
// This is a convenience functions that sets both in one go as they are
// frequently set together
func (cpu *CPU) setFlagZeroAndNegativeForResult(result uint8) {
	cpu.setFlagZeroForResult(result)
	cpu.setFlagNegativeForResult(result)
}
