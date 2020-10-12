package ans

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/analysys/ans-go-sdk/src/base"
	"github.com/analysys/ans-go-sdk/src/lib"
)

// 初始化 SDK
func InitAnalysysAgent(Collector lib.Collector, appid string, debugMode int) Analysys {
	base.BaseConfig.Debug = debugMode
	// 校验 appid
	lib.CheckAppid(appid)
	return Analysys{
		C:         Collector,
		Appid:     appid,
		DebugMode: debugMode,
	}
}

// 初始化同步上报
func InitSyncCollector(uploadURL string) *lib.SyncCollector {
	return lib.InitSyncCollector(uploadURL)
}

// 初始化异步上报
func InitBatchCollector(uploadURL string, postNumber int) *lib.BatchCollector {
	return lib.InitBatchCollector(uploadURL, postNumber)
}

// 初始化同步落日志到本地
func InitSyncLogCollector(gerFolder string, gerRule string) *lib.SyncLogCollector {
	return lib.InitSyncLogCollector(gerFolder, gerRule)
}

// 初始化异步落日志到本地
func InitBatchLogCollector(gerFolder string, gerRule string, postNumber int) *lib.BatchLogCollector {
	return lib.InitBatchLogCollector(gerFolder, gerRule, postNumber)
}

// SDK 初始化需要的传参
type Analysys struct {
	Appid         string
	C             lib.Collector
	DebugMode     int
	cancle        chan bool
	superProperty map[string]interface{}
}

// 数据初始化 五要素：
type DataPost struct {
	Appid    string                 `json:"appid"`
	Xwhat    string                 `json:"xwhat"`
	Xwho     string                 `json:"xwho"`
	Xwhen    int                    `json:"xwhen"`
	Xcontext map[string]interface{} `json:"xcontext"`
}

// xcontext 要素
type XcontextPost struct {
	Lib        string `json:"$lib"`
	LibVersion string `json:"$lib_version"`
	Debug      int    `json:"$debug"`
	IsLogin    bool   `json:"$is_login"`
	Platform   string `json:"$platform"`
}

// 超级属性  key value
func (ans *Analysys) RegisterSuperProperty(key string, value interface{}) bool {
	// 校验 key
	lib.CheckKey("$registerSuperProperty", key)
	// 校验 value
	valueStatusBool, valueStr := lib.CheckValue("$registerSuperProperty", key, value)
	var mapObj = make(map[string]interface{})
	if valueStatusBool {
		mapObj = lib.ToMap(key, value)
	} else {
		mapObj = lib.ToMap(key, valueStr)
	}
	superProperty := ans.superProperty
	// 超级属性的合并 map 合并 map
	ans.superProperty = lib.MapObjMerge(superProperty, mapObj)
	lib.SuccessLog("$registerSuperProperty", key, lib.MapToString(mapObj), 20002)
	return true
}

// 超级属性  property
func (ans *Analysys) RegisterSuperProperties(property map[string]interface{}) bool {
	// 校验 map的 string是否符合 key的要求
	superProperty := ans.superProperty
	Value := ""
	proBool, proValue := lib.CheckMap("$registerSuperProperties", property)
	if !proBool {
		ans.superProperty = lib.MapObjMerge(superProperty, proValue)
		Value = lib.MapToString(proValue)
	} else {
		ans.superProperty = lib.MapObjMerge(superProperty, property)
		Value = lib.MapToString(property)
	}
	lib.SuccessLog("$registerSuperProperties", "", Value, 20002)
	return true
}

// 获取超级属性
func (ans *Analysys) GetSuperProperty(key string) interface{} {
	//校验 key
	lib.CheckKey("$getSuperProperty", key)
	propertyValue, ok := ans.superProperty[key]
	if ok {
		lib.SuccessLog("$getSuperProperty", key, fmt.Sprint(propertyValue), 20010)
	} else {
		lib.SuccessLog("$getSuperProperty", key, "", 20009)
	}
	return propertyValue
}

//获取所有超级属性
func (ans *Analysys) GetSuperProperties() map[string]interface{} {
	lib.SuccessLog("$getSuperProperties", "", lib.MapToString(ans.superProperty), 20010)
	return ans.superProperty
}

