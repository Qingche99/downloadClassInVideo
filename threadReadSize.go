/*
IDE:GoLand
PackageName:main
FileName:threadReadSize.go
UserName:QH
CreateDate:2023/10/24
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func getURLHead(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Head(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	fmt.Printf("URL: %s\n", url)
	for header, values := range resp.Header {
		fmt.Printf("%s: %s\n", header, values)
	}
	fmt.Println()
}

func main() {
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
	}

	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go getURLHead(url, &wg)
	}

	wg.Wait()
}
