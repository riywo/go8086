package go8086

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVMFlag(t *testing.T) {
	vm := NewVM()
	assert.Equal(t, false, vm.GetFlag(CF))
	vm.flag = 0x0001
	assert.Equal(t, true, vm.GetFlag(CF))
	vm.flag = 0x0400
	assert.Equal(t, true, vm.GetFlag(OF))
}

func TestVMSetFlag(t *testing.T) {
	vm := NewVM()
	assert.Equal(t, false, vm.GetFlag(CF))
	vm.SetFlag(CF, true)
	assert.Equal(t, true, vm.GetFlag(CF))
	assert.Equal(t, 0x0001, vm.flag)
	vm.SetFlag(CF, true)
	assert.Equal(t, true, vm.GetFlag(CF))
	assert.Equal(t, 0x0001, vm.flag)
}

func TestVMClearFlag(t *testing.T) {
	vm := NewVM()
	vm.flag = 0x0001
	assert.Equal(t, true, vm.GetFlag(CF))
	vm.SetFlag(CF, false)
	assert.Equal(t, false, vm.GetFlag(CF))
	assert.Equal(t, 0x0000, vm.flag)
	vm.SetFlag(CF, false)
	assert.Equal(t, false, vm.GetFlag(CF))
	assert.Equal(t, 0x0000, vm.flag)
}
