package main

import (
	_ "oyster-iot/routers"

	_ "oyster-iot/devaccess"
	_ "oyster-iot/init/mysql"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	orm.Debug = true //开始数据库调试
	beego.Run()
}
