/*
IDE:GoLand
PackageName:main
FileName:main.go
UserName:QH
CreateDate:2023/10/30
*/
//go:generate rsrc -ico resource/icon.ico -manifest resource/goversioninfo.exe.manifest -o main.syso
package main

import (
	"fmt"
	"github.com/gonutz/w32/v2"
	"os"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

const (
	SW_MAXIMIZE = 3
)

func main() {

	// 获取控制台窗口句柄
	handle := w32.GetConsoleWindow()

	// 设置控制台窗口为最大化
	w32.ShowWindow(handle, SW_MAXIMIZE)
	// 禁止调整控制台窗口大小
	w32.SetWindowLong(handle, w32.GWL_STYLE, w32.GetWindowLong(handle, w32.GWL_STYLE)&^w32.WS_SIZEBOX)

	// 获取当前程序所在目录的路径
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("获取当前目录失败：", err)
		return
	}

	// 打开文件对话框
	openDialog := syscall.MustLoadDLL("comdlg32.dll").MustFindProc("GetOpenFileNameW")

	var file [4096]uint16
	file[0] = 0

	//filter := "\"Csv Files (*.csv)\\000*.csv\\000\" +" + "All Files\000*.*\000"

	var ofn struct {
		structSize      uint32
		ownerHWND       uintptr
		instance        uintptr
		filter          *uint16
		customFilter    *uint16
		maxCustomFilter uint32
		selectedFilter  uint32
		file            *uint16
		maxFile         uint32
		fileTitle       *uint16
		maxFileTitle    uint32
		initialDir      *uint16
		title           *uint16
		flags           uint32
		fileOffset      uint16
		fileExtension   uint16
		defExt          *uint16
		custData        uintptr
		fnHook          uintptr
		templateName    *uint16
	}

	ofn.structSize = uint32(unsafe.Sizeof(ofn))
	ofn.file = &file[0]
	ofn.maxFile = uint32(len(file))
	//ofn.filter = syscall.StringToUTF16Ptr(filter)
	ofn.initialDir = &utf16.Encode([]rune(currentDir))[0] // 设置初始目录为当前目录

	ret, _, _ := openDialog.Call(uintptr(unsafe.Pointer(&ofn)))

	if ret != 0 {
		fileName := syscall.UTF16ToString(file[:])
		fmt.Println("选择的文件:", fileName)
		openCsv(fileName)
		multiThreadedDownload(dList)

	} else {
		fmt.Println("用户取消了选择文件")
	}
	fmt.Println("程序执行结束，请按 Enter 键关闭窗口...")
	_, err = fmt.Scanln()
	if err != nil {
		return
	}
	os.Exit(0)

}

// 定义命令行选项
//download := flag.String("download", "fileName", "下载课程视频")

// 解析命令行参数
//flag.Parse()

// 根据选项进行相应操作
//if *download != "" {
//	fmt.Println("下载课程视频")
//	fmt.Println("文件路径:", *download)
//	// 执行下载操作的代码
//	openCsv(*download)
//	multiThreadedDownload(dList)
//}

//}
