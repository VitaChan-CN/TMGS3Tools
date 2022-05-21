package ofs3

import (
	"compress/gzip"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	data, _ := os.ReadFile("../data/ofs3/script/scp001.dir/WAL_00_000.evd.o")
	fmt.Println(len(data))
	data2 := GzEncode(data, &gzip.Header{
		ModTime: time.Date(2011, 12, 8, 3, 59, 15, 0, time.Local),
		Name:    "WAL_00_000.evsc",
		OS:      11,
	})
	os.WriteFile("../data/ofs3/script/scp001.dir/WAL_00_000.evd.e.gz", data2, os.ModePerm)
	fmt.Println(len(data2))
}

func TestDecode(t *testing.T) {
	data, _ := os.ReadFile("../data/ofs3/script/scp001.dir/WAL_00_000.evd.e.gz")
	data3, h := GzDecode(data, false)
	fmt.Println(len(data3), h)
	os.WriteFile("../data/ofs3/script/scp001.dir/WAL_00_000.evd.e", data3, os.ModePerm)
}
