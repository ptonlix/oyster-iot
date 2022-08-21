package controllers

import (
	"errors"
	jwt "oyster-iot/utils"

	oysterlog "oyster-iot/init/log"

	"github.com/beego/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

type PageParam struct {
	Pagesize int    `json:"pagesize" valid:"Max(255)"`
	Pagenum  int    `json:"pagenum" valid:"Max(255)"`
	Keyword  string `json:"keyword,omitempy" valid:"MaxSize(255)"`
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
	if code != 200 {
		c.Ctx.Output.Header(oysterlog.LogHeaderFlag, oysterlog.Failed)
	} else {
		c.Ctx.Output.Header(oysterlog.LogHeaderFlag, oysterlog.Success)
	}
	c.Data["json"] = mystruct
	c.ServeJSON()

}

//获取当前用户
func (c *BaseController) GetUserInfo() (int, string, error) {
	authorization := c.Ctx.Request.Header["Authorization"][0]
	userToken := authorization[len(jwt.JWTType)+1:]
	jwtInfo, err := jwt.ParseCliamsToken(userToken)
	if err != nil {
		logs.Error("Parse JWT Error!")
		return -1, "", errors.New("Parse JWT Error!")
	}

	return jwtInfo.Id, jwtInfo.Usertype, nil
}
