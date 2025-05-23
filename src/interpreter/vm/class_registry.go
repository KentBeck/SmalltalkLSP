package vm

import (
	"smalltalklsp/interpreter/pile"
)

// ClassType enum for standard class types
type ClassType string

const (
	Object          ClassType = "Object"
	UndefinedObject ClassType = "UndefinedObject" // nil class
	True            ClassType = "True"
	False           ClassType = "False"
	Integer         ClassType = "Integer"
	Float           ClassType = "Float"
	String          ClassType = "String"
	Symbol          ClassType = "Symbol"
	Array           ClassType = "Array"
	Dictionary      ClassType = "Dictionary"
	ByteArray       ClassType = "ByteArray"
	Block           ClassType = "Block"
	Method          ClassType = "Method"
	Class           ClassType = "Class"
	ContextClass    ClassType = "Context"
	Exception       ClassType = "Exception"
)

// ClassRegistry holds references to all standard Smalltalk classes
// This centralizes class management and provides a cleaner way to
// access classes than having many separate fields in the VM struct
type ClassRegistry struct {
	// Map of registered classes by ClassType
	classesByType map[ClassType]*pile.Class
	
	// Map of registered classes by name (for non-standard classes)
	classesByName map[string]*pile.Class
}

// NewClassRegistry creates a new class registry
func NewClassRegistry() *ClassRegistry {
	return &ClassRegistry{
		classesByType: make(map[ClassType]*pile.Class),
		classesByName: make(map[string]*pile.Class),
	}
}

// Register adds a class to the registry
func (r *ClassRegistry) Register(classType ClassType, class *pile.Class) {
	if class == nil {
		return
	}
	
	r.classesByType[classType] = class
	r.classesByName[class.Name] = class
}


// Get retrieves a class by type
func (r *ClassRegistry) Get(classType ClassType) *pile.Class {
	return r.classesByType[classType]
}

// GetByName retrieves a class by name
func (r *ClassRegistry) GetByName(name string) *pile.Class {
	return r.classesByName[name]
}

