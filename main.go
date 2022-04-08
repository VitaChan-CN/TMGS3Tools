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
	var idxFile, imgFile, inputDir, output, output2 string
	var ofs3File string
	var inputAppend, log, ofs3log, ofs3Mode, gz bool
	flag.StringVar(&idxFile, "idx", "", "[DFI必要]cdimg.idx文件名")
	flag.StringVar(&imgFile, "img", "", "[DFI必要]cdimg.idx文件名")
	flag.StringVar(&ofs3File, "ofs3", "", "[OFS3必要] OFS3文件名")
	flag.StringVar(&inputDir, "i", "", "[打包]输入文件夹路径")
	flag.StringVar(&output, "o", "", "[解包]输出文件夹路径, [打包]输出文件名")
	flag.StringVar(&output2, "o2", "", "[打包]输出idx文件名，若为空则为-o后增加.idx")
	flag.BoolVar(&inputAppend, "append", false, "[打包]追加写入模式")
	flag.BoolVar(&log, "log", false, "显示日志")
	flag.BoolVar(&ofs3Mode, "dfi.ofs3", false, "[DFI解包]递归解包所有OFS3格式文件")
	flag.BoolVar(&ofs3log, "ofs3.log", false, "显示OFS3日志")
	flag.BoolVar(&gz, "gz", false, "解包时是否自动解压gz文件(解压后为.dgz文件，导入需要手动压缩并去掉后缀)")

	var patchOffset int
	flag.IntVar(&patchOffset, "patch", 0, "[打包]对已存在的-o文件的指定位置进行修改而不是创建新的，仅append模式有效。输入原img大小")
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
		dfi.ReBuildImg(imgFile, output, inputAppend, patchOffset)
		idxName := output + ".idx"
		if len(output2) > 0 {
			idxName = output2
		}
		dfi.SaveIdx(idxName)
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
