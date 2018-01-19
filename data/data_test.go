package usgazetteer

import "testing"

func TestVars(t *testing.T) {
	if len(States) == 0 {
		t.Error("States are empty")
	}

	if len(Counties) == 0 {
		t.Error("Counties are empty")
	}
}
