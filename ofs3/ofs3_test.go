package ofs3

import (
	"fmt"
	"github.com/go-restruct/restruct"
	"os"
	"path"
	"testing"
)

func TestOpenOFS3(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "../data/ofs3"
	input := path.Join(dir, "005.bin")
	outputDir := path.Join(dir, "005_out")

	data, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	ofs := OpenOFS3(data, outputDir)
	fmt.Println(ofs)
}
