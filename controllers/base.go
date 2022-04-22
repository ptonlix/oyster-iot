package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

type PageParam struct {
	Pagesize int `json:"pagesize" valid:"Max(255)"`
	Pagenum  int `json:"pagenum" valid:"Max(255)"`
}

func (c *BaseController) Response(code int, msg string, data ...interface{}) {
	type JSONStruct struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	mystruct := &JSONStruct{code, msg, data[0]}
	c.Data["json"] = mystruct
	c.ServeJSON()

}
