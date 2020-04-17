package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

// 返回标准 时间 map 对象
func getTime() map[string]string {
	year := strconv.Itoa(time.Now().Year())
	var monthIntStr, dayIntString, hourIntStr string
	monthInt := int(time.Now().Month())
	if monthInt < 10 {
		monthIntStr = "0" + fmt.Sprint(monthInt)
	} else {
		monthIntStr = fmt.Sprint(monthInt)
	}
	dayInt := int(time.Now().Day())
	if dayInt < 10 {
		dayIntString = "0" + fmt.Sprint(dayInt)
	} else {
		dayIntString = fmt.Sprint(dayInt)
	}
	hourInt := int(time.Now().Hour())
	if hourInt < 10 {
		hourIntStr = "0" + fmt.Sprint(hourInt)
	} else {
		hourIntStr = fmt.Sprint(hourInt)
	}
	countryCapitalMap := map[string]string{"year": year, "month": monthIntStr, "day": dayIntString, "hour": hourIntStr}
	return countryCapitalMap
}

//调用os.MkdirAll递归创建文件夹
func createFile(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 组建map
func ToMap(key string, value interface{}) map[string]interface{} {
	mapObj := make(map[string]interface{})
	mapObj[key] = value
	return mapObj
}

// map 合并map
func MapObjMerge(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {
	mapContext := make(map[string]interface{})
	for k, v := range map1 {
		mapContext[k] = v
	}
	// 后循环的会覆盖之前循环的，比方说 有$platform "node",会覆盖之前的 "go"
	for k, v := range map2 {
		mapContext[k] = v
	}
	return mapContext
}

// map转换 string
func MapToString(map1 map[string]interface{}) string {
	mapObj := make(map[string]interface{})
	mapObj = map1
	mapObjValue, err := json.Marshal(mapObj)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
	}
	return string(mapObjValue)
}

func MapIntToMapInter(map1 map[string]int) map[string]interface{} {
	mapIN := make(map[string]interface{})
	for k, v := range map1 {
		mapIN[k] = v
	}
	return mapIN
}
