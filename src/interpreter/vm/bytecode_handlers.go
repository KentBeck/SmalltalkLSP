package vm

import (
	"encoding/binary"
	"fmt"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/pile"
)

// ExecutePushLiteral executes the PUSH_LITERAL bytecode
func (vm *VM) ExecutePushLiteral(context *Context) error {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the literal index (4 bytes)
	index := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))
	if index < 0 || index >= len(method.Literals) {
		return fmt.Errorf("literal index out of bounds: %d", index)
	}

	// Push the literal onto the stack
	literal := method.Literals[index]
	context.Push(literal)
	return nil
}

// ExecutePushInstanceVariable executes the PUSH_INSTANCE_VARIABLE bytecode
func (vm *VM) ExecutePushInstanceVariable(context *Context) error {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the instance variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))
	class := vm.GetClass(context.Receiver.(*pile.Object))
	if index < 0 || index >= len(class.InstanceVarNames) {
		return fmt.Errorf("instance variable index out of bounds: %d", index)
	}

	// Push the instance variable onto the stack
	value := context.Receiver.GetInstanceVarByIndex(index)
	context.Push(value)
	return nil
}

// ExecutePushTemporaryVariable executes the PUSH_TEMPORARY_VARIABLE bytecode
func (vm *VM) ExecutePushTemporaryVariable(context *Context) error {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the temporary variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))

	// First try to get the variable from the current context
	if index < len(context.TempVars) {
		context.Push(context.GetTempVarByIndex(index))
	} else {
		// If the index is out of bounds in the current context,
		// try to get it from the outer context (for blocks)
		outerContext := context.Sender
		if outerContext != nil && index < len(outerContext.TempVars) {
			context.Push(outerContext.GetTempVarByIndex(index))
		} else {
			return fmt.Errorf("temporary variable index out of bounds: %d", index)
		}
	}
	return nil
}

// ExecutePushSelf executes the PUSH_SELF bytecode
func (vm *VM) ExecutePushSelf(context *Context) error {
	context.Push(context.Receiver)
	return nil
}

// ExecuteStoreInstanceVariable executes the STORE_INSTANCE_VARIABLE bytecode
func (vm *VM) ExecuteStoreInstanceVariable(context *Context) error {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the instance variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))
	class := vm.GetClass(context.Receiver.(*pile.Object))

	if index < 0 || index >= len(class.InstanceVarNames) {
		return fmt.Errorf("instance variable index out of bounds: %d", index)
	}

	// Pop the value from the stack
	value := context.Pop()

	// Store the value in the instance variable
	context.Receiver.SetInstanceVarByIndex(index, value)

	// Push the value back onto the stack
	context.Push(value)
	return nil
}

// ExecuteStoreTemporaryVariable executes the STORE_TEMPORARY_VARIABLE bytecode
func (vm *VM) ExecuteStoreTemporaryVariable(context *Context) error {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the temporary variable index (4 bytes)
	index := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))

	// Pop the value from the stack
	value := context.Pop()

	// First try to store the value in the current context's temporary variable
	if index < len(context.TempVars) {
		context.SetTempVarByIndex(index, value)
	} else {
		// If the index is out of bounds in the current context,
		// try to store it in the outer context (for blocks)
		outerContext := context.Sender
		if outerContext != nil && index < len(outerContext.TempVars) {
			outerContext.SetTempVarByIndex(index, value)
		} else {
			return fmt.Errorf("temporary variable index out of bounds: %d", index)
		}
	}

	// Push the value back onto the stack
	context.Push(value)
	return nil
}

// ExecuteSendMessage executes the SEND_MESSAGE bytecode
func (vm *VM) ExecuteSendMessage(context *Context) (*pile.Object, error) {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the selector index (4 bytes)
	selectorIndex := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))
	if selectorIndex < 0 || selectorIndex >= len(method.Literals) {
		return nil, fmt.Errorf("selector index out of bounds: %d", selectorIndex)
	}

	// Get the argument count (4 bytes)
	argCount := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+5:]))

	// Get the selector
	selector := method.Literals[selectorIndex]
	if selector.Type() != pile.OBJ_SYMBOL {
		return nil, fmt.Errorf("selector is not a symbol: %s", selector)
	}

	// Pop the arguments from the stack
	args := make([]*pile.Object, argCount)
	for i := argCount - 1; i >= 0; i-- {
		args[i] = context.Pop()
	}

	// Pop the receiver
	receiver := context.Pop()

	// Check for nil receiver
	if receiver == nil {
		return nil, fmt.Errorf("nil receiver for message: %s", pile.ObjectToSymbol(selector).GetValue())
	}

	methodObj := vm.LookupMethod(receiver, selector)
	if methodObj == nil {
		return nil, fmt.Errorf("method not found: %s", pile.ObjectToSymbol(selector).GetValue())
	}

	// Handle primitive methods
	if result := vm.ExecutePrimitive(receiver, selector, args, methodObj); result != nil {
		context.Push(result)
		return result, nil
	}

	// Create a new context for the method
	newContext := NewContext(methodObj, receiver, args, context)

	// Set the current context to the new context in the executor
	vm.Executor.CurrentContext = newContext

	// Return from this context execution to start executing the new context
	// We need to execute the new context immediately
	result, err := vm.ExecuteContext(newContext)
	if err != nil {
		return nil, err
	}

	// Check for nil result
	if result == nil {
		return nil, fmt.Errorf("method not found: %s", pile.ObjectToSymbol(selector).GetValue())
	}

	// Move back to the sender context in the executor
	vm.Executor.CurrentContext = context

	// Push the result onto the stack
	context.Push(result)

	// Return the result
	return result.(*pile.Object), nil
}

