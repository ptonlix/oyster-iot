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

type SomeDevAssets struct {
	AssetsNum  []string `json:"assets_nums" valid:"MaxSize(255)"`
	DeviceType string   `json:"dev_type" valid:"MaxSize(64)"`
}

func (c *BaseController) Response(code int, msg string, data ...interface{}) {
	type JSONStruct struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data,omitempty"` //忽略空值
	}
	var mystruct *JSONStruct
	if data != nil {
		mystruct = &JSONStruct{code, msg, data[0]}
	} else {
		mystruct = &JSONStruct{code, msg, nil}
	}

	c.Data["json"] = mystruct
	c.ServeJSON()

}
