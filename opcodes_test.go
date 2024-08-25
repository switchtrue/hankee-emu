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
