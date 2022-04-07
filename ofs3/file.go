package ofs3

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"path"
	"strings"
)

func DecodeName(data []byte, filename string) string {
	switch binary.BigEndian.Uint32(data[0:4]) {
	case 0x1F8B0808:
		return filename + ".dgz"
	default:
		return filename
	}
}

func Decode(data []byte) []byte {
	var result []byte
	switch binary.BigEndian.Uint32(data[0:4]) {
	case 0x1F8B0808:
		b := bytes.NewBuffer(data)
		gz, err := gzip.NewReader(b)
		fmt.Println(gz.Header)

		if err != nil {
			fmt.Println(err)
			return data
		}
		result, err = io.ReadAll(gz)
		if err != nil {
			fmt.Println(err)
			return data
		}
		fmt.Println(2222, len(result))
	default:
		result = data
	}
	return result
}
func Encode(data []byte, filename string) ([]byte, string) {

	ext := path.Ext(filename)
	fmt.Println(filename, ext)
	switch ext {
	case ".dgz":
		b := bytes.NewBuffer(nil)
		gz := gzip.NewWriter(b)
		gz.Header.OS = 11
		_, err := gz.Write(data)
		if err != nil {
			return data, filename
		}
		err = gz.Flush()
		if err != nil {
			return data, filename
		}
		filename = strings.TrimSuffix(filename, ext)
		return b.Bytes(), filename
	default:
		return data, filename
	}

}
