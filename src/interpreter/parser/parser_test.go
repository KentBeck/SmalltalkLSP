package parser

import (
	"testing"

	"smalltalklsp/interpreter/ast"
	"smalltalklsp/interpreter/classes"
)

// TestParseYourself tests parsing the method "yourself ^self"
func TestParseYourself(t *testing.T) {
	// Create a class
	objectClass := classes.NewClass("Object", nil)

	// Create a parser
	p := NewParser("yourself ^self", classes.ClassToObject(objectClass))

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "yourself" {
		t.Errorf("Expected method selector to be 'yourself', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(methodNode.Parameters))
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(methodNode.Temporaries))
	}

	// Check the method body
	returnNode, ok := methodNode.Body.(*ast.ReturnNode)
	if !ok {
		t.Fatalf("Expected return node, got %T", methodNode.Body)
	}

	// Check the return expression
	_, ok = returnNode.Expression.(*ast.SelfNode)
	if !ok {
		t.Fatalf("Expected self node, got %T", returnNode.Expression)
	}

	// Check the method class
	if methodNode.Class != classes.ClassToObject(objectClass) {
		t.Errorf("Expected method class to be %v, got %v", classes.ClassToObject(objectClass), methodNode.Class)
	}
}

// TestParseAdd tests parsing the method "+ aNumber ^self + aNumber"
func TestParseAdd(t *testing.T) {
	// Create a class
	objectClass := classes.NewClass("Object", nil)
	integerClass := classes.NewClass("Integer", objectClass)

	// Create a parser
	p := NewParser("+ aNumber ^self + aNumber", classes.ClassToObject(integerClass))

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "+" {
		t.Errorf("Expected method selector to be '+', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(methodNode.Parameters))
	} else if methodNode.Parameters[0] != "aNumber" {
		t.Errorf("Expected parameter to be 'aNumber', got '%s'", methodNode.Parameters[0])
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(methodNode.Temporaries))
	}

	// Check the method body
	returnNode, ok := methodNode.Body.(*ast.ReturnNode)
	if !ok {
		t.Fatalf("Expected return node, got %T", methodNode.Body)
	}

	// Check the return expression
	messageSendNode, ok := returnNode.Expression.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected message send node, got %T", returnNode.Expression)
	}

	// Check the message receiver
	_, ok = messageSendNode.Receiver.(*ast.SelfNode)
	if !ok {
		t.Fatalf("Expected self node, got %T", messageSendNode.Receiver)
	}

	// Check the message selector
	if messageSendNode.Selector != "+" {
		t.Errorf("Expected message selector to be '+', got '%s'", messageSendNode.Selector)
	}

	// Check the message arguments
	if len(messageSendNode.Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(messageSendNode.Arguments))
	} else {
		// Check the argument
		variableNode, ok := messageSendNode.Arguments[0].(*ast.VariableNode)
		if !ok {
			t.Fatalf("Expected variable node, got %T", messageSendNode.Arguments[0])
		}

		if variableNode.Name != "aNumber" {
			t.Errorf("Expected variable name to be 'aNumber', got '%s'", variableNode.Name)
		}
	}

	// Check the method class
	if methodNode.Class != classes.ClassToObject(integerClass) {
		t.Errorf("Expected method class to be %v, got %v", classes.ClassToObject(integerClass), methodNode.Class)
	}
}

// TestParseWithTemporaries tests parsing a method with temporary variables
func TestParseWithTemporaries(t *testing.T) {
	// Create a class
	objectClass := classes.NewClass("Object", nil)

	// Create a parser
	p := NewParser("factorial | temp | ^temp", classes.ClassToObject(objectClass))

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "factorial" {
		t.Errorf("Expected method selector to be 'factorial', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 0 {
		t.Errorf("Expected 0 parameters, got %d", len(methodNode.Parameters))
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 1 {
		t.Errorf("Expected 1 temporary, got %d", len(methodNode.Temporaries))
	} else if methodNode.Temporaries[0] != "temp" {
		t.Errorf("Expected temporary to be 'temp', got '%s'", methodNode.Temporaries[0])
	}

	// Check the method class
	if methodNode.Class != classes.ClassToObject(objectClass) {
		t.Errorf("Expected method class to be %v, got %v", classes.ClassToObject(objectClass), methodNode.Class)
	}
}

// TestParseWithBlock tests parsing a method with a block
func TestParseWithBlock(t *testing.T) {
	// Create a class
	objectClass := classes.NewClass("Object", nil)

	// Create a parser
	p := NewParser("do: aBlock ^aBlock value", classes.ClassToObject(objectClass))

	// Parse the method
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("Error parsing method: %v", err)
	}

	// Check that the node is a method node
	methodNode, ok := node.(*ast.MethodNode)
	if !ok {
		t.Fatalf("Expected method node, got %T", node)
	}

	// Check the method selector
	if methodNode.Selector != "do:" {
		t.Errorf("Expected method selector to be 'do:', got '%s'", methodNode.Selector)
	}

	// Check the method parameters
	if len(methodNode.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(methodNode.Parameters))
	} else if methodNode.Parameters[0] != "aBlock" {
		t.Errorf("Expected parameter to be 'aBlock', got '%s'", methodNode.Parameters[0])
	}

	// Check the method temporaries
	if len(methodNode.Temporaries) != 0 {
		t.Errorf("Expected 0 temporaries, got %d", len(methodNode.Temporaries))
	}

	// Check the method body
	returnNode, ok := methodNode.Body.(*ast.ReturnNode)
	if !ok {
		t.Fatalf("Expected return node, got %T", methodNode.Body)
	}

	// Check the return expression
	messageSendNode, ok := returnNode.Expression.(*ast.MessageSendNode)
	if !ok {
		t.Fatalf("Expected message send node, got %T", returnNode.Expression)
	}

	// Check the message receiver
	variableNode, ok := messageSendNode.Receiver.(*ast.VariableNode)
	if !ok {
		t.Fatalf("Expected variable node, got %T", messageSendNode.Receiver)
	}

	if variableNode.Name != "aBlock" {
		t.Errorf("Expected variable name to be 'aBlock', got '%s'", variableNode.Name)
	}

	// Check the message selector
	if messageSendNode.Selector != "value" {
		t.Errorf("Expected message selector to be 'value', got '%s'", messageSendNode.Selector)
	}

	// Check the message arguments
	if len(messageSendNode.Arguments) != 0 {
		t.Errorf("Expected 0 arguments, got %d", len(messageSendNode.Arguments))
	}

	// Check the method class
	if methodNode.Class != classes.ClassToObject(objectClass) {
		t.Errorf("Expected method class to be %v, got %v", classes.ClassToObject(objectClass), methodNode.Class)
	}
}
