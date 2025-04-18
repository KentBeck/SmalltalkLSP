package main

import (
	"testing"
	"unsafe"
)

func TestSizeOfObject(t *testing.T) {
	obj := &Object{}
	size := unsafe.Sizeof(*obj)

	// Assert that the size is less than or equal to 232 bytes
	const maxSize = 232
	if size > maxSize {
		t.Errorf("Object size is %d bytes, which exceeds the maximum allowed size of %d bytes", size, maxSize)
	}
}
