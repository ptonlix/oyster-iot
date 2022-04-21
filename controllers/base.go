package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Response(code int, msg interface{}) {
	type JSONStruct struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}
	mystruct := &JSONStruct{code, msg}
	c.Data["json"] = mystruct
	c.ServeJSON()
}
