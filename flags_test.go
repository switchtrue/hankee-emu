package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// Test that the zero flag is set when the value in register a is zero
func Test_Flags_ZeroFlag_Set(t *testing.T) {
	cpu := NewCPU()
	cpu.setFlagZeroAndNegativeForResult(0)
	assertZeroFlagSet(t, cpu.status)
}

// Test that the zero flag is not set when the value in register a is not zero
func Test_Flags_ZeroFlag_NotSet(t *testing.T) {
	cpu := NewCPU()
	cpu.setFlagZeroAndNegativeForResult(1)
	assertZeroFlagNotSet(t, cpu.status)
}

// Test that the negative flag is set when the value in register a is negative
func Test_Flags_NegativeFlag_Set(t *testing.T) {
	cpu := NewCPU()
	cpu.setFlagZeroAndNegativeForResult(0xff)
	assertNegativeFlagSet(t, cpu.status)
}

// Test that the negative flag is not set when the value in register a is positive
func Test_Flags_NegativeFlag_NotSet(t *testing.T) {
	cpu := NewCPU()
	cpu.setFlagZeroAndNegativeForResult(1)
	assertNegativeFlagNotSet(t, cpu.status)
}
