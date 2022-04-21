package main

import (
	"log"
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
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)

	beego.Run()
}
