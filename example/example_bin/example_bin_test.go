package example_bin

import "testing"

func Test(t *testing.T) {
	res := getName("Red", "Stone")
	if res != "Red Stone Team" {
		t.Error("Error", res)
	}
}
