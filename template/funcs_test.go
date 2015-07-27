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

/*
func Test_Add2Ints_1(t *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
    if (Add2Ints(3, 4) != 7) { //try a unit test on function
        t.Error("Add2Ints did not work as expected.") // log error if it did not work as expected
    } else {
        t.Log("one test passed.") // log some info if you want
    }
}

func Test_Add2Ints_2(t *testing.T) { //test function starts with "Test" and takes a pointer to type testing.T
    t.Error("this is just hardcoded as an error.") //Indicate that this test failed and log the string as info
}
*/