// ExecuteReturnStackTop executes the RETURN_STACK_TOP bytecode
func (vm *VM) ExecuteReturnStackTop(context *Context) (*pile.Object, error) {
	if context.StackPointer <= 0 {
		// If the stack is empty, return nil
		return vm.NilObject.(*pile.Object), nil
	}
	returnValue := context.Pop()
	return returnValue, nil
}

// ExecuteJump executes the JUMP bytecode
func (vm *VM) ExecuteJump(context *Context) (bool, error) {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the jump offset (4 bytes)
	if context.PC+1 >= len(method.Bytecodes) {
		return false, fmt.Errorf("jump offset out of bounds")
	}
	offset := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))

	// The offset is relative to the current instruction
	// We need to add the size of the instruction to get past this instruction
	newPC := context.PC + bytecode.InstructionSize(bytecode.JUMP) + offset

	// Check if the new PC is valid
	if newPC < 0 || newPC >= len(method.Bytecodes) {
		return false, fmt.Errorf("jump target out of bounds: %d", newPC)
	}

	// Set the PC to the new position
	context.PC = newPC

	// Skip the normal PC increment
	return true, nil
}

// ExecuteJumpIfTrue executes the JUMP_IF_TRUE bytecode
func (vm *VM) ExecuteJumpIfTrue(context *Context) (bool, error) {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the jump offset (4 bytes)
	if context.PC+1 >= len(method.Bytecodes) {
		return false, fmt.Errorf("jump offset out of bounds")
	}
	offset := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))

	// Pop the condition from the stack
	condition := context.Pop()

	// If the condition is true, jump by the offset
	isTrue := condition.IsTrue()
	if isTrue {
		// The offset is relative to the current instruction
		// We need to add the size of the instruction to get past this instruction
		newPC := context.PC + bytecode.InstructionSize(bytecode.JUMP_IF_TRUE) + offset
		// Check if the new PC is valid
		if newPC < 0 || newPC >= len(method.Bytecodes) {
			return false, fmt.Errorf("jump target out of bounds: %d", newPC)
		}

		// Set the PC to the new position
		context.PC = newPC
		return true, nil
	}

	return false, nil
}

// ExecuteJumpIfFalse executes the JUMP_IF_FALSE bytecode
func (vm *VM) ExecuteJumpIfFalse(context *Context) (bool, error) {
	// Get the method
	method := pile.ObjectToMethod(context.Method)

	// Get the jump offset (4 bytes)
	if context.PC+1 >= len(method.Bytecodes) {
		return false, fmt.Errorf("jump offset out of bounds")
	}
	offset := int(binary.BigEndian.Uint32(method.Bytecodes[context.PC+1:]))

	// Pop the condition from the stack
	condition := context.Pop()

	// If the condition is false, jump by the offset
	isTrue := condition.IsTrue()
	if !isTrue {
		// The offset is relative to the current instruction
		// We need to add the size of the instruction to get past this instruction
		newPC := context.PC + bytecode.InstructionSize(bytecode.JUMP_IF_FALSE) + offset
		// Check if the new PC is valid
		if newPC < 0 || newPC >= len(method.Bytecodes) {
			return false, fmt.Errorf("jump target out of bounds: %d", newPC)
		}

		// Set the PC to the new position
		context.PC = newPC
		return true, nil
	}

	return false, nil
}

// ExecutePop executes the POP bytecode
func (vm *VM) ExecutePop(context *Context) error {
	context.Pop()
	return nil
}

// ExecuteDuplicate executes the DUPLICATE bytecode
func (vm *VM) ExecuteDuplicate(context *Context) error {
	value := context.Top()
	context.Push(value)
	return nil
}
