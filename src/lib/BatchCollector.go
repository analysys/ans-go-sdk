package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type BatchCollector struct {
	UploadURL  string
	PostNumber int
	postArray  []string //上报数组，从缓存数组取 ，上报之后清空
	catchArray []string //缓存数组
}

func InitBatchCollector(uploadURL string, postNumber int) *BatchCollector {
	return &BatchCollector{
		UploadURL:  uploadURL,
		PostNumber: postNumber,
	}
}

func (up *BatchCollector) UpLoad(Data string, DebugMode int) bool {
	// 数据之后推进一个 数组，满足一定的条数上传（postNumber）立即上传
	up.catchArray = append(up.catchArray, Data)
	// catchArray 赋值给 postArray
	up.postArray = up.catchArray
	if len(up.postArray) >= up.PostNumber {
		// 数据发送
		return up.Send(DebugMode)
		// return true
	}
	return true
}

func (up *BatchCollector) Send(DebugMode int) bool {
	// 防止 调用flush 进行空上报
	if len(up.postArray) == 0 {
		return true
	}
	// 上报的相关方法
	up.catchArray = []string{}
	postData := "[" + strings.Join(up.postArray, ",") + "]"
	// 校验 uploadURL
	if !CheckURL(up.UploadURL) {
		return false
	}
	uploadURL := CheckURLLast(up.UploadURL) + "up"
	if uploadURL == "up" {
		return false
	}
	SuccessLog("", up.UploadURL, postData, 20012)
	resp, err := http.Post(uploadURL, "application/x-www-form-urlencoded", strings.NewReader(postData))
	up.postArray = []string{}
	if err != nil {
		ErrorLog("", "", "", 60001)
		if DebugMode != 1 && DebugMode != 2 {
			fmt.Println("Send message to server :" + up.UploadURL)
			fmt.Println("Data :" + postData)
			fmt.Println("Send message failed")
		}
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		ErrorLog("", "", "", 60001)
		if DebugMode != 1 && DebugMode != 2 {
			fmt.Println("Send message to server :" + up.UploadURL)
			fmt.Println("Data :" + postData)
			fmt.Println("Send message failed")
		}
		return false
	}
	// code 返回 200 的操作
	if string(body) == "{\"code\":200}" || string(body) == "H4sIAAAAAAAAAKtWSs5PSVWyMjIwqAUAVAOW6gwAAAA=" {
		SuccessLog("", "", "", 20001)
		return true
	}
	ErrorLog("", string(body), "", 60001)
	if DebugMode != 1 && DebugMode != 2 {
		fmt.Println("Send message to server :" + up.UploadURL)
		fmt.Println("Data :" + postData)
		fmt.Println("Send message failed," + string(body))
	}
	return false
}
