package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func AlignUp(n, align int) int {
	/* 当n为align的整数倍 返回n */
	if n/align*align == n {
		return n
	}
	/* 当n不是align的整数倍 返回>n,且离n最近的align的倍数 */
	return (n/align + 1) * align
}
func MD5(d []byte) string {
	r := md5.Sum(d)
	return hex.EncodeToString(r[:])
}
func MD5F(fName string) string {
	f, e := os.Open(fName)
	if e != nil {
		log.Fatal(e)
	}
	h := md5.New()
	_, e = io.Copy(h, f)
	if e != nil {
		log.Fatal(e)
	}
	return hex.EncodeToString(h.Sum(nil))
}
