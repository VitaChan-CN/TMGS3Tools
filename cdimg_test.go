package main

import (
	"TMGS3Tools/utils"
	"fmt"
	"github.com/go-restruct/restruct"
	"testing"
)

func TestLoadMes(t *testing.T) {
	restruct.EnableExprBeta()
	dfi := LoadIdx("data/01/a.idx")
	for _, node := range dfi.Nodes {
		if node.IsDir() {
			fmt.Printf("%v\n", node)
		} else {
			fmt.Printf("----%v\n", node)
		}
	}

}

func TestLoadIdx(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "data/01/"
	inputIdx := dir + "a.out.idx"
	inputImg := dir + "a.out.img"
	dfi := LoadIdx(inputIdx)
	dfi.SetDir(dir+"output", false)
	dfi.LoadImg(inputImg, false, false)

}

func TestDFI_ReBuildImg(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "data/01/"
	inputIdx := dir + "a.idx"
	inputImg := dir + "a.img"
	outputIdx := dir + "a.out.idx"
	outputImg := dir + "a.out.img"

	dfi := LoadIdx(inputIdx)
	dfi.SetDir(dir+"output", true)
	dfi.ReBuildImg(inputImg, outputImg, true)
	dfi.SaveIdx(outputIdx)
	fmt.Printf("%v\n%v\n", utils.MD5F(inputImg), utils.MD5F(outputImg))
	fmt.Printf("%v\n%v\n", utils.MD5F(inputIdx), utils.MD5F(outputIdx))
}
