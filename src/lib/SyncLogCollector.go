package lib

import (
	"fmt"
	"os"
	"strings"
)

type SyncLogCollector struct {
	gerFolder  string
	gerRule    string
	postArray  []string //上报数组，从缓存数组取 ，上报之后清空
	catchArray []string //缓存数组
}

func InitSyncLogCollector(gerFolder string, gerRule string) *SyncLogCollector {
	return &SyncLogCollector{
		gerFolder: gerFolder,
		gerRule:   gerRule,
	}
}

func (log *SyncLogCollector) UpLoad(Data string, DebugMode int) bool {
	// 设置了上传时间 > = 0 设置上传条数大于0 ,会进入时间条件的上传，到了时间间隔就上传
	log.catchArray = append(log.postArray, Data)
	// catchArray 赋值给 postArray
	log.postArray = log.catchArray
	if len(log.postArray) > 0 {
		// 数据发送
		log.Send(DebugMode)
		return true
	}
	return false
}

func (log *SyncLogCollector) Send(DebugMode int) bool {
	// 防止 调用flush 进行空上报
	if len(log.postArray) == 0 {
		return false
	}
	// 校验和创建文件夹
	createFile(log.gerFolder)
	log.catchArray = []string{}
	// 根据 gerRule 生成文件名称 按 天还是按小时，
	timeStr := getTime()
	fileName := CheckURLLast(log.gerFolder) + "datas_" + timeStr["year"] + timeStr["month"] + timeStr["day"]
	if strings.ToUpper(log.gerRule) == "D" {
		fileName += ".log"
	} else {
		fileName += timeStr["hour"] + ".log"
	}
	// 要追加的字符串 gerData
	gerData := strings.Join(log.postArray, "")
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
