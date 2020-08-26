package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HTTPGet 得到 url 对应的内容。
func HTTPGet(url string) (result string, err error) {
	resp, err1 := http.Get(url) //发送get请求
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	// 读取网页内容
	buf := make([]byte, 1024*4)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		// 累加读取的内容
		result += string(buf[:n])
	}
	return
}

func main() {
	var serverKey = ""
	eventsTick := time.NewTicker(time.Duration(1) * time.Second)
	defer eventsTick.Stop()
	for {
		select {
		case <-eventsTick.C:
			var movies = [...]string{"釜山行"}
			var urls = [...]string{"https://www.dy2018.com"}

			for _, movie := range movies {
				for _, url := range urls {
					result, err := HTTPGet(url)

					if err != nil {
						return
					}

					var inde = strings.Index(result, movie)
					if inde != -1 {
						fmt.Println("没有" + movie)
					} else {
						fmt.Println("有" + movie)
						http.Get("https://sc.ftqq.com/" + serverKey + ".send?text=" + movie + url)
					}
				}
			}
		}
	}
}
