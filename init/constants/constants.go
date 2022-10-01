package constants

import (
	"oyster-iot/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

const (
	LogFileDir  = "./log/"
	LogFileName = "oyster-syslog.log"
)

var QiniuAK string
var QiniuSK string

type InfluxConfig struct {
	Host   string
	Token  string
	Org    string
	Bucket string
	Online bool
}

type WxConfig struct {
	Appid     string // 小程序配置文件
	AppSecret string
	ApiHost   string
}

var InfluxConf InfluxConfig
var WxConf WxConfig

func init() {
	QiniuAK, _ = beego.AppConfig.String("qiniuak")
	QiniuSK, _ = beego.AppConfig.String("qiniusk")
	logs.Info("QiniuAK: %s   QiniuSK: %s ", QiniuAK, QiniuSK)

	//获取influxDB配置
	InfluxConf.Host, _ = beego.AppConfig.String("influxhost")
	InfluxConf.Token, _ = beego.AppConfig.String("influxtoken")
	InfluxConf.Org, _ = beego.AppConfig.String("influxorg")
	InfluxConf.Bucket, _ = beego.AppConfig.String("influxbucket")

	logs.Info("InfluxDB Config: Host: %s | Token: %s | Organization: %s | Bucket: %s|",
		InfluxConf.Host, InfluxConf.Token, InfluxConf.Org, InfluxConf.Bucket)

	InfluxConf.Online = utils.DetectInfluxDBOnline(InfluxConf.Host, InfluxConf.Token, InfluxConf.Org)

	// 获取微信配置
	WxConf.Appid, _ = beego.AppConfig.String("wxappid")
	WxConf.AppSecret, _ = beego.AppConfig.String("wxappsecret")
	WxConf.ApiHost, _ = beego.AppConfig.String("wxapihost")

}
