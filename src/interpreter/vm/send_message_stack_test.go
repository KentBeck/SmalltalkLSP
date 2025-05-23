package vm_test

import (
	"smalltalklsp/interpreter/pile"
	"testing"

	"smalltalklsp/interpreter/compiler"
	"smalltalklsp/interpreter/vm"
)

func TestSendMessageStackManagement(t *testing.T) {
	// Pile package integration has been fixed
	virtualMachine := vm.NewVM()

	// We'll use the VM's Object and Integer classes
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])

	// Create literals
	returnValueSelector := pile.NewSymbol("returnValue")
	valueObj := virtualMachine.NewInteger(42)
	receiverObj := virtualMachine.NewInteger(10)

	// Create a simple method that returns a value using AddLiteral
	returnValueBuilder := compiler.NewMethodBuilder(integerClass)

	// Add literals to the method builder
	valueIndex, returnValueBuilder := returnValueBuilder.AddLiteral(valueObj) // Literal 0: 42

	// Create bytecodes for the method: just return 42
	returnValueBuilder.PushLiteral(valueIndex)
	returnValueBuilder.ReturnStackTop()

	// Finalize the method
	returnValueBuilder.Go("returnValue")

	// Create a caller method that will call returnValue and then use the result using AddLiteral
	callerBuilder := compiler.NewMethodBuilder(integerClass)

	// Add literals to the caller method builder
	receiverIndex, callerBuilder := callerBuilder.AddLiteral(receiverObj)                    // Literal 0: 10
	returnValueSelectorIndex, callerBuilder := callerBuilder.AddLiteral(returnValueSelector) // Literal 1: returnValue

	// Create bytecodes for the caller method:
	// 1. Push a value onto the stack that should be preserved
	// 2. Send returnValue message to receiver
	// 3. Check that both the original value and the result are on the stack

	// PUSH_LITERAL receiverIndex (10) - this is a value we want to preserve across the method call
	callerBuilder.PushLiteral(receiverIndex)

	// PUSH_SELF - this will be the receiver of the returnValue message
	callerBuilder.PushSelf()

	// SEND_MESSAGE returnValue with 0 arguments
	callerBuilder.SendMessage(returnValueSelectorIndex, 0)

	// At this point, the stack should have two values:
	// 1. The original value (10)
	// 2. The result of the returnValue method (42)

	// RETURN_STACK_TOP - just return the top of the stack (which should be 42)
	callerBuilder.ReturnStackTop()

	// Finalize the method
	callerMethod := callerBuilder.Go("caller")

	// Create a receiver for the caller method
	receiver := virtualMachine.NewInteger(5)

	// Create a context for the caller method
	context := vm.NewContext(callerMethod, receiver, []*pile.Object{}, nil)

	// Execute the context
	result, err := virtualMachine.ExecuteContext(context)
	if err != nil {
		t.Errorf("Error executing caller method: %v", err)
		return
	}

	// Check that the result is 42 (the value returned by the returnValue method)
	if pile.IsIntegerImmediate(result) {
		intValue := pile.GetIntegerImmediate(result)
		if intValue != 42 {
			t.Errorf("Expected result to be 42, got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", result)
	}
}

func TestSendMessageWithMultiplication(t *testing.T) {
	// Pile package integration has been fixed
	virtualMachine := vm.NewVM()

	// We'll use the VM's Object and Integer classes
	integerClass := pile.ObjectToClass(virtualMachine.Globals["Integer"])

	// Create literals
	returnValueSelector := pile.NewSymbol("returnValue")
	valueObj := virtualMachine.NewInteger(42)
	timesSelector := pile.NewSymbol("*")

	// Create a simple method that returns a value using AddLiteral
	returnValueBuilder := compiler.NewMethodBuilder(integerClass)

	// Add literals to the method builder
	valueIndex, returnValueBuilder := returnValueBuilder.AddLiteral(valueObj) // Literal 0: 42

	// Create bytecodes for the method: just return 42
	returnValueBuilder.PushLiteral(valueIndex)
	returnValueBuilder.ReturnStackTop()

	// Finalize the method
	returnValueBuilder.Go("returnValue")

	// Get the predefined multiplication method from the VM
	timesMethod := virtualMachine.LookupMethod(pile.ClassToObject(integerClass), timesSelector)

	// Make sure the method is in the method dictionary
	methodDict := integerClass.GetMethodDict()
	dict := pile.ObjectToDictionary(methodDict)
	dict.Entries["*"] = timesMethod

	// Create a method that will call returnValue and then use the result for multiplication using AddLiteral
	multiplyBuilder := compiler.NewMethodBuilder(integerClass)

	// Add literals to the multiply method builder
	returnValueSelectorIndex, multiplyBuilder := multiplyBuilder.AddLiteral(returnValueSelector) // Literal 0: returnValue
	timesSelectorIndex, multiplyBuilder := multiplyBuilder.AddLiteral(timesSelector)             // Literal 1: *

	// Create bytecodes for the multiply method:
	// 1. Push self (for later use in multiplication)
	// 2. Send returnValue message to self
	// 3. Multiply self by the result

	// PUSH_SELF (for later use in multiplication)
	multiplyBuilder.PushSelf()

	// DUPLICATE (to save a copy for later multiplication)
	multiplyBuilder.PushSelf()

	// SEND_MESSAGE returnValue with 0 arguments
	multiplyBuilder.SendMessage(returnValueSelectorIndex, 0)

	// SEND_MESSAGE * with 1 argument
	multiplyBuilder.SendMessage(timesSelectorIndex, 1)

	// RETURN_STACK_TOP
	multiplyBuilder.ReturnStackTop()

	// Finalize the method
	multiplyMethod := multiplyBuilder.Go("multiply")

	// Create a receiver for the multiply method
	multiplyReceiver := virtualMachine.NewInteger(5)

	// Create a context for the multiply method
	multiplyContext := vm.NewContext(multiplyMethod, multiplyReceiver, []*pile.Object{}, nil)

	// Execute the context
	multiplyResult, err := virtualMachine.ExecuteContext(multiplyContext)
	if err != nil {
		t.Errorf("Error executing multiply method: %v", err)
		return
	}

	// Check that the result is 5 * 42 = 210
	if pile.IsIntegerImmediate(multiplyResult) {
		intValue := pile.GetIntegerImmediate(multiplyResult)
		if intValue != 210 {
			t.Errorf("Expected result to be 210 (5 * 42), got %d", intValue)
		}
	} else {
		t.Errorf("Expected an immediate integer, got %v", multiplyResult)
	}
}
