package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that the zero flag is set when the value in register a is zero
func Test_Flags_ZeroFlag_Set(t *testing.T) {
	cpu := NewCPU()
	cpu.updateZeroAndNegativeFlags(0)
	assertZeroFlagSet(t, cpu.status)
}

// Test that the zero flag is not set when the value in register a is not zero
func Test_Flags_ZeroFlag_NotSet(t *testing.T) {
	cpu := NewCPU()
	cpu.updateZeroAndNegativeFlags(1)
	assertZeroFlagNotSet(t, cpu.status)
}

// Test that the negative flag is set when the value in register a is negative
func Test_Flags_NegativeFlag_Set(t *testing.T) {
	cpu := NewCPU()
	cpu.updateZeroAndNegativeFlags(0xff)
	assertNegativeFlagSet(t, cpu.status)
}

// Test that the negative flag is not set when the value in register a is positive
func Test_Flags_NegativeFlag_NotSet(t *testing.T) {
	cpu := NewCPU()
	cpu.updateZeroAndNegativeFlags(1)
	assertNegativeFlagNotSet(t, cpu.status)
}

func assertNegativeFlagSet(t *testing.T, status uint8) {
	assert.Equal(t, uint8(0b1000_0000), status&0b1000_0000, "")
}

func assertNegativeFlagNotSet(t *testing.T, status uint8) {
	assert.Equal(t, uint8(0b00), status&0b1000_0000, "")
}

func assertZeroFlagSet(t *testing.T, status uint8) {
	assert.Equal(t, uint8(0b10), status&0b0000_0010, "")
}

func assertZeroFlagNotSet(t *testing.T, status uint8) {
	assert.Equal(t, uint8(0b00), status&0b0000_0010, "")
}

// Test the 0xa9 immediate load opcode by loading 0x05 into register a and checking it's there.
// Also checks that the zero flag and negative flags are both not set
func Test_0xa9_LDA_ImmediateLoadData(t *testing.T) {
	cpu := NewCPU()
	cpu.loadAndRun([]uint8{0xa9, 0x05, 0x00})
	assert.Equal(t, uint8(0x05), cpu.register_a, "")
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Test the 0xa9 immediate load opcode by loading 0x05 into register a and checking it's there.
// Checks that the negative flag is set
func Test_0xa9_LDA_ImmediateLoadData_NegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0x00})
	assert.Equal(t, uint8(0xff), cpu.register_a, "")
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagSet(t, cpu.status)
}

// Test the 0xa9 immediate load opcode by loading 0x05 into register a and checking it's there.
// Checks that the zero flag is set
func Test_0xa9_LDA_ImmediateLoadData_ZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.loadAndRun([]uint8{0xa9, 0x00, 0x00})
	assert.Equal(t, uint8(0x00), cpu.register_a, "")
	assertZeroFlagSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Test that we can successfully copy register A to register X
func Test_0xaa_TAX_MoveAToX(t *testing.T) {
	cpu := NewCPU()
	// LDA 10, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0x0a, 0xaa, 0x00})
	assert.Equal(t, uint8(10), cpu.register_x)
}

// Test that we can successfully copy register A to register X
// Checks that the negative flag is set
func Test_0xaa_TAX_MoveAToX_NegativeFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA -1, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0xaa, 0x00})
	assert.Equal(t, uint8(0xff), cpu.register_x)
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagSet(t, cpu.status)
}

// Test that we can successfully copy register A to register X
// Checks that the zero flag is set
func Test_0xaa_TAX_MoveAToX_ZeroFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA 0, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0x00, 0xaa, 0x00})
	assert.Equal(t, uint8(0), cpu.register_x)
	assertZeroFlagSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Thet that increment X wikll increment the value of X by 1
func Test_0xe8_INX_IncrementX(t *testing.T) {
	cpu := NewCPU()
	// LDA 0, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0x00, 0xe8, 0x00})
	assert.Equal(t, uint8(1), cpu.register_x, "")
}

// Thet that increment X will overflow and wrap the x register
func Test_0xe8_INX_IncrementX_Overflow(t *testing.T) {
	cpu := NewCPU()
	// LDA -1, TAX INX, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0xaa, 0xe8, 0xe8, 0x00})
	assert.Equal(t, uint8(1), cpu.register_x, "")
}

// Test that INX will correctly set the negative flag based on its result
func Test_0xe8_INX_IncrementX_NegativeFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA -2, TAX, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xfe, 0xaa, 0xe8, 0x00})
	assert.Equal(t, uint8(0xff), cpu.register_x)
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagSet(t, cpu.status)
}

// Test that INX will correctly set the zero flag based on its result
func Test_0xe8_INX_IncrementX_ZeroFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA -1, TAX, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0xaa, 0xe8, 0x00})
	assert.Equal(t, uint8(0x00), cpu.register_x)
	assertZeroFlagSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Test six op codes working together as a mini program
func Test_SixOpsWorkingTogether(t *testing.T) {
	cpu := NewCPU()
	// LDA -64, TAX, NOP, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xea, 0xe8, 0x00})
	assert.Equal(t, uint8(0xc0), cpu.register_a, "")
	assert.Equal(t, uint8(0xc1), cpu.register_x, "")
}
