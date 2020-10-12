// 完整 demo ，使用前 请 注释调 test.go .确保只有一个 main 包
package main

import (
	"fmt"

	ans "github.com/analysys/ans-go-sdk"
)

func main() {
	// 初始化 AnalysysAgent
	appId := "1234"
	upLoadURL := "https://arksdk.analysys.cn:4089"
	postNumber := 2
	debugMode := 2
	batchCollector := ans.InitBatchCollector(upLoadURL, postNumber)
	AnalysysAgent := ans.InitAnalysysAgent(batchCollector, appId, debugMode)
	// 初始化完成
	distinctId := "1234567890987654321"
	platform := "Android"
	//浏览商品
	trackPropertie := map[string]interface{}{
		"$ip": "112.112.112.112", //IP地址
	}
	bookList := []string{"Go语言编程"}
	trackPropertie["productName"] = bookList //商品列表
	trackPropertie["productType"] = "Go书籍"   //商品类别
	trackPropertie["producePrice"] = 80      //商品价格
	trackPropertie["shop"] = "xx网上书城"        //店铺名称
	viewStatus := AnalysysAgent.Track(distinctId, true, "ViewProduct", trackPropertie, platform, ans.CurrentTime())
	fmt.Println("是否浏览商品", viewStatus)
	//用户注册登录
	registerId := "ABCDEF123456789"
	aliasStatus := AnalysysAgent.Alias(registerId, distinctId, platform, ans.CurrentTime()) //设置公共属性
	fmt.Println("是否关联用户", aliasStatus)
	superPropertie := map[string]interface{}{
		"sex": "male", //性别
		"age": 23,     //年龄
	}
	registerStatus := AnalysysAgent.RegisterSuperProperties(superPropertie) //用户信息
	fmt.Println("是否设置用户信息", registerStatus)
	profiles := map[string]interface{}{
		"$city":     "北京",    //城市
		"$province": "北京",    //省份
		"nickName":  "昵称123", //昵称
		"userLevel": 0,       //用户级别
		"userPoint": 0,       //用户积分
	}
	interestList := []string{"户外活动", "足球赛事", "游戏"}
	profiles["interest"] = interestList                                               //用户兴趣爱好
	AnalysysAgent.ProfileSet(registerId, true, profiles, platform, ans.CurrentTime()) //用户注册时间
	profile_age := map[string]interface{}{
		"registerTime": "20180101101010",
	}
	setOnceStatus := AnalysysAgent.ProfileSetOnce(registerId, true, profile_age, platform, ans.CurrentTime())
	fmt.Println("是否关联设置用户首次属性", setOnceStatus)

	//重新设置公共属性
	AnalysysAgent.ClearSuperProperties()

	superPropertie = map[string]interface{}{
		"userLevel": 0, //用户级别
		"userPoint": 0, //用户积分
	}
	AnalysysAgent.RegisterSuperProperties(superPropertie)
	//再次浏览商品
	trackPropertie["$ip"] = "112.112.112.112" //IP地址
	bookList = []string{"Go语言编程"}
	trackPropertie["productName"] = bookList //商品列表
	trackPropertie["productType"] = "Go书籍"   //商品类别
	trackPropertie["producePrice"] = 80      //商品价格
	trackPropertie["shop"] = "xx网上书城"        //店铺名称
	AnalysysAgent.Track(registerId, true, "ViewProduct", trackPropertie, platform, ans.CurrentTime())
	//订单信息
	trackPropertie["orderId"] = "ORDER_12345"
	trackPropertie["price"] = 80
	orderStatus := AnalysysAgent.Track(registerId, true, "Order", trackPropertie, platform, ans.CurrentTime())
	fmt.Println("是否下订单", orderStatus)

	//支付信息
	trackPropertie["orderId"] = "ORDER_12345"
	trackPropertie["productName"] = "Go语言编程"
	trackPropertie["productType"] = "Go书籍"
	trackPropertie["producePrice"] = 8
	trackPropertie["shop"] = "xx网上书城"
	trackPropertie["productNumber"] = 1
	trackPropertie["price"] = 80
	trackPropertie["paymentMethod"] = "AliPay"
	payStatus := AnalysysAgent.Track(registerId, true, "Payment", trackPropertie, platform, ans.CurrentTime())
	fmt.Println("是否付款", payStatus)
	AnalysysAgent.Flush()

}