// 移除单个超级属性
func (ans *Analysys) UnRegisterSuperProperty(key string) bool {
	//校验 key
	lib.CheckKey("$unRegisterSuperProperty", key)
	_, ok := ans.superProperty[key]
	if ok {
		delete(ans.superProperty, key)
		lib.SuccessLog("$unRegisterSuperProperty", "", key, 20003)
		return true
	}
	lib.SuccessLog("$unRegisterSuperProperty", "", key, 20011)
	return false
}

// 移除所有超级属性
func (ans *Analysys) ClearSuperProperties() bool {
	ans.superProperty = make(map[string]interface{})
	lib.SuccessLog("$clearSuperProperties", "", "", 20004)
	return true
}

/**
 * 设置用户的属性
 * @param distinctId 用户ID
 * @param isLogin 用户ID是否是登录 ID
 * @param properties 用户属性
 */

func (ans *Analysys) ProfileSet(distinctID string, isLogin bool, properties map[string]interface{}, platform string, upLoadTime int) bool {
	// 校验 properties ，不符合提醒，超长截取
	_, Properties := lib.CheckMap("$profile_set", properties)
	// 校验uploadTime 满足 13位长度的int  抛日志，这条数据不要，
	if !lib.CheckUploadTime("$profile_set", upLoadTime) || !lib.CheckDistinctID("$profile_set", distinctID) {
		return false
	}
	// 全部校验通过上报
	return ans.upLoad(distinctID, isLogin, "$profile_set", Properties, platform, upLoadTime, false)
}

/**
 * 首次设置用户的属性,该属性只在首次设置时有效
 * @param distinctId 用户ID
 * @param isLogin 用户ID是否是登录 ID
 * @param properties 用户属性
 */
func (ans *Analysys) ProfileSetOnce(distinctID string, isLogin bool, properties map[string]interface{}, platform string, upLoadTime int) bool {
	// 校验 properties ，不符合提醒，超长截取
	_, Properties := lib.CheckMap("$profile_set_once", properties)
	// 校验uploadTime 满足 13位长度的int  抛日志，这条数据不要，
	if !lib.CheckUploadTime("$profile_set_once", upLoadTime) || !lib.CheckDistinctID("$profile_set_once", distinctID) {
		return false
	}
	// 全部校验通过上报
	return ans.upLoad(distinctID, isLogin, "$profile_set_once", Properties, platform, upLoadTime, false)
}

/**
 * 为用户的一个或多个数值类型的属性累加一个数值
 * @param distinctId 用户ID
 * @param isLogin 用户ID是否是登录 ID
 * @param properties 用户属性
 */

func (ans *Analysys) ProfileIncrement(distinctID string, isLogin bool, properties map[string]int, platform string, upLoadTime int) bool {
	// 需要把 参数  map[string]int 转换成 map[string] interface{}
	Properties := lib.MapIntToMapInter(properties)
	// 校验 properties ，不符合提醒，超长截取
	_, Properties = lib.CheckMap("$profile_increment", Properties)
	// 校验uploadTime 满足 13位长度的int  抛日志，这条数据不要，
	if !lib.CheckUploadTime("$profile_increment", upLoadTime) || !lib.CheckDistinctID("$profile_increment", distinctID) {
		return false
	}
	// 全部校验通过上报
	return ans.upLoad(distinctID, isLogin, "$profile_increment", Properties, platform, upLoadTime, false)
}

/**
 * 追加用户列表类型的属性
 * @param distinctId 用户ID
 * @param isLogin 用户ID是否是登录 ID
 * @param properties 用户属性
 */

func (ans *Analysys) ProfileAppend(distinctID string, isLogin bool, properties map[string]interface{}, platform string, upLoadTime int) bool {
	// 校验 properties ，不符合提醒，超长截取
	_, Properties := lib.CheckMap("$profile_append", properties)
	// 校验uploadTime 满足 13位长度的int  抛日志，这条数据不要，
	if !lib.CheckUploadTime("$profile_append", upLoadTime) || !lib.CheckDistinctID("$profile_append", distinctID) {
		return false
	}
	// 全部校验通过上报
	return ans.upLoad(distinctID, isLogin, "$profile_append", Properties, platform, upLoadTime, false)
}

/**
 * 删除用户某一个属性
 * @param distinctId 用户ID
 * @param isLogin 用户ID是否是登录 ID
 * @param property 用户属性名称
 * @throws AnalysysException
 */

