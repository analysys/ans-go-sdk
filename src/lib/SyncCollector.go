package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Collector interface {
	UpLoad(Data string, DebugMode int) bool
	Send(DebugMode int) bool
}

func InitSyncCollector(uploadURL string) *SyncCollector {
	return &SyncCollector{
		UploadURL: uploadURL,
	}
}

type SyncCollector struct {
	UploadURL  string
	postArray  []string //上报数组，从缓存数组取 ，上报之后清空
	catchArray []string //缓存数组
}

func (up *SyncCollector) UpLoad(Data string, DebugMode int) bool {
	// 设置了上传时间 > = 0 设置上传条数大于0 ,会进入时间条件的上传，到了时间间隔就上传
	up.catchArray = append(up.postArray, Data)
	// catchArray 赋值给 postArray
	up.postArray = up.catchArray
	if len(up.postArray) > 0 {
		// 数据发送
		return up.Send(DebugMode)
		// return true
	}
	return true
}

func (up *SyncCollector) Send(DebugMode int) bool {
	// 防止 调用flush 进行空上报
	if len(up.postArray) == 0 {
		return true
	}
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
	return true
}
