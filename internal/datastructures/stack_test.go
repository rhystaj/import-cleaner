package datastructures_test

import (
	"importcleaner/internal/datastructures"
	"testing"
)

func Test(t *testing.T) {
	testElements := []string{
		"a value",
		"another value",
		"yet another value",
	}

	testStack := datastructures.CreateNewStack[string]()

	for _, e := range testElements {
		testStack.Push(&e)
	}

	for i := len(testElements) - 1; i >= 0; i-- {
		peekedValue := *testStack.Peek()
		if peekedValue != testElements[i] {
			t.Errorf("Expected '%s' from peek, but got '%s'", testElements[i], peekedValue)
		}

		poppedValue := *testStack.Pop()
		if poppedValue != testElements[i] {
			t.Errorf("Expected '%s' from pop, but got '%s'", testElements[i], peekedValue)
		}
	}

	peekedValue := testStack.Peek()
	if peekedValue != nil {
		t.Error("Expected a nil value from peek.")
	}

	poppedValue := testStack.Peek()
	if poppedValue != nil {
		t.Error("Expected a nil value from pop.")
	}
}
