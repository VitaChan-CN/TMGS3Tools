package main

import (
	"flag"
	"fmt"
	"github.com/go-restruct/restruct"
)

var ShowLog = false

func main() {
	fmt.Println("WeTor wetorx@qq.com")
	var idxFile, imgFile, inputDir, output string
	var log bool
	flag.StringVar(&idxFile, "idx", "", "cdimg.idx文件名")
	flag.StringVar(&imgFile, "img", "", "cdimg.idx文件名")
	flag.StringVar(&inputDir, "i", "", "[仅打包需要]输入文件夹路径")
	flag.StringVar(&output, "o", "", "[解包]输出文件夹路径；[打包]输出文件名")
	flag.BoolVar(&log, "o", false, "显示日志")
	flag.Parse()
	restruct.EnableExprBeta()
	ShowLog = log
	dfi := LoadIdx(idxFile)
	if len(inputDir) > 0 {
		dfi.ReBuildImg(imgFile, inputDir, output)
	} else {
		dfi.LoadImg(imgFile, output)
	}

}
