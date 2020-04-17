package lib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func CheckURLLast(serverURL string) string {
	// 此方法使用 url 不能为空，不然报错
	lengthInt := strings.Count(serverURL, "") - 1
	if lengthInt > 0 {
		content := serverURL[len(serverURL)-1 : len(serverURL)]
		if content == "/" {
			return serverURL
		} else {
			return serverURL + "/"
		}
	}
	ErrorLog("InitAnalysysAgent", "serverURL", serverURL, 600021)
	return ""
}

// 校验 key
func CheckKey(fn string, key string) bool {
	// 字符串的长度限制，99.不能关键字，特殊字符
	length99Status, length99Code := length99(key)
	notSpecialStatus, notSpecialCode := notSpecialCharacters(key)
	keywordsStatus, keywordsCode := keywords(key)
	if !length99Status {
		ErrorLog(fn, key, "", length99Code)
		return false
	}
	if !notSpecialStatus {
		ErrorLog(fn, key, "", notSpecialCode)
		return false
	}
	if !keywordsStatus {
		ErrorLog(fn, key, "", keywordsCode)
		return false
	}
	return true
}

// 校验 value
func CheckValue(fn string, key string, value interface{}) (bool, interface{}) {
	// 数字  布尔  字符串（长度255，超过255，截取）
	switch value.(type) {
	case int:
		return true, value
	case float64:
		return true, value
	case string:
		// 要进行超长判断，和超长截取
		lengthStr := fmt.Sprint(value)
		lengthInt := strings.Count(lengthStr, "") - 1
		length255status, length255Code := length255(lengthStr)
		if !length255status {
			if lengthInt > 8192 {
				lengthStr = lengthStr[0:8191] + "$"
				ErrorLog(fn, key, fmt.Sprint(value), length255Code)
				return false, lengthStr
			}
			ErrorLog(fn, key, fmt.Sprint(value), length255Code)
			return false, value
		}
		return true, value
	case bool:
		return true, value
	case []string:
		// 这里可以抓到 slice string ，(之后强制转换)但是不能抓到 array string,所以 value 只支持 []string
		valArr, ok := value.([]string)
		if ok {
			valStr := []string{}
			if len(valArr) > 100 {
				ErrorLog(fn, key, fmt.Sprint(value), 60007)
			}
			for i := 0; i < len(valArr); i++ {
				// 超长截取
				itemSta, itemCode := length255(valArr[i])
				if !itemSta {
					lengthInt := strings.Count(valArr[i], "") - 1
					if lengthInt > 8192 {
						valReturn := valArr[i][0:8191] + "$"
						valStr = append(valStr, valReturn)
					} else {
						valStr = append(valStr, valArr[i])
					}
					// value 打印转化一下，变成有，防止打印不标准
					valueLog, _ := json.Marshal(value)
					ErrorLog(fn, key, string(valueLog), itemCode)
					continue
				}
				valStr = append(valStr, valArr[i])
			}
			return false, valStr
		}
		return false, value
	default:
		// 不满足以上的几种要求，提醒支持的类型不对，依旧注册成功
		ErrorLog(fn, key, "", 60008)
		return false, value
	}
}

// 校验 uploadTime 是否满足13位的 int
func CheckUploadTime(fn string, uploadTime int) bool {
	//Itoa方法可以把整数转换为字符串
	data := strconv.Itoa(uploadTime)
	lenthInt := strings.Count(data, "") - 1
	if lenthInt != 13 {
		ErrorLog(fn, fmt.Sprint(uploadTime), "upLoadTime", 600024)
		return false
	}
	return true
}

// 校验 distincID
func CheckDistinctID(fn string, distinctID string) bool {
	if !checkID(distinctID) {
		ErrorLog(fn, "distinctID", "", 60006)
		return false
	}
	return true
}

// 校验 aliasID
func CheckAliasID(fn string, aliasID string) bool {
	if !checkID(aliasID) {
		ErrorLog(fn, "aliasID", "", 60006)
		return false
	}
	return true
}

