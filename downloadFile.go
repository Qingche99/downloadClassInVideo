/*
IDE:GoLand
PackageName:main
FileName:downloadFile.go
UserName:QH
CreateDate:2023/10/26
*/
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	fileName = "ClassIn平台课节视频汇总_20231027185611_1021998.csv"
	dList    []downInfo
)

func main() {
	openCsv(fileName)
	multiThreadedDownload(dList)

}

type downInfo struct {
	filePath    string
	downloadUrl string
	fileType    string
	startPos    int64
}

type ProgressWriter struct {
	Total       int64
	Progress    int64
	DownloadUrl string
	mu          sync.Mutex
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.mu.Lock()
	pw.Progress += int64(n)
	progress := float64(pw.Progress) / float64(pw.Total) * 100
	pw.mu.Unlock()

	fmt.Printf("Downloaded: %s / %s (%.2f%%) - %s\r", bytesToString(pw.Progress), bytesToString(pw.Total), progress, pw.DownloadUrl)

	return n, nil
}

// 字节转字符串
func bytesToString(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// 递归创建文件夹
func mkDir(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Println("创建文件夹失败:", err)
		return
	}
	//fmt.Println("文件夹创建成功")
}

// 字符串删除\n\t
func stringRmNT(strContent string) string {

	retStrN := strings.ReplaceAll(strContent, "\n", "")
	retStrT := strings.ReplaceAll(retStrN, "\t", "")
	return retStrT
}

// 打开csv文件
func openCsv(filePath string) {
	// 打开CSV文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	// 创建一个CSV reader
	reader := csv.NewReader(file)

	// 读取所有行
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("读取CSV文件失败:", err)
		return
	}
	readRows(rows)
}

// 读取二维切片
func readRows(rows [][]string) {

	keys := []string{"classId", "courseId", "courseName", "className", "classStart", "classEnd", "teacher", "videoTime", "videoId", "url"} // keys
	fileRot := fileName[0 : len(fileName)-4]
	lastId := ""
	n := 1
	// 遍历每一行并输出
	for _, eachRow := range rows {
		eachDownMap := make(map[string]string)
		if !strings.HasPrefix(eachRow[9], "http") {
			continue
		}
		// 遍历键切片
		for i := 0; i < len(keys); i++ {
			// 检查值切片是否越界
			if i < len(eachRow) {
				eachDownMap[keys[i]] = eachRow[i]
			} else {
				eachDownMap[keys[i]] = ""
			}

		}

		if lastId == eachDownMap["classId"] {
			//fmt.Println(row[0], downMap["classId"])
			n++
		} else {
			lastId = eachDownMap["classId"]
			n = 1
		}
		//fmt.Println(eachRow[0], n)

		classId := eachDownMap["classId"]
		courseName := stringRmNT(eachDownMap["courseName"])
		className := eachDownMap["className"]
		url := eachDownMap["url"]
		mkPath := fmt.Sprintf("./%s/%s", fileRot, courseName) // 拼接课程文件夹路径
		mkDir(mkPath)
		dirPath := fmt.Sprintf("%s/%s-%s-%d", mkPath, classId, className, n)

		//fmt.Println(dirPath, url)
		var info downInfo

		info.downloadUrl = url
		info.filePath = stringRmNT(dirPath)
		info.fileType = ".mp4"

		dList = append(dList, info)

	}

}

// 多线程下载方法
func multiThreadedDownload(downloadSlice []downInfo) {
	var wg sync.WaitGroup
	wg.Add(len(downloadSlice))
	for _, v := range downloadSlice {
		go func(info downInfo) {
			defer wg.Done()

			resp, err := http.Get(info.downloadUrl)
			if err != nil {
				fmt.Printf("下载文件失败 %s: %v\n", info.downloadUrl, err)
				return
			}
			defer resp.Body.Close()

			downloadFilePath := info.filePath + info.fileType

			if _, err := os.Stat(downloadFilePath); os.IsNotExist(err) { // 检查文件是否存在

				out, err := os.Create(downloadFilePath)
				if err != nil {
					fmt.Printf("创建文件失败%s: %v\n", downloadFilePath, err)
					return
				}
				defer out.Close()

				progressWriter := &ProgressWriter{
					Total:       resp.ContentLength,
					Progress:    0,
					DownloadUrl: info.downloadUrl,
				}

				_, err = io.Copy(out, io.TeeReader(resp.Body, progressWriter))
				if err != nil {
					fmt.Printf("写入文件失败 | %s | %v\n", downloadFilePath, err)
				} else {
					fmt.Printf("下载成功 | %s | %s\n", downloadFilePath, info.downloadUrl)
				}
			} else {
				if _, err := os.Stat(downloadFilePath); err == nil {
					fileInfo, err := os.Stat(downloadFilePath)
					if err != nil {
						fmt.Printf("获取文件信息失败 %s: %v\n", downloadFilePath, err)
						return
					}

					info.startPos = fileInfo.Size()
					out, err := os.OpenFile(downloadFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
					if err != nil {
						fmt.Printf("创建文件失败 %s: %v\n", downloadFilePath, err)
						return
					}
					defer out.Close()

					req, err := http.NewRequest("GET", info.downloadUrl, nil)
					if err != nil {
						fmt.Printf("创建请求失败 %s: %v\n", info.downloadUrl, err)
						return
					}

					req.Header.Set("Range", fmt.Sprintf("bytes=%d-", info.startPos))

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						fmt.Printf("请求文件失败 %s: %v\n", info.downloadUrl, err)
						return
					}
					defer resp.Body.Close()

					progressWriter := &ProgressWriter{
						Total:       info.startPos + resp.ContentLength,
						Progress:    info.startPos,
						DownloadUrl: info.downloadUrl,
					}

					_, err = io.Copy(out, io.TeeReader(resp.Body, progressWriter))
					if err != nil {
						fmt.Printf("写入文件失败 | %s | %v\n", downloadFilePath, err)
					} else {
						fmt.Printf("下载成功 | %s \n", downloadFilePath)
					}
				}
				//fmt.Printf("文件已存在，忽略 | %s | %s\n", downloadFilePath, info.downloadUrl)
			}

		}(v)
	}

	wg.Wait()
	fmt.Println("\n全部文件下载完成。")
}
