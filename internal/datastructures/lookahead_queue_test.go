package datastructures_test

import (
	"importcleaner/internal/datastructures"
	"testing"
)

func TestPeekStart(t *testing.T) {
	lookahead := 5
	testValues := []string{
		"0_K7mQz9xFp2Lw",
		"1_aTbV#58nRqJh",
		"2_Zx34M!pLuWs7",
		"3_fG9kD&2vNeYc",
		"4_Hj5rXm@8QzTb",
		"5_WpL3$vN9sKfd",
		"6_yQ7zM*xT4RgBe",
		"7_E2uH6%cJwMnA",
	}

	queue, err := datastructures.CreateNewLookahead[string](lookahead)
	if err != nil {
		t.Fatalf("Failed to create queue: %v", err)
	}

	if queue.PeekStart() != nil {
		t.Error("PeekStart should return nil when the queue is empty")
	}

	// Up until lookahead -1 items are pushed, the peek should return nil
	for i := 0; i < lookahead-1; i++ {
		queue.PushItem(&testValues[i])

		peekedValue := queue.PeekStart()
		if peekedValue != nil {
			t.Errorf("PeekStart should return nil, but got %s", *peekedValue)
		}
	}

	// After lookahead -1 items are pushed, the peek should return the first item, shifting it each time a new item is pushed
	expectedValueIndex := 0
	for _, value := range testValues[lookahead-1:] {
		queue.PushItem(&value)

		peekedValue := queue.PeekStart()
		expectedValue := testValues[expectedValueIndex]
		if *peekedValue != expectedValue {
			t.Errorf("PeekStart should return %s, but got %s", expectedValue, *peekedValue)
		}
		expectedValueIndex++
	}

	for {
		queue.PushItem(nil)
		peekedValue := queue.PeekStart()

		if expectedValueIndex < len(testValues) {
			// Keep adding the test values and expect the start value to shift
			expectedValue := testValues[expectedValueIndex]
			if *peekedValue != expectedValue {
				t.Errorf("PeekStart should return %s, but got %s", expectedValue, *peekedValue)
			}
			expectedValueIndex++
		} else {
			// When nil has been pushed lookahead times, the peek should return nil
			if peekedValue != nil {
				t.Errorf("PeekStart should return nil, but got %s", *peekedValue)
			}
			break
		}
	}
}

func TestPeekEnd(t *testing.T) {
	lookahead := 5
	testValues := []string{
		"K7mQz9xFp2Lw",
		"aTbV#58nRqJh",
		"Zx34M!pLuWs7",
		"fG9kD&2vNeYc",
		"Hj5rXm@8QzTb",
		"WpL3$vN9sKfd",
		"yQ7zM*xT4RgBe",
		"E2uH6%cJwMnA",
	}

	queue, err := datastructures.CreateNewLookahead[string](lookahead)
	if err != nil {
		t.Fatalf("Failed to create queue: %v", err)
	}

	if queue.PeekEnd() != nil {
		t.Error("PeekEnd should return nil when the queue is empty")
	}

	for _, value := range testValues {
		queue.PushItem(&value)
		if *queue.PeekEnd() != value {
			t.Errorf("PeekEnd should return %s, but got %s", value, *queue.PeekEnd())
		}
	}

	queue.PushItem(nil)
	if queue.PeekEnd() != nil {
		t.Errorf("PeekEnd should return nil, but got %s", *queue.PeekEnd())
	}
}
