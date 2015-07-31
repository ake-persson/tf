package template

import (
	"testing"
)

func Test_Last(t *testing.T) {
	a := []string{"1", "2"}
	if r, _ := Last(1, a); r != true {
		t.Error("Last didn't return expected result")
	} else {
		t.Log("Last test passes")
	}

	a2 := []int{1, 2}
	if r, _ := Last(1, a2); r != true {
		t.Error("Last didn't return expected result")
	} else {
		t.Log("Last test passes")
	}

	a3 := 1
	if r, _ := Last(1, a3); r != false {
		t.Error("Last didn't return expected result")
	} else {
		t.Log("Last test passes")
	}

}
