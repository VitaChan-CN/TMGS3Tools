package ofs3

import (
	"fmt"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	data, _ := os.ReadFile("../data/ofs3/script/scp001.dir/WAL_0A_000.evd.dgz")
	fmt.Println(len(data))
	data2, b := Encode(data, "file.dgz")
	os.WriteFile("../data/ofs3/script/scp001.dir/WAL_0A_000.evd", data2, os.ModePerm)
	fmt.Println(b)

	fmt.Println(len(data2))
	data3 := Decode(data2)
	fmt.Println(len(data3))
	os.WriteFile("../data/ofs3/script/scp001.dir/WAL_0A_000.evd.2", data3, os.ModePerm)
}
