package main

import (
	"TMGS3Tools/ofs3"
	"TMGS3Tools/utils"
	"bytes"
	"compress/gzip"
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
	IsDirN    int           `struct:"int16"`
	FileCount int           `struct:"int16"`
	Unknown1  int           `struct:"int32"`
	Offset    int           `struct:"int32"`
	Length    int           `struct:"int32"`
	FileName  string        `struct:"-"`
	FilePath  string        `struct:"-"` // 完整路径
	FileType  ofs3.FileType `struct:"-"` // 文件类型
	GzHeader  *gzip.Header  `struct:"-"`
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
	dfi.ImgSize *= 2048
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

func (d *DFI) LoadImg(imgFile, installFile string, openOfs3 bool, gz bool) {

	var f, f2 *os.File
	f, _ = os.Open(imgFile)
	defer f.Close()
	if len(installFile) > 0 {
		f2, _ = os.Open(installFile)
		defer f2.Close()
	}
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
		} else if node.Offset+node.Length <= int(size) || len(installFile) > 0 {

			data := make([]byte, node.Length)
			if node.Offset+node.Length <= int(size) {
				f.ReadAt(data, int64(node.Offset))
			} else {
				f2.ReadAt(data, int64(node.Offset)-size)
			}
			node.FileType = ofs3.GetFileType(data)
			switch node.FileType {
			case ofs3.FILE_GZ:
				if gz {
					data, node.GzHeader = ofs3.GzDecode(data, false)
					if ShowLog {
						fmt.Println(node.GzHeader)
					}
				}
				fallthrough
			case ofs3.FILE_OFS3:
				if openOfs3 {
					if ShowLog {
						fmt.Printf("OFS3文件 %v\n", node)
					}
					ofs := ofs3.OpenOFS3(data, node.FilePath)
					ofs.WriteFile(data, node.FilePath, gz)
					if ShowLog {
						fmt.Printf("\t %v\n", ofs.Header)
					}
				}
				fallthrough
			case ofs3.FILE_OTHER:
				fallthrough
			default:
				if ShowLog {
					fmt.Printf("文件 %v\n", node)
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

func (d *DFI) ReBuildImg(imgFile, outputFile, installFile string, gz bool, appendMode bool, patchOffset int) {
	var f, f2 *os.File
	var out, out2 *os.File
	f, _ = os.Open(imgFile)
	defer f.Close()
	if len(installFile) > 0 {
		f2, _ = os.Open(installFile)
		defer f2.Close()
		out2, _ = os.Create(path.Join(path.Dir(outputFile), "out_INSTALL.DAT"))
		defer out2.Close()
		appendMode = false
		patchOffset = 0
	}
	info, _ := f.Stat()
	size := info.Size()

	patch := false
	if patchOffset > 0 && appendMode && utils.FileExists(outputFile) {
		// 对output的指定位置进行打补丁
		out, _ = os.OpenFile(outputFile, os.O_RDWR, os.ModePerm)
		patch = true
	} else {
		out, _ = os.Create(outputFile)
	}

	defer out.Close()

	endIndex := len(d.Nodes)
	installIndex := endIndex
	// 排除INSTALL中的数据，用于标记img最后一个文件的下标(不含)
	for i, node := range d.Nodes {
		if node.Offset+node.Length > int(size) {
			if len(installFile) == 0 {
				endIndex = i
			} else {
				installIndex = i
			}

			break
		}
	}

	var data []byte
	offset := 0
	if appendMode {
		if patch {
			offset = patchOffset
			fmt.Println(offset)
		} else {
			if ShowLog {
				fmt.Printf("追加模式，正在复制原文件...\n根据硬盘读写速度可能需要一段时间\n")
			}
			io.Copy(out, f)
			offset = int(size)
			offset = utils.AlignUp(offset, 2048)
		}
		d.ImgSize = offset
	}
	install := false
	for i := 0; i < endIndex; i++ {
		if d.Nodes[i].IsDir() {
			continue
		}
		if i == installIndex {
			install = true
			d.ImgSize = offset
		}
		// 读取原文件
		dataSrc := make([]byte, d.Nodes[i].Length)
		if !install {
			f.ReadAt(dataSrc, int64(d.Nodes[i].Offset))
		} else {
			f2.ReadAt(dataSrc, int64(d.Nodes[i].Offset)-size)
		}
		d.Nodes[i].FileType = ofs3.GetFileType(dataSrc)

		if !utils.FileExists(d.Nodes[i].FilePath) {
			if appendMode {
				continue
			}
			if ShowLog {
				fmt.Printf("文件不存在，将使用原数据。%v\n", d.Nodes[i].FilePath)
			}
			data = dataSrc

		} else {
			data, _ = os.ReadFile(d.Nodes[i].FilePath)
			if gz && d.Nodes[i].FileType == ofs3.FILE_GZ {
				_, d.Nodes[i].GzHeader = ofs3.GzDecode(dataSrc, true)
				data = ofs3.GzEncode(data, d.Nodes[i].GzHeader)
				if ShowLog {
					fmt.Println(d.Nodes[i].GzHeader)
				}
			} else if appendMode {
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
		if !install {
			out.WriteAt(data, int64(offset))
		} else {
			// install
			out2.WriteAt(data, int64(offset-d.ImgSize))
		}

		d.Nodes[i].Offset = offset
		d.Nodes[i].Length = len(data)
		offset += len(data)
		offset = utils.AlignUp(offset, 2048)
	}
	// 字节对齐
	if install {
		out.Truncate(int64(d.ImgSize))
		out2.Truncate(int64(offset - d.ImgSize))
	} else {
		out.Truncate(int64(offset))
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
	d.ImgSize /= 2048

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
