package main

import (
	_ "oyster-iot/devaccess"
	_ "oyster-iot/init/constants"
	_ "oyster-iot/init/log"
	_ "oyster-iot/init/mysql"
	_ "oyster-iot/init/session"
	_ "oyster-iot/routers"

	_ "oyster-iot/init/cache"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

var (
	AppName      string // 应用名称
	AppVersion   string // 应用版本
	BuildVersion string // 编译版本
	BuildTime    string // 编译时间
	GitRevision  string // Git版本
	GitBranch    string // Git分支
	GoVersion    string // Golang信息
)

func main() {
	Version()
	//orm.Debug = true //开始数据库调试
	logs.SetLogger(logs.AdapterFile, `{"filename":"./log/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	beego.Run()
}

// Version 版本信息
func Version() {
	logs.Info("App Name:\t\t%s\n", AppName)
	logs.Info("App Version:\t%s\n", AppVersion)
	logs.Info("Build version:\t%s\n", BuildVersion)
	logs.Info("Build time:\t%s\n", BuildTime)
	logs.Info("Git revision:\t%s\n", GitRevision)
	logs.Info("Git branch:\t%s\n", GitBranch)
	logs.Info("Golang Version:\t%s\n", GoVersion)
}
