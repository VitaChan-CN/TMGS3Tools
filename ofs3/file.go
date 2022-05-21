package ofs3

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
)

type FileType uint8

const (
	FILE_OTHER FileType = iota
	FILE_OFS3
	FILE_GZ
)

func GetFileType(data []byte) FileType {
	v1 := binary.BigEndian.Uint16(data[0:2])
	v2 := binary.BigEndian.Uint16(data[2:4])
	if v1 == 0x1F8B {
		return FILE_GZ
	} else if v1 == 0x4F46 && v2 == 0x5333 { //OFS3
		return FILE_OFS3
	} else {
		return FILE_OTHER
	}
}

func GzDecode(data []byte, onlyHeader bool) ([]byte, *gzip.Header) {
	var result []byte
	var header *gzip.Header
	b := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(b)
	if err != nil {
		fmt.Println(err)
		return data, nil
	}
	defer gz.Close()
	header = &gz.Header
	if onlyHeader {
		return nil, header
	}

	result, err = io.ReadAll(gz)
	if err != nil {
		fmt.Println(err)
		return data, nil
	}
	return result, header
}
func GzEncode(data []byte, header *gzip.Header) []byte {
	b := bytes.NewBuffer(nil)
	gz, err := gzip.NewWriterLevel(b, 9)
	if err != nil {
		fmt.Println(err)
	}
	gz.Header = *header
	_, err = gz.Write(data)
	if err != nil {
		fmt.Println(err)
		return data
	}
	err = gz.Close()
	if err != nil {
		fmt.Println(err)
		return data
	}

	return b.Bytes()
}
