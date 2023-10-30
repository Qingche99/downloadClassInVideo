/*
IDE:GoLand
PackageName:ClassIn_edu
FileName:01sayHello.go
UserName:QH
CreateDate:2023/10/20
*/

package main

import (
	"fmt"
)

var ()

func main() {
	fmt.Println("sayHello")
	// %d 表示整型数字，%s 表示字符串
	var stockCode = 123
	var endDate = "2020-12-31"
	var url = "Code=%d&endDate=%s"
	var targetUrl = fmt.Sprintf(url, stockCode, endDate)
	fmt.Println(targetUrl)
	S1 := 123
	fmt.Println(S1)

}
