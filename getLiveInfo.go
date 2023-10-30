/*
IDE:GoLand
PackageName:main
FileName:getLiveInfo.go
UserName:QH
CreateDate:2023/10/24
*/

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var ()

func main() {
	queryDate()
}

var sid = 18266146

func getLiveData(sTime string, eTime string) {
	// 发起GET请求
	url := "https://www.eeo.cn/test/webcast_tongji_time.php?startTime=%s&endTime=%s&sid=%d"
	targetUrl := fmt.Sprintf(url, sTime, eTime, sid)
	resp, err := http.Get(targetUrl)
	if err != nil {
		fmt.Println("GET请求发送失败：", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败：", err)
		return
	}
	fileName := fmt.Sprintf("直播回放观看数据_%d.txt", sid)
	writeInfo(fileName, strings.ReplaceAll(string(body), "<br>", "\n"))
	fmt.Println(string(body))
}

func queryDate() {

	// 定义开始日期和结束日期
	startDate := time.Date(2023, time.March, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2023, time.October, 31, 0, 0, 0, 0, time.Local)

	// 计算月份差
	monthsDiff := 0

	// 遍历月份
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		// 输出月份
		fmt.Println(currentDate.Format("2006-01"))
		cYear := currentDate.Year()
		cMonth := currentDate.Month()
		cFirstDay, cLastDay := getFirstAndLastDayOfMonth(cYear, cMonth)
		FirstDay := fmt.Sprintf("%s", cFirstDay.Format("2006-01-02"))
		LastDay := fmt.Sprintf("%s", cLastDay.Format("2006-01-02"))
		fmt.Println(FirstDay, LastDay)
		getLiveData(FirstDay, LastDay)
		//fmt.Println("本月月初:", cFirstDay.Format("2006-01-02"))
		//fmt.Println("本月月尾:", cLastDay.Format("2006-01-02"))

		// 递增月份
		currentDate = currentDate.AddDate(0, 1, 0)
		// 计算月份差
		monthsDiff++
	}
}

// 返回指定日期的月份的函数
func getMonth(date time.Time) time.Month {
	return date.Month()
}

// 获取指定月份的第一天和最后一天
func getFirstAndLastDayOfMonth(year int, month time.Month) (time.Time, time.Time) {
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := firstDay.AddDate(0, 1, -1)
	// 输出本月月初和月尾日期
	//fmt.Println("本月月初:", firstDay.Format("2006-01-02"))
	//fmt.Println("本月月尾:", lastDay.Format("2006-01-02"))
	return firstDay, lastDay
}
func writeInfo(fileName string, msg string) {
	data := msg + "\n"
	// 打开文件，如果文件不存在则创建
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 将数据写入文件
	_, err = file.Write([]byte(data))
	if err != nil {
		fmt.Println("无法写入文件:", err)
		return
	}

	fmt.Println("数据已写入文件:", fileName)
}
