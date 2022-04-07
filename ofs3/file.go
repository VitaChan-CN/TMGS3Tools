package ofs3

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
	"path"
	"strings"
)

func Decode(data []byte, filename string) ([]byte, string) {
	var result []byte
	switch binary.BigEndian.Uint32(data[0:4]) {
	case 0x1F8B0808:
		gz, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return data, filename
		}
		result, err = io.ReadAll(gz)
		if err != nil {
			return data, filename
		}
		filename += ".dgz"
	default:
		result = data
	}
	return result, filename
}
func Encode(data []byte, filename string) ([]byte, string) {

	var result []byte
	ext := path.Ext(filename)
	switch ext {
	case ".dgz":
		b := bytes.NewBuffer(data)
		gz := gzip.NewWriter(b)
		_, err := gz.Write(data)
		if err != nil {
			return data, filename
		}
		err = gz.Flush()
		if err != nil {
			return data, filename
		}
		result = b.Bytes()
		filename = strings.TrimSuffix(filename, ext)
	default:
		result = data
	}
	return result, filename
}
