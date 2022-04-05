package main

import (
	"fmt"
	"github.com/go-restruct/restruct"
	"testing"
)

func TestLoadMes(t *testing.T) {
	restruct.EnableExprBeta()
	dfi := LoadIdx("data/01/a.idx")
	for _, node := range dfi.Nodes {
		if node.IsDir {
			fmt.Printf("%v\n", node)
		} else {
			fmt.Printf("----%v\n", node)
		}
	}

}

func TestLoadIdx(t *testing.T) {
	restruct.EnableExprBeta()
	dir := "data/01/"
	dfi := LoadIdx(dir + "a.idx")
	dfi.LoadImg(dir+"a.img", dir+"output")

}
