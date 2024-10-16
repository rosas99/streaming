package main

import (
	"fmt"
	"github.com/rosas99/streaming/internal/pkg/client"
	"github.com/rosas99/streaming/pkg/log"
)

type AuthSuccess struct {
}

func main() {
	url := "http://example.com/api"
	request := client.NewRequest()
	response, err := request.
		SetBody([]byte(`{"username":"testuser", "password":"testpass"}`)).
		SetResult(&AuthSuccess{}).
		Post(url)
	if err != nil {
		log.Errorf("请求失败: %v\n", err)
	}

	fmt.Print(response)

	if response.StatusCode() >= 400 {
		fmt.Printf("服务器返回错误状态码: %d\n", response.StatusCode())
	}

	//resp, err := client.R().
	//	SetBody(Article{
	//		Tags: []string{"new tag1", "new tag2"},
	//	}).
	//	SetAuthToken("C6A79608-782F-4ED0-A11D-BD82FAD829CD").
	//	SetError(&errorStruct). // 注意这里传入了错误结构体的地址
	//	Patch("https://myapp.com/articles/1234")
	//
	//if err != nil {
	//	fmt.Printf("请求失败: %v\n", err)
	//} else if resp.IsError() { // 检查是否是错误响应
	//	fmt.Printf("错误码: %d\n", errorStruct.Code)
	//	// 根据需要处理错误码
	//}

}
