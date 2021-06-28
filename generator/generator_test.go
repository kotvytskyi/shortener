package generator

import (
	"fmt"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	var lengthTests = []struct {
		len    uint
		gotLen uint
	}{
		{5, 5},
		{3, 3},
	}

	for _, test := range lengthTests {
		t.Run(fmt.Sprintf("L=%d", test.len), func(t *testing.T) {
			got := generate(test.len)
			length := len(got)
			if length != int(test.gotLen) {
				t.Errorf("len of Generate(%d) = %d.", test.len, length)
			}
		})
	}
}
