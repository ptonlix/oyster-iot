package controllers

import (
	"encoding/json"
	"oyster-iot/models"
	"oyster-iot/services"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type OperLogController struct {
	BaseController
}

type PageParamLog struct {
	Pagesize int `json:"pagesize" valid:"Max(255)"`
	Pagenum  int `json:"pagenum" valid:"Max(255)"`
}

type SearchParamLog struct {
	Pagesize  int    `json:"pagesize" valid:"Max(255)"`
	Pagenum   int    `json:"pagenum" valid:"Max(255)"`
	Imei      string `json:"imei,omitempy" valid:"MaxSize(255)"`
	Context   string `json:"context,omitempy" valid:"MaxSize(255)"`
	Starttime string `json:"starttime,omitempy" valid:"MaxSize(255)"`
	Endtime   string `json:"endtime,omitempy" valid:"MaxSize(255)"`
}

type DeleteParamLog struct {
	Id int `json:"id" valid:"Required"`
}

type OperlogList struct {
	TotalNum   int                `json:"totalnum"`
	TotalPages int                `json:"totalpages"`
	List       *[]*models.Operlog `json:"list"`
}

// 搜索请求
func (op *OperLogController) Search() {
	//获取URL参数
	pageparam := SearchParamLog{}
	err := json.Unmarshal(op.Ctx.Input.RequestBody, &pageparam)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		op.Response(500, "系统内部错误")
		return
	}

	logs.Debug("pagesize is %#v, pagenum is %#v Imei is %#v  Context is %#v Starttime is %#v Endtime is %#v", pageparam.Pagesize, pageparam.Pagenum, pageparam.Imei,
		pageparam.Context, pageparam.Starttime, pageparam.Endtime)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		op.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		op.Response(400, "输入参数错误")
		return
	}
	sParam := &services.SearchOperlogParam{
		Starttime:     pageparam.Starttime,
		Endtime:       pageparam.Endtime,
		SearchImei:    pageparam.Imei,
		SearchContext: pageparam.Context,
	}
	var operService services.OperLogService

	totalNum, totalPages, operlogs, err := operService.GetOperlogBySearch(pageparam.Pagesize, pageparam.Pagenum, sParam)
	if err != nil {
		op.Response(400, "查找不到日志")
		return
	}

	retList := OperlogList{
		TotalNum:   totalNum,
		TotalPages: totalPages,
		List:       &operlogs,
	}
	op.Response(200, "获取日志列表成功", retList)
}

// 获取日志列表 ?pagesize=10&pagenum=1
func (op *OperLogController) List() {
	//获取URL参数
	pageparam := PageParamLog{}
	op.Ctx.Input.Bind(&pageparam.Pagesize, "pagesize")
	op.Ctx.Input.Bind(&pageparam.Pagenum, "pagenum")

	logs.Debug("pagesize is %#v, pagenum is %#v ", pageparam.Pagesize, pageparam.Pagenum)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		op.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		op.Response(400, "输入参数错误")
		return
	}

	var operService services.OperLogService

	totalNum, totalPages, operlogs, err := operService.GetOperlogByPage(pageparam.Pagesize, pageparam.Pagenum)
	if err != nil {
		op.Response(400, "查找不到日志")
		return
	}

	retList := OperlogList{
		TotalNum:   totalNum,
		TotalPages: totalPages,
		List:       &operlogs,
	}
	op.Response(200, "获取日志列表成功", retList)
}

func (op *OperLogController) Delete() {
	param := DeleteParamLog{}
	err := json.Unmarshal(op.Ctx.Input.RequestBody, &param)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		op.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&param)
	if err != nil {
		// handler error
		op.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		op.Response(400, "输入参数错误")
		return
	}

	opermodel := &models.Operlog{
		Id: param.Id,
	}
	var operService services.OperLogService
	if err := operService.Delete(opermodel); err != nil {
		op.Response(500, "数据库操作错误")
	}

	op.Response(200, "删除操作日志成功")
}
