package base

var BaseConfig = BaseConf{
	Appid: "",
	Debug: 0,
}

type BaseConf struct {
	Appid string
	Debug int `json:"$debug"`
}
