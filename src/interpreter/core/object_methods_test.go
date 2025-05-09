package core_test

import (
	"testing"
	"unsafe"

	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
	"smalltalklsp/interpreter/vm"
)

func TestObjectIsTrue(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	tests := []struct {
		name     string
		obj      core.ObjectInterface
		expected bool
	}{
		{
			name:     "Boolean true",
			obj:      core.NewBoolean(true),
			expected: true,
		},
		{
			name:     "Boolean false",
			obj:      core.NewBoolean(false),
			expected: false,
		},
		{
			name:     "Nil",
			obj:      core.NewNil(),
			expected: false,
		},
		{
			name:     "Integer",
			obj:      virtualMachine.NewInteger(42),
			expected: false,
		}, {
			name:     "Float",
			obj:      virtualMachine.NewFloat(42.0),
			expected: false,
		},
		{
			name:     "String",
			obj:      classes.StringToObject(classes.NewString("hello")),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.obj.IsTrue()
			if result != test.expected {
				t.Errorf("Expected %v, got %v for %s", test.expected, result, test.name)
			}
		})
	}
}

func TestObjectInstanceVarMethods(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	// Create a class with instance variables
	class := classes.NewClass("TestClass", nil)
	class.InstanceVarNames = append(class.InstanceVarNames, "var1", "var2")

	// Create an instance
	instance := core.NewInstance((*core.Class)(unsafe.Pointer(class)))

	// Test GetInstanceVarByIndex
	instance.SetInstanceVarByIndex(0, virtualMachine.NewInteger(42))
	instance.SetInstanceVarByIndex(1, classes.StringToObject(classes.NewString("hello")))

	// Get the instance variable and check its value
	var0 := instance.GetInstanceVarByIndex(0)
	if core.IsIntegerImmediate(var0) {
		intValue := core.GetIntegerImmediate(var0)
		if intValue != 42 {
			t.Errorf("Expected instance var 0 to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", var0)
	}

	var1 := instance.GetInstanceVarByIndex(1)
	var1Str := classes.ObjectToString(var1)
	if var1.Type() != core.OBJ_STRING || var1Str.Value != "hello" {
		t.Errorf("Expected instance variable 1 to be a string with value 'hello', got %v", var1)
	}

	// Test GetInstanceVarByIndex with out of bounds index
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on out of bounds access, but no panic occurred")
		}
	}()
	instance.GetInstanceVarByIndex(2) // This should panic
}

func TestObjectSetInstanceVarByIndex(t *testing.T) {
	// Create a VM for testing
	virtualMachine := vm.NewVM()

	// Create a class with instance variables
	class := classes.NewClass("TestClass", nil)
	class.InstanceVarNames = append(class.InstanceVarNames, "var1", "var2")

	// Create an instance
	instance := core.NewInstance((*core.Class)(unsafe.Pointer(class)))

	// Test SetInstanceVarByIndex
	instance.SetInstanceVarByIndex(0, virtualMachine.NewInteger(42))
	instance.SetInstanceVarByIndex(1, classes.StringToObject(classes.NewString("hello")))

	// Check the instance variables
	var instanceVars2 = instance.InstanceVars()
	var0 := instanceVars2[0]
	if core.IsIntegerImmediate(var0) {
		intValue := core.GetIntegerImmediate(var0)
		if intValue != 42 {
			t.Errorf("Expected instance var 0 to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", var0)
	}

	var1 := instanceVars2[1]
	var1Str := classes.ObjectToString(var1)
	if var1.Type() != core.OBJ_STRING || var1Str.Value != "hello" {
		t.Errorf("Expected instance variable 1 to be a string with value 'hello', got %v", var1)
	}

	// Test SetInstanceVarByIndex with out of bounds index
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on out of bounds access, but no panic occurred")
		}
	}()
	instance.SetInstanceVarByIndex(2, virtualMachine.NewInteger(42)) // This should panic
}

func TestObjectGetMethodDict(t *testing.T) {
	// Create a VM for testing
	_ = vm.NewVM()

	// Test with a class
	class := classes.NewClass("TestClass", nil)
	methodDict := class.GetMethodDict()
	if methodDict.Type() != core.OBJ_DICTIONARY {
		t.Errorf("Expected method dictionary to be a dictionary, got %v", methodDict.Type())
	}

	// Test with a non-class object
	instance := core.NewInstance((*core.Class)(unsafe.Pointer(class)))
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling GetMethodDict on a non-class object, but no panic occurred")
		}
	}()
	instance.GetMethodDict() // This should panic
}
