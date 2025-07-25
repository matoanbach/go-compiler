package vm

import (
	"fmt"
	"testing"
	"turtle/ast"
	"turtle/compiler"
	"turtle/lexer"
	"turtle/object"
	"turtle/parser"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"2 * 3", 6},
		{"8 / 4", 2},
		{"50 - 2 - 2 + 4", 50},
		{"5 + 5 + 2 + 1", 13},
		{"5 * 2 * 2 * 2", 40},
		{"5 * (2 + 2)", 20},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	runVmTests(t, tests)
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		// get the ast program
		program := parse(tt.input)

		// init the compiler
		comp := compiler.New()

		// run the compiler on the ast program
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler err: %s", err)
		}
		// init the vm (not implemented yet)
		vm := New(comp.Bytecode())
		fmt.Printf("runVmTests: %q\n", vm.instructions)
		// run the vm
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm err: %s", err)
		}

		// get the first element in the stack
		// stackElem := vm.StackTop()
		// testExpectedObject(t, tt.expected, stackElem)

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	}
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. expected=%t, got=%t", expected, result.Value)
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. expected=%d, got=%d", expected, result.Value)
	}

	return nil
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
