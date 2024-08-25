package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// More of a validation that a test, this ensures that all keys in the map match
// the OpCode in the ObCode struct map value
func TestValidateMapKeysAndOpCode(t *testing.T) {
	for key, value := range CPU_OP_CODE_TABLE {
		assert.Equal(t, key, value.Opcode, "OpCode key %x does not match %x for %s", key, value.Opcode, value.Name)
	}
}

// Test executing all of the opcodes just to ensure they've all been implemented
// and don't panic. This doesn't check any of them actually work.'
func TestAllOpCodes(t *testing.T) {
	cpu := NewCPU()

	var program []uint8
	for key, _ := range CPU_OP_CODE_TABLE {
		// Skip the BRK command as this will end the program.
		if key != 0x00 {
			program = append(program, key)
		}
	}
	program = append(program, 0x00)

	cpu.loadAndRun(program)
}
