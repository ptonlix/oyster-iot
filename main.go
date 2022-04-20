package main

import (
	_ "oyster-iot/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

