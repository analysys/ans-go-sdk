package lib

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/analysys/ans-go-sdk/src/base"
)

func SuccessMessage() map[interface{}]string {
	var successM map[interface{}]string
	successM = make(map[interface{}]string)
	successM[20001] = "Send message success"
	successM[20002] = "{FN}: set success ({VALUE})"
	successM[20003] = "{FN}:({VALUE}) delete success"
	successM[20004] = "{FN}:clear success"
	successM[20009] = "{FN}:[{KEY}] : get failed"
	successM[20010] = "{FN}:[{KEY}] : get success. ({VALUE})"
	successM[20011] = "{FN}:({VALUE}) delete failed"
	successM[20012] = "Send message to server : {KEY} \n" +
		"data:{VALUE}"
	successM[20014] = "Save message to file: {KEY} \n" +
		"data:{VALUE}"
	successM[20015] = "Save message success"
	successM[20016] = "{FN}: set failed ({VALUE})"
	return successM
}

func ErrorMessage() map[interface{}]string {
	var errorM map[interface{}]string
	errorM = make(map[interface{}]string)
	errorM[60001] = "Send message failed,{KEY}"
	errorM[60006] = "{FN}:The length of the property key (string[{KEY}]) needs to be 1-255 !"
	errorM[60007] = "{FN}:The length of the property[{KEY}] value (string[{VALUE}]) needs Less than 100 !"
	errorM[60008] = "{FN}:Property value invalid of key[{KEY}], support type: string,int,float64,bool,[]string "
	errorM[60009] = "{FN}:The length of the property key (string[{KEY}]) needs to be 1-125 !"
	errorM[600010] = "{FN}:The length of the property key (string[{KEY}]) needs to be 1-99 !"
	errorM[600011] = "{FN}:[{KEY}] does not conform to naming rules!"
	errorM[600012] = "{FN}:Property key invalid, nonsupport value: $lib/$lib_version/$platform/$first_visit_time/$debug/$is_login \n" +
		"current KEY:{KEY}"
	errorM[600013] = "{FN}:Property value invalid of key[{KEY}], support type: slice with String([]string) as its internal element \n" +
		"current value:{VALUE}\n" +
		"current type: {VALUETYPE}"
	errorM[600019] = "{FN}:The length of the property[{KEY}] value (string[{VALUE}]) needs to be 1-255 !"
	errorM[600021] = "{FN}:Property value invalid of key[{KEY}],  Cannot be set to empty \n" +
		"current value:{VALUE}"
	errorM[600023] = "{FN}:{KEY} do not macth {KEY} rules!"
	errorM[600024] = "{FN}:{KEY} do not macth {VALUE} rules!, length needs 13!"
	errorM[600027] = "Save message failed"
	return errorM
}

func SuccessLog(fn string, key string, value string, successCode int) {
	DebugMode := base.BaseConfig.Debug
	if DebugMode == 1 || DebugMode == 2 {
		mesTemp := SuccessMessage()[successCode]
		// fmt.Println("type:", reflect.TypeOf(fn))
		showMsg := strings.Replace(mesTemp, "{FN}", fn, -1)
		showMsg = strings.Replace(showMsg, "{KEY}", key, -1)
		showMsg = strings.Replace(showMsg, "{VALUE}", value, -1)
		fmt.Println(showMsg)
	}
}

func ErrorLog(fn string, key string, value string, errorCode int) {
	DebugMode := base.BaseConfig.Debug
	if DebugMode == 1 || DebugMode == 2 {
		valueType := fmt.Sprint(reflect.TypeOf(value).Kind())
		mesTemp := ErrorMessage()[errorCode]
		// fmt.Println("type:", reflect.TypeOf(fn))
		showMsg := strings.Replace(mesTemp, "{FN}", fn, -1)
		showMsg = strings.Replace(showMsg, "{KEY}", key, -1)
		showMsg = strings.Replace(showMsg, "{VALUE}", value, -1)
		showMsg = strings.Replace(showMsg, "{VALUETYPE}", valueType, -1)
		fmt.Println(showMsg)
	}
}
