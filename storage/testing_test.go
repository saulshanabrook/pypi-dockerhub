package storage

import "testing"

func TestTestingClient(t *testing.T) {
	tc := NewTestingClient(123)
	initial, _ := tc.GetTime()

	if initial != int64(123) {
		t.Error("Expected 123, got ", initial)
	}

	tc.SetTime(456)

	second, _ := tc.GetTime()

	if second != int64(456) {
		t.Error("Expected 456, got ", second)
	}
}
