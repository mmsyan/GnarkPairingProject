package utils

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	a1 := []int{1, 2, 3, 4}
	a2 := []int{2, 3, 4, 5}

	result := FindCommonAttributes(a1, a2, 3)

	if !reflect.DeepEqual(result, []int{2, 3, 4}) {
		t.Error("FindCommonAttributes failed")
	}
}
