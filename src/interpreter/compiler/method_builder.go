package compiler

import (
	"encoding/binary"

	"smalltalklsp/interpreter/bytecode"
	"smalltalklsp/interpreter/classes"
	"smalltalklsp/interpreter/core"
)

// MethodBuilder provides a fluent interface for creating methods
type MethodBuilder struct {
	class          *classes.Class
	selectorName   string
	selectorObj    *core.Object
	bytecodes      []byte
	literals       []*core.Object
	tempVarNames   []string
	isPrimitive    bool
	primitiveIndex int
}

// NewMethodBuilder creates a new MethodBuilder for the given class
func NewMethodBuilder(class *classes.Class) *MethodBuilder {
	return &MethodBuilder{
		class:          class,
		bytecodes:      make([]byte, 0),
		literals:       make([]*core.Object, 0),
		tempVarNames:   make([]string, 0),
		isPrimitive:    false,
		primitiveIndex: 0,
	}
}

// Selector sets the selector for the method
func (mb *MethodBuilder) Selector(name string) *MethodBuilder {
	mb.selectorName = name
	mb.selectorObj = classes.NewSymbol(name)
	return mb
}

// Primitive marks the method as a primitive with the given index
func (mb *MethodBuilder) Primitive(index int) *MethodBuilder {
	mb.primitiveIndex = index
	mb.isPrimitive = true
	return mb
}

// AddLiterals adds multiple literals to the method
func (mb *MethodBuilder) AddLiterals(literals []*core.Object) *MethodBuilder {
	mb.literals = append(mb.literals, literals...)
	return mb
}

// AddLiteral adds a single literal to the method and returns its index
func (mb *MethodBuilder) AddLiteral(literal *core.Object) (int, *MethodBuilder) {
	index := len(mb.literals)
	mb.literals = append(mb.literals, literal)
	return index, mb
}

// TempVars adds temporary variable names to the method
func (mb *MethodBuilder) TempVars(names []string) *MethodBuilder {
	mb.tempVarNames = append(mb.tempVarNames, names...)
	return mb
}

// addUint32 adds a 32-bit unsigned integer to the bytecodes in big-endian format
func (mb *MethodBuilder) addUint32(value uint32) *MethodBuilder {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)
	mb.bytecodes = append(mb.bytecodes, bytes...)
	return mb
}

// PushLiteral adds a PUSH_LITERAL bytecode with the given literal index
func (mb *MethodBuilder) PushLiteral(index int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.PUSH_LITERAL)
	return mb.addUint32(uint32(index))
}

// PushInstanceVariable adds a PUSH_INSTANCE_VARIABLE bytecode with the given offset
func (mb *MethodBuilder) PushInstanceVariable(offset int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.PUSH_INSTANCE_VARIABLE)
	return mb.addUint32(uint32(offset))
}

// PushTemporaryVariable adds a PUSH_TEMPORARY_VARIABLE bytecode with the given offset
func (mb *MethodBuilder) PushTemporaryVariable(offset int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.PUSH_TEMPORARY_VARIABLE)
	return mb.addUint32(uint32(offset))
}

// PushSelf adds a PUSH_SELF bytecode
func (mb *MethodBuilder) PushSelf() *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.PUSH_SELF)
	return mb
}

// StoreInstanceVariable adds a STORE_INSTANCE_VARIABLE bytecode with the given offset
func (mb *MethodBuilder) StoreInstanceVariable(offset int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.STORE_INSTANCE_VARIABLE)
	return mb.addUint32(uint32(offset))
}

// StoreTemporaryVariable adds a STORE_TEMPORARY_VARIABLE bytecode with the given offset
func (mb *MethodBuilder) StoreTemporaryVariable(offset int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.STORE_TEMPORARY_VARIABLE)
	return mb.addUint32(uint32(offset))
}

// SendMessage adds a SEND_MESSAGE bytecode with the given selector index and argument count
func (mb *MethodBuilder) SendMessage(selectorIndex, argCount int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.SEND_MESSAGE)
	mb.addUint32(uint32(selectorIndex))
	return mb.addUint32(uint32(argCount))
}

// ReturnStackTop adds a RETURN_STACK_TOP bytecode
func (mb *MethodBuilder) ReturnStackTop() *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.RETURN_STACK_TOP)
	return mb
}

// Jump adds a JUMP bytecode with the given target offset
func (mb *MethodBuilder) Jump(target int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.JUMP)
	return mb.addUint32(uint32(target))
}

// JumpIfTrue adds a JUMP_IF_TRUE bytecode with the given target offset
func (mb *MethodBuilder) JumpIfTrue(target int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.JUMP_IF_TRUE)
	return mb.addUint32(uint32(target))
}

// JumpIfFalse adds a JUMP_IF_FALSE bytecode with the given target offset
func (mb *MethodBuilder) JumpIfFalse(target int) *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.JUMP_IF_FALSE)
	return mb.addUint32(uint32(target))
}

// Pop adds a POP bytecode
func (mb *MethodBuilder) Pop() *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.POP)
	return mb
}

// Duplicate adds a DUPLICATE bytecode
func (mb *MethodBuilder) Duplicate() *MethodBuilder {
	mb.bytecodes = append(mb.bytecodes, bytecode.DUPLICATE)
	return mb
}

// Go finalizes the method creation and adds it to the class's method dictionary
func (mb *MethodBuilder) Go() *core.Object {
	if mb.selectorObj == nil {
		panic("Selector not set. Call Selector() first.")
	}

	// Create the method object
	method := classes.NewMethod(mb.selectorObj, mb.class)

	// Set the method properties
	methodObj := classes.ObjectToMethod(method)
	methodObj.SetBytecodes(mb.bytecodes)
	methodObj.Literals = mb.literals
	methodObj.TempVarNames = mb.tempVarNames
	methodObj.SetPrimitive(mb.isPrimitive)
	methodObj.SetPrimitiveIndex(mb.primitiveIndex)

	// Add the method to the method dictionary
	symbolValue := classes.GetSymbolValue(mb.selectorObj)
	methodDict := mb.class.GetMethodDictionary()
	methodDict.SetEntry(symbolValue, method)

	// Reset the builder state for reuse
	mb.bytecodes = make([]byte, 0)
	mb.literals = make([]*core.Object, 0)
	mb.tempVarNames = make([]string, 0)
	mb.isPrimitive = false
	mb.primitiveIndex = 0
	// Note: We don't reset the class or selector as they might be reused

	return method
}
