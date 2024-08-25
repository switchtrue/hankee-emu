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

func (cpu *CPU) setFlag(flag Flag, isSet bool) {
	if isSet {
		cpu.status |= flag
		return
	}

	cpu.status &= 0xFF - flag
}

func (cpu *CPU) setFlagZeroForResult(result uint8) {
	cpu.setFlag(FlagZero, result == 0)
}

func (cpu *CPU) setFlagInterruptDisable(isSet bool) {
	cpu.setFlag(FlagInterruptDiable, isSet)
}

func (cpu *CPU) setFlagDecimalMode(isSet bool) {
	cpu.setFlag(FlagDecimalMode, isSet)
}

func (cpu *CPU) setFlagCarry(isSet bool) {
	cpu.setFlag(FlagCarry, isSet)
}

func (cpu *CPU) getFlagCarry() bool {
	return cpu.status&FlagCarry != 0
}

func (cpu *CPU) setFlagCarryForResult(result uint8) {
	cpu.setFlag(FlagCarry, result > 0xff)
}

func (cpu *CPU) setFlagBreakCommand(isSet bool) {
	cpu.setFlag(FlagBreakCommand, isSet)
}

func (cpu *CPU) setFlagOverflow(isSet bool) {
	cpu.setFlag(FlagOverflow, isSet)
}

func (cpu *CPU) setFlagNegativeForResult(result uint8) {
	cpu.setFlag(FlagNegative, result&0b1000_0000 != 0)
}

func (cpu *CPU) setFlagZeroAndNegativeForResult(result uint8) {
	cpu.setFlagZeroForResult(result)
	cpu.setFlagNegativeForResult(result)
}
