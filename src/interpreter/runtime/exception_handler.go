package runtime

import (
	"smalltalklsp/interpreter/pile"
)

// ExceptionHandler represents an exception handler
type ExceptionHandler struct {
	ExceptionClass *pile.Object
	HandlerBlock   *pile.Object
	NextHandler    *ExceptionHandler
}

// CurrentExceptionHandler is the current active exception handler
var CurrentExceptionHandler *ExceptionHandler

// IsKindOf checks if an object is an instance of a class or one of its subclasses
// This is a simplified implementation that just checks if the classes are the same
// In a real implementation, we would check the class hierarchy
func IsKindOf(obj *pile.Object, class *pile.Object) bool {
	return obj.Class() == class
}

// SignalException signals an exception
// If there's a handler for the exception, it will be executed
// Otherwise, it will panic with the exception
func SignalException(exception *pile.Object) *pile.Object {
	// If there's no handler, just panic with the exception
	if CurrentExceptionHandler == nil {
		panic(exception)
	}

	// Find a handler for this exception
	handler := CurrentExceptionHandler
	for handler != nil {
		if IsKindOf(exception, handler.ExceptionClass) {
			// Found a handler, execute it
			return ExecuteBlock(handler.HandlerBlock, []*pile.Object{exception})
		}
		handler = handler.NextHandler
	}

	// No handler found, panic with the exception
	panic(exception)
	return nil // This will never be reached, but Go requires a return value
}