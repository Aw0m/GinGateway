package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func BenchmarkLimiter(b *testing.B) {
	for j := 0; j <= b.N; j++ {
		go sendMessage()
		time.Sleep(time.Millisecond)
	}
}

func sendMessage() {
	req, _ := http.NewRequest("GET", "http://localhost:8082/api/work/test", nil)
	// 比如说设置个token
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJyb290IiwidXNlck5hbWUiOiJyb290IiwiZXhwIjoxNjU1NjI4ODEzLCJpc3MiOiJ0ZXN0IiwibmJmIjoxNjUzMDM2ODEzfQ.06f6hpkxKpR7pLwisi5MA4zwKvSZRR3x_9DVtHkTHto")
	// 再设置个json
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	//resp, err := http.Get(serviceUrl + "/topic/query/false/lsj")
	if err != nil {
		log.Println("失败！")
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(respByte)
}
