package classes

import (
	"fmt"
	"unsafe"

	"smalltalklsp/interpreter/core"
)

// Method represents a Smalltalk method
type Method struct {
	core.Object
	Bytecodes      []byte
	Literals       []*core.Object
	Selector       *core.Object
	TempVarNames   []string
	MethodClass    *Class
	IsPrimitive    bool
	PrimitiveIndex int
}

// NewMethod creates a new method object
func NewMethod(selector *core.Object, class *Class) *core.Object {
	method := &Method{
		Object: core.Object{
			TypeField: core.OBJ_METHOD,
		},
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*core.Object, 0),
		Selector:     selector,
		MethodClass:  class,
		TempVarNames: make([]string, 0),
	}

	return MethodToObject(method)
}

// MethodToObject converts a Method to an Object
func MethodToObject(m *Method) *core.Object {
	return (*core.Object)(unsafe.Pointer(m))
}

// ObjectToMethod converts an Object to a Method
func ObjectToMethod(o *core.Object) *Method {
	if o == nil || o.Type() != core.OBJ_METHOD {
		return nil
	}
	return (*Method)(unsafe.Pointer(o))
}

// String returns a string representation of the method object
func (m *Method) String() string {
	if m.Selector != nil {
		return fmt.Sprintf("Method %s", GetSymbolValue(m.Selector))
	}
	return "a Method"
}

// GetBytecodes returns the bytecodes of the method
func (m *Method) GetBytecodes() []byte {
	return m.Bytecodes
}

// SetBytecodes sets the bytecodes of the method
func (m *Method) SetBytecodes(bytecodes []byte) {
	m.Bytecodes = bytecodes
}

// GetLiterals returns the literals of the method
func (m *Method) GetLiterals() []*core.Object {
	return m.Literals
}

// AddLiteral adds a literal to the method
func (m *Method) AddLiteral(literal *core.Object) {
	m.Literals = append(m.Literals, literal)
}

// GetSelector returns the selector of the method
func (m *Method) GetSelector() *core.Object {
	return m.Selector
}

// SetSelector sets the selector of the method
func (m *Method) SetSelector(selector *core.Object) {
	m.Selector = selector
}

// GetTempVarNames returns the temporary variable names of the method
func (m *Method) GetTempVarNames() []string {
	return m.TempVarNames
}

// AddTempVarName adds a temporary variable name to the method
func (m *Method) AddTempVarName(name string) {
	m.TempVarNames = append(m.TempVarNames, name)
}

// GetMethodClass returns the class of the method
func (m *Method) GetMethodClass() *Class {
	return m.MethodClass
}

// SetMethodClass sets the class of the method
func (m *Method) SetMethodClass(class *Class) {
	m.MethodClass = class
}

// IsPrimitiveMethod returns true if the method is a primitive
func (m *Method) IsPrimitiveMethod() bool {
	return m.IsPrimitive
}

// SetPrimitive sets whether the method is a primitive
func (m *Method) SetPrimitive(isPrimitive bool) {
	m.IsPrimitive = isPrimitive
}

// GetPrimitiveIndex returns the primitive index of the method
func (m *Method) GetPrimitiveIndex() int {
	return m.PrimitiveIndex
}

// SetPrimitiveIndex sets the primitive index of the method
func (m *Method) SetPrimitiveIndex(index int) {
	m.PrimitiveIndex = index
}
