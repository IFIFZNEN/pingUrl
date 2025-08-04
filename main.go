package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func ping(url string, respCh chan string, errCh chan error) {
	resp, err := http.Get(url)
	if err != nil {
		errCh <- fmt.Errorf("%v - невалидный", url)
		return
	}
	respCh <- fmt.Sprintf("%v - %d", url, resp.StatusCode)
}

func main() {
	path := flag.String("file", "url.txt", "path to URL file")
	flag.Parse()
	file, err := os.ReadFile(*path)
	if err != nil {
		panic(err.Error())
	}
	urlSlice := strings.Split(string(file), "\n")
	respChStr := make(chan string)
	errCh := make(chan error)
	for _, url := range urlSlice {
		go ping(url, respChStr, errCh)
	}
	for range urlSlice {
		select {
		case errResp := <-errCh:
			fmt.Println(errResp)
		case resp := <-respChStr:
			fmt.Println(resp)
		}
	}
}
