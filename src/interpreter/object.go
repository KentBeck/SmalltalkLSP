package main

import (
	"fmt"
)

// ObjectType represents the type of a Smalltalk object
type ObjectType int

const (
	OBJ_INTEGER ObjectType = iota
	OBJ_BOOLEAN
	OBJ_NIL
	OBJ_STRING
	OBJ_ARRAY
	OBJ_DICTIONARY
	OBJ_BLOCK
	OBJ_INSTANCE
	OBJ_CLASS
	OBJ_METHOD
	OBJ_SYMBOL
)

// Object represents a Smalltalk object
type Object struct {
	Type             ObjectType
	Class            *Object
	IntegerValue     int64
	BooleanValue     bool
	StringValue      string
	SymbolValue      string
	InstanceVars     []*Object // Instance variables stored by index
	Elements         []*Object
	Entries          map[string]*Object
	Method           *Method
	Bytecodes        []byte
	Literals         []*Object
	Selector         *Object
	SuperClass       *Object
	InstanceVarNames []string
	Moved            bool    // Used for garbage collection
	ForwardingPtr    *Object // Used for garbage collection
}

const METHOD_DICTIONARY_IV = 0

// Method represents a Smalltalk method
type Method struct {
	Bytecodes      []byte
	Literals       []*Object
	Selector       *Object
	Class          *Object
	TempVarNames   []string
	IsPrimitive    bool
	PrimitiveIndex int
}

// NewInteger creates a new integer object
func NewInteger(value int64) *Object {
	return &Object{
		Type:         OBJ_INTEGER,
		IntegerValue: value,
	}
}

// NewBoolean creates a new boolean object
func NewBoolean(value bool) *Object {
	return &Object{
		Type:         OBJ_BOOLEAN,
		BooleanValue: value,
	}
}

// NewNil creates a new nil object
func NewNil() *Object {
	return &Object{
		Type: OBJ_NIL,
	}
}

// NewString creates a new string object
func NewString(value string) *Object {
	return &Object{
		Type:        OBJ_STRING,
		StringValue: value,
	}
}

// NewSymbol creates a new symbol object
func NewSymbol(value string) *Object {
	return &Object{
		Type:        OBJ_SYMBOL,
		SymbolValue: value,
	}
}

// NewArray creates a new array object
func NewArray(size int) *Object {
	return &Object{
		Type:     OBJ_ARRAY,
		Elements: make([]*Object, size),
	}
}

// NewDictionary creates a new dictionary object
func NewDictionary() *Object {
	return &Object{
		Type:    OBJ_DICTIONARY,
		Entries: make(map[string]*Object),
	}
}

// NewInstance creates a new instance of a class
func NewInstance(class *Object) *Object {
	// Initialize instance variables array with nil values
	instVarsSize := 0
	if class != nil && len(class.InstanceVarNames) > 0 {
		instVarsSize = len(class.InstanceVarNames)
	}
	instVars := make([]*Object, instVarsSize)
	for i := range instVars {
		instVars[i] = NewNil()
	}

	return &Object{
		Type:         OBJ_INSTANCE,
		Class:        class,
		InstanceVars: instVars,
	}
}

// NewClass creates a new class object
func NewClass(name string, superClass *Object) *Object {
	// For classes, we need a special instance variable for the method dictionary
	// We'll store it at index 0
	instVars := make([]*Object, 1)
	instVars[0] = NewDictionary() // methodDict at index 0

	return &Object{
		Type:             OBJ_CLASS,
		SymbolValue:      name,
		SuperClass:       superClass,
		InstanceVarNames: make([]string, 0),
		InstanceVars:     instVars,
	}
}

// NewMethod creates a new method object
func NewMethod(selector *Object, class *Object) *Object {
	method := &Method{
		Bytecodes:    make([]byte, 0),
		Literals:     make([]*Object, 0),
		Selector:     selector,
		Class:        class,
		TempVarNames: make([]string, 0),
	}

	return &Object{
		Type:   OBJ_METHOD,
		Method: method,
	}
}

// IsTrue returns true if the object is considered true in Smalltalk
func (o *Object) IsTrue() bool {
	if o.Type == OBJ_BOOLEAN {
		return o.BooleanValue
	}
	return o.Type != OBJ_NIL
}

// String returns a string representation of the object
func (o *Object) String() string {
	switch o.Type {
	case OBJ_INTEGER:
		return fmt.Sprintf("%d", o.IntegerValue)
	case OBJ_BOOLEAN:
		if o.BooleanValue {
			return "true"
		}
		return "false"
	case OBJ_NIL:
		return "nil"
	case OBJ_STRING:
		return fmt.Sprintf("'%s'", o.StringValue)
	case OBJ_SYMBOL:
		return fmt.Sprintf("#%s", o.SymbolValue)
	case OBJ_ARRAY:
		return fmt.Sprintf("Array(%d)", len(o.Elements))
	case OBJ_DICTIONARY:
		return fmt.Sprintf("Dictionary(%d)", len(o.Entries))
	case OBJ_INSTANCE:
		if o.Class != nil {
			return fmt.Sprintf("a %s", o.Class.SymbolValue)
		}
		return "an Object"
	case OBJ_CLASS:
		return fmt.Sprintf("Class %s", o.SymbolValue)
	case OBJ_METHOD:
		if o.Method.Selector != nil {
			return fmt.Sprintf("Method %s", o.Method.Selector.SymbolValue)
		}
		return "a Method"
	default:
		return "Unknown object"
	}
}

// GetInstanceVarByIndex gets an instance variable by index
func (o *Object) GetInstanceVarByIndex(index int) *Object {
	if index < 0 || index >= len(o.InstanceVars) {
		return NewNil()
	}

	return o.InstanceVars[index]
}

// SetInstanceVarByIndex sets an instance variable by index
func (o *Object) SetInstanceVarByIndex(index int, value *Object) {
	if index < 0 || index >= len(o.InstanceVars) {
		return
	}

	o.InstanceVars[index] = value
}

// GetMethodDict gets the method dictionary for a class
func (o *Object) GetMethodDict() *Object {
	if o.Type != OBJ_CLASS || len(o.InstanceVars) == 0 {
		return NewNil()
	}

	// Method dictionary is stored at index 0 for classes
	return o.InstanceVars[METHOD_DICTIONARY_IV]
}
