package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/axgle/mahonia"
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

	result = ConvertToString(result, "gbk", "utf-8")
	return
}

// 编码
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

var movies = list.New()
var urls = list.New()

func main() {
	http.HandleFunc("/getMovie", movieFunc)
	http.HandleFunc("/start", start)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func start(writer http.ResponseWriter, request *http.Request) {
	movies.PushFront("釜山行")
	urls.PushFront("https://www.dy2018.com")

	var serverKey = "server 酱"
	eventsTick := time.NewTicker(time.Duration(1) * time.Second)
	defer eventsTick.Stop()
	for {
		select {
		case <-eventsTick.C:
			fmt.Println(time.Now())
			for movie := movies.Front(); movie != nil; movie = movie.Next() {
				for url := urls.Front(); url != nil; url = url.Next() {
					var movieStr = movie.Value.(string)
					result, err := HTTPGet(url.Value.(string))

					if err != nil {
						return
					}
					var inde = strings.Index(result, movieStr)
					if inde == -1 {
						fmt.Println("没有" + movieStr)
					} else {
						fmt.Println("有" + movieStr)
						movies.Remove(movie)
						http.Get("https://sc.ftqq.com/" + serverKey + ".send?text=" + movieStr + url.Value.(string))
					}
				}
			}
		}
	}
}

type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func movieFunc(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	movie, uError := request.Form["movie"]

	var result Resp

	fmt.Println(uError)
	if movie != nil && uError {
		movies.PushFront(movie[0])
		result.Msg = "添加成功"
		result.Code = "200"

		if err := json.NewEncoder(writer).Encode(result); err != nil {
			log.Fatal(err)
		}
	}
}
