package ofs3

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"TMGS3Tools/utils"
	"github.com/go-restruct/restruct"
)

func TestOpenOFS3(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "../data/ofs3/script"
	name := "scp001"
	input := path.Join(dir, name)
	outputDir := path.Join(dir, name+".dir")
	ShowLog = true

	data, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	t1 := time.Now()
	ofs := OpenOFS3(data, outputDir)
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
	fmt.Println("==========开始写出数据==========")
	t1 = time.Now()
	ofs.WriteFile(data, outputDir, true)
	elapsed = time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
	fmt.Println(ofs)
}

func TestOFS3_ReBuild(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "../data/ofs3/script"
	name := "scp001"
	outputName := name + ".out"
	input := path.Join(dir, name)
	outputDir := path.Join(dir, name+".dir")
	output := path.Join(dir, outputName)
	ShowLog = true

	data, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	t1 := time.Now()
	ofs := OpenOFS3(data, outputDir)
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
	fmt.Println("==========开始写出数据==========")
	t1 = time.Now()

	ofs.ReBuild(data, output, true)
	elapsed = time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)

	fmt.Printf("%v\n%v\n", utils.MD5F(input), utils.MD5F(output))
}

func TestOFS3_WriteFile(t *testing.T) {

	fmt.Println(utils.AlignUp(0x11, 0x10))
}
