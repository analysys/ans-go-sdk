package lib

import (
	"fmt"
	"os"
	"strings"
)

type BatchLogCollector struct {
	GerFolder  string
	GerRule    string
	PostNumber int
	postArray  []string //上报数组，从缓存数组取 ，上报之后清空
	catchArray []string //缓存数组
}

func InitBatchLogCollector(gerFolder string, gerRule string, postNumber int) *BatchLogCollector {
	return &BatchLogCollector{
		GerFolder:  gerFolder,
		GerRule:    gerRule,
		PostNumber: postNumber,
	}
}

func (log *BatchLogCollector) UpLoad(Data string, DebugMode int) bool {
	// 数据之后推进一个 数组，满足一定的条数上传（postNumber）立即上传
	log.catchArray = append(log.postArray, Data)
	// catchArray 赋值给 postArray
	log.postArray = log.catchArray
	if len(log.postArray) >= log.PostNumber {
		// 数据发送
		return log.Send(DebugMode)
		// return true
	}
	return true
}

func (log *BatchLogCollector) Send(DebugMode int) bool {
	// 防止 调用flush 进行空上报
	if len(log.postArray) == 0 {
		return true
	}
	// 校验和创建文件夹
	createFile(log.GerFolder)
	log.catchArray = []string{}
	// 根据 gerRule 生成文件名称 按 天还是按小时，
	timeStr := getTime()
	fileName := CheckURLLast(log.GerFolder) + "datas_" + timeStr["year"] + timeStr["month"] + timeStr["day"]
	if strings.ToUpper(log.GerRule) == "D" {
		fileName += ".log"
	} else {
		fileName += timeStr["hour"] + ".log"
	}
	// 要追加的字符串 gerData
	gerData := strings.Join(log.postArray, "\n")
	str := []byte(gerData + "\n")
	// 以追加模式打开文件，当文件不存在时生成文件
	txt, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	log.postArray = []string{}
	defer txt.Close()
	if err != nil {
		ErrorLog("", "", "", 600027)
		if DebugMode != 1 && DebugMode != 2 {
			fmt.Println("Save message to file :" + fileName)
			fmt.Println("Data :" + gerData)
			fmt.Println("Save message failed")
		}
		return false
	}
	// 写入文件
	n, err := txt.Write(str)
	// 当 n != len(b) 时，返回非零错误
	if err == nil && n != len(str) {
		ErrorLog("", "", "", 600027)
		if DebugMode != 1 && DebugMode != 2 {
			fmt.Println("Save message to file :" + fileName)
			fmt.Println("Data :" + gerData)
			fmt.Println("Save message failed")
		}
		return false
	}
	// 写入成功的打印日志
	SuccessLog("", fileName, gerData, 20014)
	SuccessLog("", "", "", 20015)
	return true
}
