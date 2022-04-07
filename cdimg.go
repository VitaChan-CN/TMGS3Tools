package main

import (
	"TMGS3Tools/ofs3"
	"TMGS3Tools/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
	"io"
	"os"
	"path"
)

type DFI struct {
	Count    int    `struct:"-"`
	Size     int    `struct:"-"`
	Magic    string `struct:"[4]byte"`
	Unknown1 int    `struct:"int32"`
	ImgSize  int    `struct:"int32"`
	Unknown3 int    `struct:"int32"`
	Nodes    []Node `struct:"size=Count"`
}

type Node struct {
	IsDirN    int    `struct:"int16"`
	FileCount int    `struct:"int16"`
	Unknown1  int    `struct:"int32"`
	Offset    int    `struct:"int32"`
	Length    int    `struct:"int32"`
	FileName  string `struct:"-"`
	FilePath  string `struct:"-"` // 完整路径
}

func (n *Node) IsDir() bool {
	return n.IsDirN != 0
}

func LoadIdx(idxFile string) *DFI {
	idx, _ := os.ReadFile(idxFile)
	offset := 0x10
	dfi := &DFI{
		Size: len(idx),
	}
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

func (d *DFI) SetDir(dir string, isInput bool) {
	if !isInput && !utils.DirExists(dir) {
		// output
		os.Mkdir(dir, os.ModePerm)
	}
	lastDir := ""
	for i, node := range d.Nodes {
		if node.IsDir() {
			d.Nodes[i].FilePath = path.Join(dir, node.FileName)
			lastDir = d.Nodes[i].FilePath
		} else {
			d.Nodes[i].FilePath = path.Join(lastDir, node.FileName)
		}
	}
}

func (d *DFI) LoadImg(imgFile string, openOfs3 bool, gz bool) {

	f, _ := os.Open(imgFile)
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()
	for _, node := range d.Nodes {
		if node.IsDir() {
			if ShowLog {
				fmt.Printf("文件夹 %v\n", node)
			}
			if !utils.DirExists(node.FilePath) {
				os.Mkdir(node.FilePath, os.ModePerm)
			}
		} else if node.Offset+node.Length <= int(size) {

			data := make([]byte, node.Length)
			f.ReadAt(data, int64(node.Offset))

			if openOfs3 && len(data) > 4 && string(data[0:4]) == "OFS3" {
				if ShowLog {
					fmt.Printf("OFS3文件 %v\n", node)
				}
				ofs := ofs3.OpenOFS3(data, node.FilePath)
				ofs.WriteFile(data, node.FilePath, gz)
				if ShowLog {
					fmt.Printf("\t %v\n", ofs.Header)
				}
			} else {
				if ShowLog {
					fmt.Printf("文件 %v\n", node)
				}
				if gz {
					data = ofs3.Decode(data)
					node.FilePath = ofs3.DecodeName(data, node.FilePath)
				}
				tf, _ := os.Create(node.FilePath)
				tf.Write(data)
				tf.Close()
			}

		} else {
			if ShowLog {
				fmt.Printf("不在此文件中 %v\n", node)
			}
		}
	}

}

func (d *DFI) ReBuildImg(imgFile, outputFile string, appendMode bool) {
	f, _ := os.Open(imgFile)
	defer f.Close()
	info, _ := f.Stat()
	size := info.Size()

	out, _ := os.Create(outputFile)
	defer out.Close()

	endIndex := len(d.Nodes) // 因为要排除INSTALL中的数据，用于标记img最后一个文件的下标(不含)

	for i, node := range d.Nodes {
		if node.Offset+node.Length > int(size) {
			endIndex = i
			break
		}
	}

	var data []byte
	offset := 0
	if appendMode {
		if ShowLog {
			fmt.Printf("追加模式，正在复制原文件...\n根据硬盘读写速度可能需要一段时间\n")
		}
		io.Copy(out, f)
		offset = int(size)
		offset = utils.AlignUp(offset, 2048)
	}
	for i := 0; i < endIndex; i++ {
		if d.Nodes[i].IsDir() {
			continue
		}
		if !utils.FileExists(d.Nodes[i].FilePath) {
			if appendMode {
				continue
			}
			if ShowLog {
				fmt.Printf("文件不存在，将使用原数据。%v\n", d.Nodes[i].FilePath)
			}
			data = make([]byte, d.Nodes[i].Length)
			f.ReadAt(data, int64(d.Nodes[i].Offset))
		} else {
			data, _ = os.ReadFile(d.Nodes[i].FilePath)
			if appendMode {
				dataSrc := make([]byte, d.Nodes[i].Length)
				f.ReadAt(dataSrc, int64(d.Nodes[i].Offset))
				if utils.MD5(data) == utils.MD5(dataSrc) {
					if ShowLog {
						fmt.Printf("文件未更改，跳过。%v\n", d.Nodes[i].FilePath)
					}
					continue
				}
			}
			if ShowLog {
				fmt.Printf("写入文件 %v\n", d.Nodes[i].FilePath)
			}
		}
		out.WriteAt(data, int64(offset))
		d.Nodes[i].Offset = offset
		d.Nodes[i].Length = len(data)
		offset += len(data)
		offset = utils.AlignUp(offset, 2048)
	}
	// 字节对齐
	out.Truncate(int64(offset))

	d.ImgSize = offset / 2048
	// INSTALL 偏移修改
	for i := endIndex; i < len(d.Nodes); i++ {
		if d.Nodes[i].IsDir() {
			continue
		}
		d.Nodes[i].Offset = offset
		offset += d.Nodes[i].Length
		offset = utils.AlignUp(offset, 2048)
	}

}
func (d *DFI) SaveIdx(outputFile string) {
	for i, _ := range d.Nodes {
		d.Nodes[i].Offset = d.Nodes[i].Offset / 2048
	}
	data, _ := restruct.Pack(binary.LittleEndian, d)
	f, _ := os.Create(outputFile)
	defer f.Close()
	f.Write(data)
	for _, node := range d.Nodes {
		f.WriteString(node.FileName)
		f.Write([]byte{0})
	}
	f.Truncate(int64(d.Size))

}
