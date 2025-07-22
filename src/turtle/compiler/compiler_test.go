package compiler

import (
	"fmt"
	"testing"
	"turtle/ast"
	"turtle/code"
	"turtle/lexer"
	"turtle/object"
	"turtle/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	// make some test cases here
	tests := []compilerTestCase{
		{
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
			},
		},
	}
	// run the real tests on those cases
	runCompilerTests(t, tests)
}

func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	// for each test case
	// init the Compiler
	// parse to get the AST
	// run the compiler
	// get the bytecode

	// call testInstructions() to check the generated instructions
	// if err, then return t.Fatalf("testInstructions failed: %s", err)
	// call testConstants() to check the generated constants
	// if err, then return t.Fatalf("testInstructions failed: %s", err)
	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("testInstructions failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("testConstants failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	output := code.Instructions{}

	for _, ins := range instructions {
		output = append(output, ins...)
	}

	return output
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)
	if len(concatted) != len(actual) {
		return fmt.Errorf("wrong number of instructions. expected=%q, got=%q", concatted, actual)
	}

	for i, ins := range concatted {
		if ins != actual[i] {
			return fmt.Errorf("wrong instruction at %d. expected=%v, got=%v", i, concatted[i], expected[i])
		}
	}

	return nil
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. expected=%d, got=%d", len(expected), len(actual))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		}
	}

	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)

	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if expected != result.Value {
		return fmt.Errorf("object has the wrong value. expected=%d, got=%d", expected, result.Value)
	}

	return nil
}
