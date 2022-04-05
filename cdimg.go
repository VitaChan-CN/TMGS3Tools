package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
	"os"
	"path"
)

type DFI struct {
	Count    int    `struct:"-"`
	Magic    string `struct:"[4]byte"`
	Unknown1 int    `struct:"int32"`
	Unknown2 int    `struct:"int32"`
	Unknown3 int    `struct:"int32"`
	Nodes    []Node `struct:"size=Count"`
}

type Node struct {
	IsDir     bool   `struct:"int16,variantbool"`
	FileCount int    `struct:"int16"`
	Unknown1  int    `struct:"int32"`
	Offset    int    `struct:"int32"`
	Length    int    `struct:"int32"`
	FileName  string `struct:"-"`
}

func LoadIdx(idxFile string) *DFI {
	idx, _ := os.ReadFile(idxFile)
	offset := 0x10
	dfi := &DFI{}
	for offset+1 < len(idx) && idx[offset] <= 0x01 && idx[offset+1] == 0x00 {
		dfi.Count++
		offset += 0x10
	}
	if ShowLog {
		fmt.Printf("文件数量：%d\n", dfi.Count)
	}
	err := restruct.Unpack(idx, binary.LittleEndian, dfi)
	if err != nil {
		panic(err)
	}
	offset = 0x10 + dfi.Count*0x10
	str := bytes.NewBuffer(nil)
	for i, _ := range dfi.Nodes {
		for idx[offset] != 0 {
			str.WriteByte(idx[offset])
			offset++
		}
		offset++
		dfi.Nodes[i].Offset *= 2048
		dfi.Nodes[i].FileName = str.String()
		str.Reset()
	}
	return dfi
}

func (d *DFI) LoadImg(imgFile string, outputDir string) {
	if !PathExists(outputDir) {
		os.Mkdir(outputDir, os.ModePerm)
	}
	lastDir := ""
	f, _ := os.Open(imgFile)
	info, _ := f.Stat()
	size := info.Size()
	for _, node := range d.Nodes {
		if node.IsDir {
			if ShowLog {
				fmt.Printf("文件夹 %v\n", node)
			}

			filename := path.Join(outputDir, node.FileName)
			if !PathExists(filename) {
				os.Mkdir(filename, os.ModePerm)
			}
			lastDir = filename
		} else if node.Offset+node.Length <= int(size) {
			if ShowLog {
				fmt.Printf("文件 %v\n", node)
			}

			filename := path.Join(lastDir, node.FileName)
			tf, _ := os.Create(filename)
			data := make([]byte, node.Length)
			f.ReadAt(data, int64(node.Offset))
			tf.Write(data)
			tf.Close()
		} else {
			if ShowLog {
				fmt.Printf("不在此文件中 %v\n", node)
			}

		}
	}

	f.Close()

}
func (d *DFI) ReBuildImg(imgFile string, inputDir string, outputFile string) {

}
