package main

import (
	"TMGS3Tools/ofs3"
	"flag"
	"fmt"
	"github.com/go-restruct/restruct"
	"os"
)

var ShowLog = true

func main() {
	fmt.Println("WeTor wetorx@qq.com")
	var idxFile, imgFile, inputDir, output string
	var inputAppend, log, ofs3log, ofs3Mode bool
	flag.StringVar(&idxFile, "idx", "", "[必要]cdimg.idx文件名")
	flag.StringVar(&imgFile, "img", "", "[必要]cdimg.idx文件名")
	flag.StringVar(&inputDir, "i", "", "[打包.必要]输入文件夹路径；[OFS3解包]输入OFS3文件名")
	flag.StringVar(&output, "o", "", "[解包]输出文件夹路径；[打包]输出文件名")
	flag.BoolVar(&inputAppend, "append", false, "[打包]追加写入模式，待测试")
	flag.BoolVar(&log, "log", true, "显示日志")
	flag.BoolVar(&ofs3log, "ofs3.log", false, "显示OFS3日志")
	flag.BoolVar(&ofs3Mode, "ofs3", false, "[解包]递归解包所有OFS3格式文件，无法打包")
	flag.Parse()
	restruct.EnableExprBeta()

	ShowLog = log
	ofs3.ShowLog = ofs3log

	if ofs3Mode && len(inputDir) > 0 {
		if len(output) == 0 {
			fmt.Println("需要-o")
			return
		}
		data, _ := os.ReadFile(inputDir)
		ofs3.OpenOFS3(data, output)
		return
	} else if len(idxFile)*len(imgFile) == 0 {
		fmt.Println("需要-idx和-img")
		return
	}

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
		dfi.LoadImg(imgFile, ofs3Mode)
	}

}
