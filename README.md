# Analysys Golang SDK

========

This is the official Golang SDK for Analysys.

# Golang SDK 目录说明：

* AnalysysAgent.go —— SDK文件
* src —— SDK 公共方法
* testDemo —— API调用演示

# Golang SDK 基础说明：

## 快速集成
如果您是第一次使用易观方舟产品，可以通过阅读本文快速了解此产品
1. 集成 SDK

* 先获取 易观 Golang SDK 

      go get github.com/analysys/ans-go-sdk

* 或更新 本地已经存在的

      go get -u github.com/analysys/ans-go-sdk  

* 引入 易观 Golang SDK

      import sdk "github.com/analysys/ans-go-sdk"

2. 设置初始化接口
```Go
AnalysysAgent = ans.InitAnalysysAgent(Collector lib.Collector, appid string, debugMode int)
```
* appid : 方舟项目对应的唯一标识
* debug : debug模式，有 0、1、2 三种枚举值。
    * 0 表示关闭 Debug 模式 （默认状态）
    * 1 表示打开 Debug 模式，但该模式下发送的数据仅用于调试，不计入平台数据统计
    * 2 表示打开 Debug 模式，该模式下发送的数据可计入平台数据统计 注意：发布版本时debug模式设置为0。
* Collector : 实时或者批量上报事件或者落文件到本地
    * InitSyncCollecter(uploadURL string)  实时上报数据
    * InitBatchCollector(uploadURL string,postNumber int) 批量上报数据
        * uploadURL ：服务上报地址
        * postNumber ：满足数量之后批量上报
    * InitSyncLogCollector(gerFolder string,gerRule string) 实时落文件到本地
    * InitBatchLogCollector(gerFolder string,gerRule string,postNumber int) 批量落文件到本地
        * gerFolder:文件存放的路径
        * gerRule:文件的切割规则，"D" 按天分割, "H",按小时分割 ,传"" ,则默认按小时切割
        * postNumber:批量落日志的数量，满足条数落日志到本地：


通过以上步骤您即可验证SDK是否已经集成成功。更多接口说明请您查看 API 文档。

更多Api使用方法参考：https://docs.analysys.cn/ark/integration/sdk/go

# 讨论
* 微信号：nlfxwz
* 钉钉群：30099866
* 邮箱：nielifeng@analysys.com.cn
  

# License

[gpl-3.0](https://www.gnu.org/licenses/gpl-3.0.txt)

