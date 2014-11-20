package go8086

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVMFlag(t *testing.T) {
	vm := NewVM()
	assert.Equal(t, 0, vm.GetFlag(CF))
	assert.Equal(t, 0, vm.GetFlag(OF))
	vm.flag = 0x0001
	assert.Equal(t, 1, vm.GetFlag(CF))
	assert.Equal(t, 0, vm.GetFlag(OF))
	vm.flag = 0x0400
	assert.Equal(t, 0, vm.GetFlag(CF))
	assert.Equal(t, 1, vm.GetFlag(OF))
}

func TestVMSetFlag(t *testing.T) {
	vm := NewVM()
	assert.Equal(t, 0, vm.GetFlag(CF))
	vm.SetFlag(CF, true)
	assert.Equal(t, 1, vm.GetFlag(CF))
	assert.Equal(t, 0x0001, vm.flag)
	vm.SetFlag(CF, true)
	assert.Equal(t, 1, vm.GetFlag(CF))
	assert.Equal(t, 0x0001, vm.flag)
}

func TestVMClearFlag(t *testing.T) {
	vm := NewVM()
	vm.flag = 0x0001
	assert.Equal(t, 1, vm.GetFlag(CF))
	vm.SetFlag(CF, false)
	assert.Equal(t, 0, vm.GetFlag(CF))
	assert.Equal(t, 0x0000, vm.flag)
	vm.SetFlag(CF, false)
	assert.Equal(t, 0, vm.GetFlag(CF))
	assert.Equal(t, 0x0000, vm.flag)
}
