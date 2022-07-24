package utils

import (
	"reflect"
	"testing"
)

func TestInsert(t *testing.T) {
	dataTable := []struct {
		In1    []int // Slice
		In2    int   // index
		In3    int   // value
		Expect []int
	}{
		{[]int{1, 2, 3}, 2, 1, []int{1, 2, 1, 3}},
	}

	for _, c := range dataTable {
		if reflect.DeepEqual(Insert(c.In1, c.In2, c.In3), c.Expect) {
			t.Log("utils.Min PASS")
		} else {
			t.Error("utils.Min FAIL")
		}
	}
}
