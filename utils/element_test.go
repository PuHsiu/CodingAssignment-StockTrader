package utils

import "testing"

func TestMin(t *testing.T) {
	dataTable := []struct {
		In1    int
		In2    int
		Expect int
	}{
		{1, 2, 1},
		{1, 0, 0},
		{0, 0, 0},
		{1, 1, 1},
	}

	for _, c := range dataTable {
		if Min(c.In1, c.In2) == c.Expect {
			t.Log("utils.Min PASS")
		} else {
			t.Error("utils.Min FAIL")
		}
	}
}
