package controllers

import (
	"encoding/json"
	"oyster-iot/models"
	"oyster-iot/services"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type BusinessController struct {
	BaseController
}

type BusinessInfo struct {
	BusinessID int    `json:"business_id"`
	Name       string `json:"business_name" valid:"Required;MaxSize(255);Match(/[\u4e00-\u9fa5a-zA-Z0-9_]{3,16}/)"`
	Remark     string `json:"remark" valid:"MaxSize(255)"`
}

// 增加一个业务
func (b *BusinessController) Add() {
	businessInfo := BusinessInfo{}
	err := json.Unmarshal(b.Ctx.Input.RequestBody, &businessInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		b.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	result, err := v.Valid(&businessInfo)
	if err != nil {
		// handler error
		b.Response(500, "系统内部错误")
		return
	}
	if !result {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		b.Response(400, "输入参数错误")
		return
	}
	businessM := &models.Business{
		Name:   businessInfo.Name,
		Remark: businessInfo.Remark,
	}
	//插入业务数据
	var businessS services.BusinessService
	id, err := businessS.Add(businessM)
	if err != nil {
		b.Response(500, "数据库操作错误")
		return
	}
	r := struct {
		BusinessID int64 `json:"business_id"`
	}{
		id,
	}
	b.Response(200, "新增业务成功", r)
}

func (b *BusinessController) Edit() {
	businessInfo := BusinessInfo{}
	err := json.Unmarshal(b.Ctx.Input.RequestBody, &businessInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		b.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	result, err := v.Valid(&businessInfo)
	if err != nil {
		// handler error
		b.Response(500, "系统内部错误")
		return
	}
	if !result {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		b.Response(400, "输入参数错误")
		return
	}
	businessM := &models.Business{
		Id:     businessInfo.BusinessID,
		Name:   businessInfo.Name,
		Remark: businessInfo.Remark,
	}
	//插入业务数据
	var businessS services.BusinessService
	err = businessS.Update(businessM)
	if err != nil {
		b.Response(500, "数据库操作错误")
		return
	}

	b.Response(200, "更新业务成功")
}

func (b *BusinessController) Delete() {
	businessInfo := BusinessInfo{}
	err := json.Unmarshal(b.Ctx.Input.RequestBody, &businessInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		b.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	result, err := v.Valid(&businessInfo)
	if err != nil {
		// handler error
		b.Response(500, "系统内部错误")
		return
	}
	if !result {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		b.Response(400, "输入参数错误")
		return
	}
	businessM := &models.Business{
		Id:     businessInfo.BusinessID,
		Name:   businessInfo.Name,
		Remark: businessInfo.Remark,
	}
	//插入业务数据
	var businessS services.BusinessService
	err = businessS.Delete(businessM)
	if err != nil {
		b.Response(500, "数据库操作错误")
		return
	}

	b.Response(200, "删除业务成功")
}

// 获取业务列表
func (b *BusinessController) List() {
	//获取URL参数
	pageparam := PageParam{}
	b.Ctx.Input.Bind(&pageparam.Pagesize, "pagesize")
	b.Ctx.Input.Bind(&pageparam.Pagenum, "pagenum")
	logs.Debug("pagesize is %#v, pagenum is %#v", pageparam.Pagesize, pageparam.Pagenum)

	// 校验输入参数是否合法
	v := validation.Validation{}
	result, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		b.Response(500, "系统内部错误")
		return
	}
	if !result {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		b.Response(400, "输入参数错误")
		return
	}
	//插入业务数据
	var businessS services.BusinessService
	business, err := businessS.GetBusinessByPage(pageparam.Pagesize, pageparam.Pagenum)
	if err != nil {
		b.Response(400, "查找不到设备")
		return
	}

	b.Response(200, "获取业务列表成功", business)
}
