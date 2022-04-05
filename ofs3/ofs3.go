package ofs3

import (
	"TMGS3Tools/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
	"os"
	"path"
	"strconv"
)

var ShowLog = true

type Header struct {
	Magic   string `struct:"[4]byte"`
	Length  int    `struct:"int32"`
	Type    int    `struct:"int16"`
	Padding int    `struct:"int8"`
	SubType int    `struct:"int8"`
	RawSize int    `struct:"int32"` // 数据段大小
	Count   int    `struct:"int32"`
}

type OFS3 struct {
	Header
	Files []*File
}
type File struct {
	Offset     int //`struct:"int32"`
	Size       int //`struct:"int32"`
	NameOffset int //`struct:"int32"`
	Name       string
	FilePath   string
	*OFS3
}

func OpenOFS3(data []byte, dir string) *OFS3 {
	if string(data[0:4]) != "OFS3" {
		fmt.Println("不是OFS3文件")
		return nil
	}
	var err error
	if !utils.DirExists(dir) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("创建文件夹失败 %v\n", dir)
			return nil
		}
	}
	ofs3 := &OFS3{}
	err = restruct.Unpack(data, binary.LittleEndian, &ofs3.Header)
	if err != nil {
		fmt.Printf("解析OFS3文件头失败 %v\n", err)
		return nil
	}
	nameStr := bytes.NewBuffer(nil)
	// 0x10(Header.Length) + 4(Header.Count)
	offset := ofs3.Length + 4
	ofs3.Files = make([]*File, ofs3.Count)
	for i, _ := range ofs3.Files {
		ofs3.Files[i] = &File{}
	}
	for i, file := range ofs3.Files {
		// 文件数据偏移，默认不含Header.Length
		file.Offset = utils.ReadUInt32(data[offset:offset+4]) + ofs3.Length
		offset += 4
		// 文件大小，SubType=1时为0 ？
		file.Size = utils.ReadUInt32(data[offset : offset+4])
		offset += 4
		// Type == 2 ,含文件名偏移
		if ofs3.Type == 2 {
			// 文件名偏移，默认不含Header.Length
			file.NameOffset = utils.ReadUInt32(data[offset:offset+4]) + ofs3.Length
			offset += 4
			nameOffset := file.NameOffset
			nameStr.Reset()
			for offset < len(data) && data[nameOffset] != 0 {
				nameStr.WriteByte(data[nameOffset])
				nameOffset++
			}
			file.Name = nameStr.String()
		} else {
			file.Name = strconv.Itoa(i)
		}
		file.FilePath = path.Join(dir, file.Name)
	}

	// TODO: 未测试
	// SubType == 1 不含File.Size，需要计算?
	if ofs3.SubType == 1 {
		for i := 0; i <= ofs3.Count; i++ {
			if i == ofs3.Count {
				ofs3.Files[i-1].Size = len(data) - ofs3.Files[i-1].Offset
			} else if i > 0 {
				ofs3.Files[i-1].Size = ofs3.Files[i].Offset - ofs3.Files[i-1].Offset
			}
		}
	}
	// 开始处理数据段

	for _, file := range ofs3.Files {
		// 子文件数据
		subData := data[file.Offset : file.Offset+file.Size]
		if string(subData[0:4]) == "OFS3" {
			// 子文件为OFS3，递归解析
			if ShowLog {
				fmt.Printf("OFS3 %v\n", file.FilePath)
			}
			file.OFS3 = OpenOFS3(subData, file.FilePath)
		} else {
			// 非OFS3，一般文件
			err = os.WriteFile(file.FilePath, subData, os.ModePerm)
			if err != nil {
				fmt.Printf("Error 写入[%v]文件失败！%v\n", file.FilePath, err)
			} else if ShowLog {
				fmt.Printf("文件 %v\n", *file)
			}
		}
	}
	return ofs3

}