func (ans *Analysys) ProfileUnSet(distinctID string, isLogin bool, property string, platform string, upLoadTime int) bool {
	// 校验属性值
	lib.CheckKey("$profile_unset", property)
	// 组建map value为空
	properties := make(map[string]interface{})
	properties[property] = ""
	if !lib.CheckUploadTime("$profile_unset", upLoadTime) || !lib.CheckDistinctID("$profile_unset", distinctID) {
		return false
	}
	return ans.upLoad(distinctID, isLogin, "$profile_unset", properties, platform, upLoadTime, false)
}

/**
 * 删除用户所有属性
 * @param distinctId 用户ID
 * @param isLogin 用户ID是否是登录 ID
 * @throws AnalysysException
 */

func (ans *Analysys) ProfileDelete(distinctID string, isLogin bool, platform string, upLoadTime int) bool {
	// 创建新的 空 map 上传
	properties := make(map[string]interface{})
	if !lib.CheckUploadTime("$profile_delete", upLoadTime) || !lib.CheckDistinctID("$profile_delete", distinctID) {
		return false
	}
	return ans.upLoad(distinctID, isLogin, "$profile_delete", properties, platform, upLoadTime, false)
}

/**
 * 关联用户匿名ID和登录ID
 * @param aliasId 用户登录ID
 * @param distinctId 用户匿名ID
 * @throws AnalysysException
 */

func (ans *Analysys) Alias(aliasID string, distinctID string, platform string, upLoadTime int) bool {
	// 组建 map
	properties := make(map[string]interface{})
	properties["$original_id"] = distinctID
	if !lib.CheckUploadTime("$alias", upLoadTime) || !lib.CheckAliasID("$alias", aliasID) || !lib.CheckDistinctID("$alias", distinctID) {
		return false
	}
	return ans.upLoad(aliasID, true, "$alias", properties, platform, upLoadTime, false)
}

/**
 * 追踪用户多个属性的事件
 * @param distinctId 用户ID
 * @param isLogin 用户ID是否是登录 ID
 * @param eventName 事件名称
 * @param properties 事件属性
 * @throws AnalysysException
 */
func (ans *Analysys) Track(distinctID string, isLogin bool, eventName string, properties map[string]interface{}, platform string, upLoadTime int) bool {
	lib.CheckEventName(eventName)
	_, Properties := lib.CheckMap(eventName, properties)
	// 允许空map
	if !lib.CheckUploadTime(eventName, upLoadTime) || !lib.CheckDistinctID(eventName, distinctID) {
		return false
	}
	return ans.upLoad(distinctID, isLogin, eventName, Properties, platform, upLoadTime, true)
}

/**
 * 立即发送所有收集的信息到服务器
 */
func (ans *Analysys) Flush() bool {
	return ans.C.Send(ans.DebugMode)
}

/**
 * 获取当前时间，方便调用 API 传当前时间
 */

func CurrentTime() int {
	return int(time.Now().UnixNano() / 1e6)
}

// 上传方法 ，此方法是对现有数据的组装和上传条件的校验，符合条件后 触发 发送 send
func (ans *Analysys) upLoad(distinctID string, isLogin bool, eventName string, properties map[string]interface{}, platform string, upLoadTime int, merFlag bool) bool {
	// 组建五要素
	DataComponent := &DataPost{}
	DataComponent.Appid = ans.Appid
	DataComponent.Xwhat = eventName
	DataComponent.Xwho = distinctID
	DataComponent.Xwhen = upLoadTime
	// 属性合并，只有 track 合并超级属性
	if merFlag == true {
		properties = lib.MapObjMerge(ans.superProperty, properties)
	}
	// 校验 platform
	platform = lib.CheckPlatform(platform)
	// 组建 xcontext ,用map 格式，方便与 properties 进行数据合并
	contextComponent := make(map[string]interface{})
	contextComponent["$lib"] = "Go"
	contextComponent["$lib_version"] = "4.3.1"
	contextComponent["$debug"] = ans.DebugMode
	contextComponent["$is_login"] = isLogin
	contextComponent["$platform"] = platform
	// 遍历map 合并   合并完的 赋值 xcontext
	DataComponent.Xcontext = lib.MapObjMerge(properties, contextComponent) // 后循环的会覆盖之前循环的，比方说 有$platform "node",会覆盖之前的 "go"
	data, _ := json.Marshal(DataComponent)
	dataStr := string(data)
	// 数据之后推进一个 数组，设置同步发送，立即上报
	ans.cancle = make(chan bool, 1)
	return ans.C.UpLoad(dataStr, ans.DebugMode)
}
