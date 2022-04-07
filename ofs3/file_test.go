package ofs3

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	a, b := Encode(nil, "file.dgz")
	fmt.Println(a, b)
}
