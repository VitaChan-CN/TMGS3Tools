package main

import (
	"TMGS3Tools/ofs3"
	"TMGS3Tools/utils"
	"flag"
	"fmt"
	"github.com/go-restruct/restruct"
	"os"
)

var ShowLog = true

func main() {
	fmt.Println("WeTor wetorx@qq.com")
	fmt.Println("Version: 0.4.1")
	var idxFile, imgFile, inputDir, output, output2 string
	var ofs3File, installFile string
	var inputAppend, log, ofs3log, ofs3Mode, gz bool
	var idxFile2, imgFile2, installFile2 string
	flag.StringVar(&idxFile, "idx", "", "[DFI必要]cdimg.idx文件名")
	flag.StringVar(&imgFile, "img", "", "[DFI必要]cdimg.idx文件名")
	flag.StringVar(&installFile, "install", "", "[DFI]解密后的INSTALL.DAT文件，启用此项，将会无法使用append和patch")
	flag.StringVar(&idxFile2, "idx2", "", "[对比]cdimg.idx文件名")
	flag.StringVar(&imgFile2, "img2", "", "[对比]cdimg.idx文件名")
	flag.StringVar(&installFile2, "install2", "", "[对比]解密后的INSTALL.DAT文件")
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

	if len(idxFile2)*len(imgFile2) != 0 && len(idxFile)*len(imgFile) != 0 && len(installFile)*len(installFile2) != 0 {
		CompareDFI(idxFile, idxFile2, imgFile, imgFile2, installFile, installFile2)
		return
	}

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
		dfi.ReBuildImg(imgFile, output, installFile, inputAppend, patchOffset)
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
		dfi.LoadImg(imgFile, installFile, ofs3Mode, gz)
	}

}

func CompareDFI(idx, idx2, img, img2, inst, inst2 string) {
	dfi1 := LoadIdx(idx)
	dfi2 := LoadIdx(idx2)
	dfi1.SetDir("/", true)
	dfi2.SetDir("/", true)

	fimg1, _ := os.Open(img)
	defer fimg1.Close()
	fimg2, _ := os.Open(img2)
	defer fimg2.Close()

	finst1, _ := os.Open(inst)
	defer finst1.Close()
	finst2, _ := os.Open(inst2)
	defer finst2.Close()

	if dfi1.Count != dfi2.Count {
		fmt.Printf("子文件数量不一致！1:%d\t2:%d\n", dfi1.Count, dfi2.Count)
		return
	}
	var data1, data2 []byte
	fmt.Printf("序号\t%-9s\t%-32s\t%-32s\t文件名\n", "描述", "file1", "file2")
	for i := 0; i < dfi1.Count; i++ {
		if dfi1.Nodes[i].FilePath != dfi2.Nodes[i].FilePath {
			fmt.Printf("%d\t%v\t%v\t文件名不一致\n", i, dfi1.Nodes[i].FilePath, dfi2.Nodes[i].FilePath)
		}
		if dfi1.Nodes[i].Length != dfi2.Nodes[i].Length {
			fmt.Printf("%d\t文件大小不一致\t% 32d\t% 32d\t%v\n", i, dfi1.Nodes[i].Length, dfi2.Nodes[i].Length, dfi1.Nodes[i].FilePath)
		} else {
			data1 = make([]byte, dfi1.Nodes[i].Length)
			if dfi1.Nodes[i].Offset < dfi1.ImgSize {
				fimg1.ReadAt(data1, int64(dfi1.Nodes[i].Offset))
			} else {
				finst1.ReadAt(data1, int64(dfi1.Nodes[i].Offset-dfi1.ImgSize))
			}

			data2 = make([]byte, dfi2.Nodes[i].Length)
			if dfi2.Nodes[i].Offset < dfi2.ImgSize {
				fimg2.ReadAt(data2, int64(dfi2.Nodes[i].Offset))
			} else {
				finst2.ReadAt(data2, int64(dfi2.Nodes[i].Offset-dfi2.ImgSize))
			}

			md51 := utils.MD5(data1)
			md52 := utils.MD5(data2)
			if md51 != md52 {
				fmt.Printf("%d\t文件Hash不一致\t%v\t%v\t%v\n", i, md51, md52, dfi1.Nodes[i].FilePath)
			}
		}
	}
}
