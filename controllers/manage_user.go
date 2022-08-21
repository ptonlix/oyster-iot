package controllers

import (
	"encoding/json"
	"oyster-iot/models"
	"oyster-iot/services"

	bcrypt "oyster-iot/utils"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type ManageUserController struct {
	BaseController
}

type ManageUserInfo struct {
	Username  string `json:"username" valid:"Required;MaxSize(100);Match(/[a-zA-Z0-9_]{3,100}/)"`
	Password  string `json:"password" valid:"Required;MinSize(6);MaxSize(100)"`
	Email     string `json:"email" valid:"Required;Email"`
	Firstname string `json:"first_name" valid:"Required;MaxSize(64)"`
	Lastname  string `json:"last_name" valid:"Required;MaxSize(64)"`
	Mobile    string `json:"mobile" valid:"Required;Mobile"`
	Remark    string `json:"remark" valid:"MaxSize(256)"`
	IsAdmin   bool   `json:"is_admin" valid:"Required;"`
	Enabled   bool   `json:"enabled" valid:"Required;"`
}

type EditManageUserInfo struct {
	Id        int    `json:"id" valid:"Required"`
	Username  string `json:"username" valid:"Required;MaxSize(100);Match(/[a-zA-Z0-9_]{3,100}/)"`
	Email     string `json:"email" valid:"Required;Email"`
	Firstname string `json:"first_name" valid:"Required;MaxSize(64)"`
	Lastname  string `json:"last_name" valid:"Required;MaxSize(64)"`
	Mobile    string `json:"mobile" valid:"Required;Mobile"`
	Remark    string `json:"remark" valid:"MaxSize(256)"`
	IsAdmin   bool   `json:"is_admin" valid:"Required;"`
	Enabled   bool   `json:"enabled" valid:"Required;"`
}

type ManageUserList struct {
	TotalNum   int              `json:"totalnum"`
	TotalPages int              `json:"totalpages"`
	List       *[]*models.Users `json:"list"`
}

type ManageUserIndex struct {
	Id int `json:"id" `
}

type ManageResetPasswdInfo struct {
	Id       int    `json:"id"`
	Password string `json:"password" valid:"Required;MinSize(6);MaxSize(100)"`
}

//添加用户
func (u *ManageUserController) Add() {
	userInfo := UserInfo{}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &userInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		u.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&userInfo)
	if err != nil {
		// handler error
		u.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		u.Response(400, "输入参数错误")
		return
	}

	// 创建用户
	user := &models.Users{
		Enabled:   userInfo.Enabled,
		Email:     userInfo.Email,
		Username:  userInfo.Username,
		Password:  bcrypt.HashAndSalt([]byte(userInfo.Password)),
		Firstname: userInfo.Firstname,
		Lastname:  userInfo.Lastname,
		Mobile:    userInfo.Mobile,
		Remark:    userInfo.Remark,
		IsAdmin:   userInfo.IsAdmin,
	}

	var userService services.UserService
	if err := userService.Add(user); err != nil {
		u.Response(500, "数据库操作错误")
		return
	}

	u.Response(200, "添加用户成功")

}

//修改用户
func (u *ManageUserController) Edit() {
	userInfo := EditUserInfo{}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &userInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		u.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&userInfo)
	if err != nil {
		// handler error
		u.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		u.Response(400, "输入参数错误")
		return
	}

	var UserService services.UserService
	user, err := UserService.GetUserById(userInfo.Id)
	if err != nil {
		u.Response(500, "数据库操作错误")
		return
	}

	user.Username = userInfo.Username
	user.Email = userInfo.Email
	user.Firstname = userInfo.Firstname
	user.Lastname = userInfo.Lastname
	user.Mobile = userInfo.Mobile
	user.IsAdmin = userInfo.IsAdmin
	user.Enabled = userInfo.Enabled
	user.Remark = userInfo.Remark

	if err := UserService.Update(user); err != nil {
		u.Response(500, "数据库操作错误")
	}

	//
	u.Response(200, "更新用户成功")

}

// 重置密码
func (u *ManageUserController) ResetPassword() {
	resetPasswd := ResetPasswdInfo{}
	logs.Info(string(u.Ctx.Input.RequestBody))
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &resetPasswd)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		u.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&resetPasswd)
	if err != nil {
		// handler error
		u.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		u.Response(400, "输入参数错误")
		return
	}
	var UserService services.UserService
	user, err := UserService.GetUserById(resetPasswd.Id)
	if err != nil {
		u.Response(500, "数据库操作错误")
		return
	}
	user.Password = bcrypt.HashAndSalt([]byte(resetPasswd.Password))
	if err := UserService.Update(user); err != nil {
		u.Response(500, "数据库操作错误")
	}

	u.Response(200, "重置用户密码成功")
}

//删除用户
func (u *ManageUserController) Delete() {
	userIndex := UserIndex{}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &userIndex)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		u.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&userIndex)
	if err != nil {
		// handler error
		u.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		u.Response(400, "输入参数错误")
		return
	}

	user := &models.Users{
		Id: userIndex.Id,
	}
	var UserService services.UserService
	err = UserService.Delete(user)
	if err != nil {
		u.Response(500, "数据库操作错误")
		return
	}
	u.Response(200, "删除用户成功")
}

func (u *ManageUserController) List() {
	//获取URL参数
	pageparam := PageParam{}
	u.Ctx.Input.Bind(&pageparam.Pagesize, "pagesize")
	u.Ctx.Input.Bind(&pageparam.Pagenum, "pagenum")
	u.Ctx.Input.Bind(&pageparam.Keyword, "keyword")
	logs.Debug("pagesize is %#v, pagenum is %#v keyword is %#v", pageparam.Pagesize, pageparam.Pagenum, pageparam.Keyword)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		u.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		u.Response(400, "输入参数错误")
		return
	}

	// 查询数据
	var UserService services.UserService

	totalNum, totalPages, users, err := UserService.GetUserByPageAndKey(pageparam.Pagesize, pageparam.Pagenum, pageparam.Keyword)
	if err != nil {
		u.Response(400, "查找不到用户")
		return
	}

	retList := UserList{
		TotalNum:   totalNum,
		TotalPages: totalPages,
		List:       &users,
	}
	u.Response(200, "获取用户列表成功", retList)
}
