package main

import (
	"flag"
	"fmt"
	"github.com/go-restruct/restruct"
)

var ShowLog = true

func main() {
	fmt.Println("WeTor wetorx@qq.com")
	var idxFile, imgFile, inputDir, output string
	var inputAppend, log bool
	flag.StringVar(&idxFile, "idx", "", "[必要]cdimg.idx文件名")
	flag.StringVar(&imgFile, "img", "", "[必要]cdimg.idx文件名")
	flag.StringVar(&inputDir, "i", "", "[打包.必要]输入文件夹路径")
	flag.StringVar(&output, "o", "", "[解包]输出文件夹路径；[打包]输出文件名")
	flag.BoolVar(&inputAppend, "append", false, "[打包]追加写入模式，待测试")
	flag.BoolVar(&log, "log", true, "显示日志")
	flag.Parse()
	restruct.EnableExprBeta()

	if len(idxFile)*len(imgFile) == 0 {
		fmt.Println("需要-idx和-img")
		return
	}
	ShowLog = log
	dfi := LoadIdx(idxFile)
	if len(inputDir) > 0 {
		// 打包
		if len(output) == 0 {
			fmt.Println("需要-o")
			return
		}
		dfi.SetDir(inputDir, true)
		dfi.ReBuildImg(imgFile, output+".img", inputAppend)
		dfi.SaveIdx(output + ".idx")
	} else {
		// 解包
		if len(output) == 0 {
			fmt.Println("需要-o")
			return
		}
		dfi.SetDir(output, false)
		dfi.LoadImg(imgFile)
	}

}
