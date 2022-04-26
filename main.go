package main

import (
	_ "oyster-iot/init/log"
	_ "oyster-iot/routers"

	_ "oyster-iot/devaccess"
	_ "oyster-iot/init/mysql"

	_ "oyster-iot/init/cache"
	_ "oyster-iot/init/session"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	orm.Debug = true //开始数据库调试
	//logs.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

	beego.Run()
}
