package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the 0xa9 immediate load opcode by loading 0x05 into register a and checking it's there.
// Also checks that the zero flag and negative flags are both not set
func Test_0xa9_LDA_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.loadAndRun([]uint8{0xa9, 0x05, 0x00})
	assert.Equal(t, uint8(0x05), cpu.registerA, "")
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Test the 0xa9 immediate load opcode by loading 0x05 into register a and checking it's there.
// Checks that the negative flag is set
func Test_0xa9_LDA_Immediate_NegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0x00})
	assert.Equal(t, uint8(0xff), cpu.registerA, "")
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagSet(t, cpu.status)
}

// Test the 0xa9 immediate load opcode by loading 0x05 into register a and checking it's there.
// Checks that the zero flag is set
func Test_0xa9_LDA_Immediate_ZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.loadAndRun([]uint8{0xa9, 0x00, 0x00})
	assert.Equal(t, uint8(0x00), cpu.registerA, "")
	assertZeroFlagSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Test that we can load immediate into A from memory
func Test_0xa5_LDA_ZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.memWrite(0x10, 0x55)
	cpu.loadAndRun([]uint8{0xa5, 0x10, 0x00})
	assert.Equal(t, uint8(0x55), cpu.registerA, "")
}

// Test that we can successfully copy register A to register X
func Test_0xaa_TAX_MoveAToX(t *testing.T) {
	cpu := NewCPU()
	// LDA 10, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0x0a, 0xaa, 0x00})
	assert.Equal(t, uint8(10), cpu.registerX)
}

// Test that we can successfully copy register A to register X
// Checks that the negative flag is set
func Test_0xaa_TAX_MoveAToX_NegativeFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA -1, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0xaa, 0x00})
	assert.Equal(t, uint8(0xff), cpu.registerX)
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagSet(t, cpu.status)
}

// Test that we can successfully copy register A to register X
// Checks that the zero flag is set
func Test_0xaa_TAX_MoveAToX_ZeroFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA 0, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0x00, 0xaa, 0x00})
	assert.Equal(t, uint8(0), cpu.registerX)
	assertZeroFlagSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Thet that increment X wikll increment the value of X by 1
func Test_0xe8_INX_IncrementX(t *testing.T) {
	cpu := NewCPU()
	// LDA 0, TAX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0x00, 0xe8, 0x00})
	assert.Equal(t, uint8(1), cpu.registerX, "")
}

// Thet that increment X will overflow and wrap the x register
func Test_0xe8_INX_IncrementX_Overflow(t *testing.T) {
	cpu := NewCPU()
	// LDA -1, TAX INX, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0xaa, 0xe8, 0xe8, 0x00})
	assert.Equal(t, uint8(1), cpu.registerX, "")
}

// Test that INX will correctly set the negative flag based on its result
func Test_0xe8_INX_IncrementX_NegativeFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA -2, TAX, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xfe, 0xaa, 0xe8, 0x00})
	assert.Equal(t, uint8(0xff), cpu.registerX)
	assertZeroFlagNotSet(t, cpu.status)
	assertNegativeFlagSet(t, cpu.status)
}

// Test that INX will correctly set the zero flag based on its result
func Test_0xe8_INX_IncrementX_ZeroFlag(t *testing.T) {
	cpu := NewCPU()
	// LDA -1, TAX, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xff, 0xaa, 0xe8, 0x00})
	assert.Equal(t, uint8(0x00), cpu.registerX)
	assertZeroFlagSet(t, cpu.status)
	assertNegativeFlagNotSet(t, cpu.status)
}

// Test six op codes working together as a mini program
func Test_SixOpsWorkingTogether(t *testing.T) {
	cpu := NewCPU()
	// LDA -64, TAX, NOP, INX, BRK
	cpu.loadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xea, 0xe8, 0x00})
	assert.Equal(t, uint8(0xc0), cpu.registerA, "")
	assert.Equal(t, uint8(0xc1), cpu.registerX, "")
}
