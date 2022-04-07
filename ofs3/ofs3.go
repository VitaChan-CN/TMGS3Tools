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

// OpenOFS3
//  Description
//  Param data []byte OFS3文件数据
//  Param dir string 保存到目录。进用来计算目录层级，不会创建文件夹
//  Return *OFS3
//
func OpenOFS3(data []byte, dir string) *OFS3 {
	if string(data[0:4]) != "OFS3" {
		fmt.Println("不是OFS3文件")
		return nil
	}
	var err error

	ofs3 := &OFS3{}
	err = restruct.Unpack(data, binary.LittleEndian, &ofs3.Header)
	if err != nil {
		fmt.Printf("解析OFS3文件头失败 %v\n", err)
		return nil
	}
	fmt.Println(ofs3.Header)
	nameStr := bytes.NewBuffer(nil)
	// 0x10(Header.Length) + 4(Header.Count)
	offset := ofs3.Length + 4
	ofs3.Files = make([]*File, ofs3.Count)
	for i, _ := range ofs3.Files {
		ofs3.Files[i] = &File{}
	}
	offsetMap := make(map[int]int, ofs3.Count)
	for i, file := range ofs3.Files {
		// 文件数据偏移，默认不含Header.Length
		file.Offset = utils.ReadUInt32(data[offset:offset+4]) + ofs3.Length
		if index, has := offsetMap[file.Offset]; has {
			// 指向共同offset
			ofs3.Files[i] = ofs3.Files[index]
			offset += 8
			if ofs3.Type == 2 {
				offset += 4
			}
			continue
		} else {
			offsetMap[file.Offset] = i
		}
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

	// SubType == 1 File.Size=0，需要重新计算
	if ofs3.SubType == 1 {
		for i := 0; i <= ofs3.Count; i++ {
			if i == ofs3.Count {
				ofs3.Files[i-1].Size = len(data) - ofs3.Files[i-1].Offset
			} else if i > 0 {
				ofs3.Files[i-1].Size = ofs3.Files[i].Offset - ofs3.Files[i-1].Offset
			}
		}
	}
	// 递归读取嵌套数据，不写出数据
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
			if ShowLog {
				fmt.Printf("文件 %v\n", *file)
			}
		}
	}
	return ofs3

}

// WriteFile
//  Description 根据OFS3树导出文件
//  Receiver ofs3 *OFS3
//  Param data []byte 原文件数据
//  Param dir string 所在目录，会创建文件夹
//
func (ofs3 *OFS3) WriteFile(data []byte, dir string, gz bool) {
	var err error
	if !utils.DirExists(dir) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("创建文件夹失败 %v\n", dir)
			return
		}
	}
	for _, file := range ofs3.Files {
		// 子文件数据
		subData := data[file.Offset : file.Offset+file.Size]
		if file.OFS3 != nil {
			// 子文件为OFS3，递归读取
			file.WriteFile(subData, file.FilePath, gz)
		} else {
			// 非OFS3，一般文件
			if gz {
				subData = Decode(subData)
				file.FilePath = DecodeName(data, file.FilePath)
			}
			err = os.WriteFile(file.FilePath, subData, os.ModePerm)
			if err != nil {
				fmt.Printf("Error 写入[%v]文件失败！%v\n", file.FilePath, err)
			}
		}
	}
}

// ReBuild
//  Description
//  Receiver ofs3 *OFS3
//  Param data []byte 原ofs3数据
//  Param output string
//
func (ofs3 *OFS3) ReBuild(data []byte, output string, gz bool) {
	result := ofs3.createOFS3(data)
	err := os.WriteFile(output, result, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (ofs3 *OFS3) createOFS3(srcData []byte) []byte {

	// Header 数据，共0x14
	headerData, err := restruct.Pack(binary.LittleEndian, &ofs3.Header)
	if err != nil {
		fmt.Printf("解析OFS3文件头失败 %v\n", err)
		return nil
	}
	var subFileData []byte
	offsetDataSize := 4 + ofs3.Count*8 // offset size
	if ofs3.Type == 2 {
		offsetDataSize += ofs3.Count * 4 // nameOffset
	}
	// 对齐
	offsetDataSize = utils.AlignUp(offsetDataSize, ofs3.Padding)
	// subFile offset
	offset := offsetDataSize
	// subFile数据
	fileData := bytes.NewBuffer(nil)

	nameMap := make(map[string]int, ofs3.Count)
	// step 1 读取Files的数据，并修改Offset和Size
	for i, file := range ofs3.Files {

		if _, has := nameMap[file.Name]; has {
			// 指向同一个文件
			continue
		} else {
			// 不同文件
			nameMap[file.Name] = i

		}
		if file.OFS3 != nil {
			// 递归写入ofs3数据
			subData := srcData[file.Offset : file.Offset+file.Size]
			subFileData = file.OFS3.createOFS3(subData)
		} else {
			subFileData, err = os.ReadFile(file.FilePath)
			if err != nil {
				if ShowLog {
					fmt.Printf("文件不存在，将使用原数据 %v\n", file.FilePath)
				}
				// 截取原数据
				subFileData = srcData[file.Offset : file.Offset+file.Size]
			}
		}
		fileData.Write(subFileData)
		file.Offset = offset
		if ofs3.SubType == 1 {
			file.Size = 0
		} else {
			file.Size = len(subFileData)
			if file.Size == 0 {
				file.Offset = 0
			}
		}
		offset += len(subFileData)
		save := offset
		offset = utils.AlignUp(offset, ofs3.Padding)
		// 对齐填充
		for i := 0; i < offset-save; i++ {
			fileData.WriteByte(0)
		}
	}

	// step 2 计算NameOffset并写入name
	if ofs3.Type == 2 {
		for _, file := range ofs3.Files {
			file.NameOffset = offset
			n, _ := fileData.WriteString(file.Name)
			fileData.WriteByte(0)
			offset += n + 1
		}
		save := offset
		offset = utils.AlignUp(offset, ofs3.Padding)
		// 对齐填充
		for i := 0; i < offset-save; i++ {
			fileData.WriteByte(0)
		}
	}

	// Offset和Size数据，和Count有关
	offsetData := make([]byte, offsetDataSize-4)
	// offsetData offset
	offset = 0
	// step 3 写入offsetData数据
	for _, file := range ofs3.Files {
		copy(offsetData[offset:offset+4], utils.WriteUInt32(file.Offset))
		offset += 4
		copy(offsetData[offset:offset+4], utils.WriteUInt32(file.Size))
		offset += 4
		if ofs3.Type == 2 {
			copy(offsetData[offset:offset+4], utils.WriteUInt32(file.NameOffset))
			offset += 4
		}
	}
	buf := bytes.NewBuffer(headerData)
	buf.Write(offsetData)
	buf.Write(fileData.Bytes())
	return buf.Bytes()
}