func checkID(id string) bool {
	idStatus, _ := keyLength255(id)
	if !idStatus {
		return false
	}
	return true
}

// 校验 eventName
func CheckEventName(eventName string) bool {
	lengthInt := strings.Count(eventName, "") - 1
	//  99长度  不包含特殊字符
	FnName := "track"
	Key := "eventName"
	if lengthInt > 0 {
		Key = eventName
		FnName = eventName
	}
	testBool1, testCode1 := length99(eventName)
	testBool2, testCode2 := notSpecialCharacters(eventName)
	if !testBool1 {
		ErrorLog(FnName, Key, "", testCode1)
		return false
	}
	if !testBool2 {
		ErrorLog(FnName, Key, "", testCode2)
		return false
	}
	return true
}

// 校验 platform
func CheckPlatform(platform string) string {
	if strings.ToUpper(platform) == "GO" || strings.ToUpper(platform) == "" {
		platform = "Go"
	}
	if strings.ToUpper(platform) == "ANDROID" {
		platform = "Android"
	}
	if strings.ToUpper(platform) == "IOS" {
		platform = "iOS"
	}
	if strings.ToUpper(platform) == "JS" {
		platform = "JS"
	}
	if strings.ToUpper(platform) == "WECHAT" {
		platform = "WeChat"
	}
	return platform
}

// 校验 appid
func CheckAppid(appid string) bool {
	if appid == "" {
		ErrorLog("InitAnalysysAgent", "appid", appid, 600021)
		return false
	}
	return true
}

func IsRequestURL(rawurl string) bool {
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return false //Couldn't even parse the rawurl
	}
	if len(url.Scheme) == 0 {
		return false //No Scheme found
	}
	return true
}

// 校验 uploadURL
func CheckURL(url string) bool {
	if !IsRequestURL(url) {
		ErrorLog("InitAnalysysAgent", "uploadURL", "", 600023)
		return false
	}
	return true
}

// 校验map
func CheckMap(fn string, property map[string]interface{}) (bool, map[string]interface{}) {
	// 重新生成一个新的 map,假如超长截取,返回符合的map
	mapFlage := true
	for kStr, Value := range property {
		CheckKey(fn, kStr)
		// 校验 Value
		VB, VS := CheckValue(fn, kStr, Value)
		if !VB {
			mapFlage = false
			property[kStr] = VS
			continue
		}
		property[kStr] = Value
	}
	if !mapFlage {
		return false, property
	}
	return true, property
}

// 校验 keyLengh 255
func keyLength255(val string) (bool, int) {
	lengthStatus, lengthCode := length255(val)
	if !lengthStatus {
		return false, lengthCode
	}
	return true, 0
}

// 判断字符串是否包含特殊字符 并返回错误code
func notSpecialCharacters(key string) (bool, int) {
	// 汉字的正则校验
	var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]|[\uFE30-\uFFA0]$")
	// $开头的 争取校验
	var zmRegexp = regexp.MustCompile("^[$a-zA-Z][a-zA-Z0-9_$]{0,}$")
	if hzRegexp.MatchString(key) || !zmRegexp.MatchString(key) {
		return false, 600011
	}
	return true, 0
}

// 判断字符串 长度 是否 99 ，并返回错误code
func length99(key string) (bool, int) {
	lengthInt := strings.Count(key, "") - 1
	if lengthInt > 99 || lengthInt < 1 {
		return false, 600010
	}
	return true, 0
}

// 判断字符串 长度 是否 99 ，并返回错误code
func length255(key string) (bool, int) {
	lengthInt := strings.Count(key, "") - 1
	if lengthInt > 255 || lengthInt == 0 {
		return false, 600019
	}
	return true, 0
}

// 判断是否包含关键字 ，并返回 是否包含  和 code ，方便日志打印
func keywords(key string) (bool, int) {
	keyArr := [6]string{"$lib", "$lib_version", "$platform", "$debug", "$is_login", "$original_id"}
	for _, v := range keyArr {
		if v == key {
			return false, 600012
		}
	}
	return true, 0
}
