package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"

	"oyster-iot/init/constants"
	"oyster-iot/models"
	"oyster-iot/services"
	access "oyster-iot/services/thirdpartaccess"
	bcrypt "oyster-iot/utils"
)

type WxAccessController struct {
	BaseController
}

type WxLogin struct {
	Code     string `json:"code" valid:"Required;"`
	Nickname string `json:"nickname" valid:"Required;"`
	Phone    string `json:"phone" valid:"Required;Mobile"`
}

type WxCode struct {
	Code string `json:"code" valid:"Required;"`
}

// 登陆
func (wx *WxAccessController) WxLogin() {
	logs.Info(string(wx.Ctx.Input.RequestBody))
	wxInfo := WxLogin{}
	err := json.Unmarshal(wx.Ctx.Input.RequestBody, &wxInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		wx.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&wxInfo)
	if err != nil {
		// handler error
		wx.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		wx.Response(400, "输入参数错误")
		return
	}
	manager, err := access.NewWxLoginService(constants.WxConf.Appid, constants.WxConf.AppSecret, constants.WxConf.ApiHost)
	if err != nil {
		logs.Error(err)
		wx.Response(500, "系统内部错误")
		return
	}

	info, err := manager.QuerySession(wxInfo.Code)
	if err != nil {
		if err != nil {
			logs.Error(err)
			wx.Response(500, "请求微信服务器错误", err)
			return
		}
	}
	if info.Errcode != 0 {
		wx.Response(500, "请求微信服务器错误", info)
		return
	}

	// 查找是否存在该用户
	var UserService services.UserService
	userInfo, err := UserService.GetUserByUsername(wxInfo.Phone)
	if err != nil {
		// 没有找到，生成新用户
		// 创建用户
		user := &models.Users{
			Enabled:   true,
			Username:  wxInfo.Phone,
			Password:  bcrypt.HashAndSalt([]byte(info.SessionKey)),
			Firstname: "WeChat",
			Lastname:  wxInfo.Nickname,
			Mobile:    wxInfo.Phone,
			Wxunionid: info.Unionid,
			Wxopenid:  info.Openid,
			IsAdmin:   false,
		}

		if err := UserService.Add(user); err != nil {
			wx.Response(500, "数据库操作错误", err)
			return
		}

	} else {
		userInfo.Password = bcrypt.HashAndSalt([]byte(info.SessionKey))
		if err := UserService.Update(userInfo); err != nil {
			wx.Response(500, "数据库操作错误", err)
		}
	}

	wx.Response(200, "微信登陆成功", info)

}

// 获取手机号
func (wx *WxAccessController) GetPhoneNumber() {
	logs.Info(string(wx.Ctx.Input.RequestBody))
	wxInfo := WxCode{}
	err := json.Unmarshal(wx.Ctx.Input.RequestBody, &wxInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		wx.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&wxInfo)
	if err != nil {
		// handler error
		wx.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		wx.Response(400, "输入参数错误")
		return
	}

	manager, err := access.NewWxLoginService(constants.WxConf.Appid, constants.WxConf.AppSecret, constants.WxConf.ApiHost)
	if err != nil {
		logs.Error(err)
		wx.Response(500, "系统内部错误")
		return
	}

	info, err := manager.GetPhoneNumber(wxInfo.Code)
	if err != nil {
		if err != nil {
			logs.Error(err)
			wx.Response(500, "请求微信服务器错误", err)
			return
		}
	}
	if info.Errcode != 0 {
		wx.Response(500, "请求微信服务器错误", info)
		return
	}

	wx.Response(200, "获取用户手机号成功", info)
}
