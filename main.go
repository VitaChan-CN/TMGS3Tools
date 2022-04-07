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
	var ofs3File string
	var inputAppend, log, ofs3log, ofs3Mode, gz bool
	flag.StringVar(&idxFile, "idx", "", "[DFI必要]cdimg.idx文件名")
	flag.StringVar(&imgFile, "img", "", "[DFI必要]cdimg.idx文件名")
	flag.StringVar(&ofs3File, "ofs3", "", "[OFS3必要] OFS3文件名")
	flag.StringVar(&inputDir, "i", "", "[打包]输入文件夹路径")
	flag.StringVar(&output, "o", "", "[解包]输出文件夹路径, [打包]输出文件名")
	flag.BoolVar(&inputAppend, "append", false, "[打包]追加写入模式，待测试")
	flag.BoolVar(&log, "log", false, "显示日志")
	flag.BoolVar(&ofs3Mode, "dfi.ofs3", false, "[DFI解包]递归解包所有OFS3格式文件")
	flag.BoolVar(&ofs3log, "ofs3.log", false, "显示OFS3日志")
	flag.BoolVar(&gz, "gz", true, "解包打包是否自动解压、压缩gz文件(.dgz)")

	flag.Parse()
	restruct.EnableExprBeta()

	ShowLog = log
	ofs3.ShowLog = ofs3log
	if len(ofs3File) > 0 && len(output) > 0 && len(inputDir) > 0 {
		// ofs3打包
		data, _ := os.ReadFile(ofs3File)
		ofs := ofs3.OpenOFS3(data, inputDir)
		ofs.ReBuild(data, output, gz)
		return
	} else if len(ofs3File) > 0 && len(output) > 0 {
		// ofs3解包
		data, _ := os.ReadFile(ofs3File)
		ofs := ofs3.OpenOFS3(data, output)
		ofs.WriteFile(data, output, gz)
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
		dfi.LoadImg(imgFile, ofs3Mode, gz)
	}

}
