package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/vm"
)

// TestArrayAtPrimitive tests the Array at: primitive
func TestArrayAtPrimitive(t *testing.T) {
	virtualMachine := vm.NewVM()

	// Get the predefined primitive methods from the VM
	arrayClass := pile.ObjectToClass(virtualMachine.Globals["Array"])
	atSelector := pile.NewSymbol("at:")
	atMethod := virtualMachine.LookupMethod(pile.ClassToObject(arrayClass), atSelector)

	// Create a test array with 3 elements
	array := virtualMachine.NewArray(3)
	arrayObj := pile.ObjectToArray(array)
	
	// Fill the array with values
	arrayObj.AtPut(0, virtualMachine.NewInteger(1))
	arrayObj.AtPut(1, virtualMachine.NewInteger(2))
	arrayObj.AtPut(2, virtualMachine.NewInteger(3))
	
	// Create the index argument (1-based in Smalltalk)
	indexArg := virtualMachine.NewInteger(2)
	
	// Execute the primitive
	result := virtualMachine.ExecutePrimitive(array, atSelector, []*pile.Object{indexArg}, atMethod)

	// Check that the result is not nil
	if result == nil {
		t.Errorf("Array at: primitive returned nil")
		return
	}

	// Check that the result is an integer
	if !pile.IsIntegerImmediate(result) {
		t.Errorf("Expected integer result, got %v", result.Type())
		return
	}

	// Check that the result is 2 (the value at index 1, which is the second element)
	value := pile.GetIntegerImmediate(result)
	if value != 2 {
		t.Errorf("Expected value 2, got %d", value)
	}

	// Test with an out-of-bounds index
	outOfBoundsArg := virtualMachine.NewInteger(10)
	
	// This should panic, so we need to recover
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for out-of-bounds index, but no panic occurred")
		}
	}()
	
	// This should panic
	virtualMachine.ExecutePrimitive(array, atSelector, []*pile.Object{outOfBoundsArg}, atMethod)
}
