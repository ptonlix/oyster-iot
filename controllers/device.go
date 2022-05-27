package controllers

import (
	"encoding/json"
	"oyster-iot/models"
	"oyster-iot/services"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"

	uuid "oyster-iot/utils"
)

type DeviceController struct {
	BaseController
}

// TODO:正则校验，防止SQL注入
type DeviceInfo struct {
	AssetsNum  string `json:"assets_num"  valid:"MaxSize(255)"`
	Token      string `json:"token"       valid:"Required;MaxSize(255)"`
	DeviceName string `json:"device_name" valid:"Required;MaxSize(255);Match(/[\u4e00-\u9fa5a-zA-Z0-9_]{3,16}/)"`
	Protocol   string `json:"protocol"    valid:"Required;MaxSize(64)"`
	Publish    string `json:"publish"     valid:"Required;MaxSize(64)"`
	Subscribe  string `json:"subscribe"   valid:"Required;MaxSize(64)"`
	Type       string `json:"type"        valid:"Required;MaxSize(64)"`
}
type DeviceOfBusiness struct {
	BusinessId int `json:"business_id" valid:"Required"`
}

type DevUpBusiness struct {
	BusinessId int      `json:"business_id" valid:"Min(0)"`
	AssetsNum  []string `json:"assets_num"  valid:"Required;MaxSize(255)"`
}

type DeviceList struct {
	TotalNum   int               `json:"totalnum"`
	TotalPages int               `json:"totalpages"`
	DeviceList *[]*models.Device `json:"list"`
}

func (d *DeviceController) Add() {
	deviceInfo := DeviceInfo{}
	err := json.Unmarshal(d.Ctx.Input.RequestBody, &deviceInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		d.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&deviceInfo)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}

	// 生成资产编码
	uNum := uuid.GetUuid()
	// 插入设备数据
	device := &models.Device{
		AssetsNum:  uNum,
		Token:      deviceInfo.Token,
		DeviceName: deviceInfo.DeviceName,
		Protocol:   deviceInfo.Protocol,
		Publish:    deviceInfo.Publish,
		Subscribe:  deviceInfo.Subscribe,
		Type:       deviceInfo.Type,
	}
	var deviceService services.DeviceService

	if err := deviceService.Add(device); err != nil {
		d.Response(500, "数据库操作错误")
		return
	}

	// return
	r := struct {
		AssetsNum string `json:"assets_num"`
	}{
		uNum,
	}
	d.Response(200, "添加设备成功", r)
}

// 编辑设备
func (d *DeviceController) Edit() {
	deviceInfo := DeviceInfo{}
	err := json.Unmarshal(d.Ctx.Input.RequestBody, &deviceInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		d.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&deviceInfo)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}

	// 插入设备数据
	var device *models.Device
	var deviceService services.DeviceService

	if device, err = deviceService.GetDeviceByAssetsNum(deviceInfo.AssetsNum); err != nil {
		d.Response(500, "查找不到该设备")
		return
	}

	device.Token = deviceInfo.Token
	device.DeviceName = deviceInfo.DeviceName
	device.Protocol = deviceInfo.Protocol
	device.Publish = deviceInfo.Publish
	device.Subscribe = deviceInfo.Subscribe
	device.Type = deviceInfo.Type

	if err := deviceService.Update(device); err != nil {
		d.Response(500, "数据库操作错误")
	}

	d.Response(200, "更新设备成功")
}

// 删除设备
func (d *DeviceController) Delete() {
	deviceInfo := DeviceInfo{}
	err := json.Unmarshal(d.Ctx.Input.RequestBody, &deviceInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		d.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&deviceInfo)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}

	// 删除设备数据
	var device *models.Device
	var deviceService services.DeviceService

	if device, err = deviceService.GetDeviceByAssetsNum(deviceInfo.AssetsNum); err != nil {
		d.Response(500, "查找不到该设备")
		return
	}

	if err := deviceService.Delete(device); err != nil {
		d.Response(500, "数据库操作错误")
	}

	d.Response(200, "删除设备成功")
}

// 获取设备列表
func (d *DeviceController) List() {
	//获取URL参数
	pageparam := PageParam{}
	d.Ctx.Input.Bind(&pageparam.Pagesize, "pagesize")
	d.Ctx.Input.Bind(&pageparam.Pagenum, "pagenum")
	logs.Debug("pagesize is %#v, pagenum is %#v", pageparam.Pagesize, pageparam.Pagenum)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}

	// 获取设备数据
	var deviceService services.DeviceService

	totalNum, totalPages, devices, err := deviceService.GetDevicesByPage(pageparam.Pagesize, pageparam.Pagenum)
	if err != nil {
		d.Response(400, "查找不到设备")
		return
	}

	retList := DeviceList{
		TotalNum:   totalNum,
		TotalPages: totalPages,
		DeviceList: &devices,
	}
	d.Response(200, "获取设备列表成功", retList)
}

// 通过业务ID，获取设备列表
func (d *DeviceController) ListForBusiness() {
	deviceOfBusiness := DeviceOfBusiness{}
	err := json.Unmarshal(d.Ctx.Input.RequestBody, &deviceOfBusiness)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		d.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&deviceOfBusiness)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}

	// 获取设备数据
	var deviceService services.DeviceService

	devices, err := deviceService.GetDeviceByBusiness(deviceOfBusiness.BusinessId)
	if err != nil {
		d.Response(400, "查找不到设备")
		return
	}

	d.Response(200, "获取设备列表成功", devices)
}

// 获取空业务的设备列表
func (d *DeviceController) ListForNilBusiness() {
	//获取URL参数
	pageparam := PageParam{}
	d.Ctx.Input.Bind(&pageparam.Pagesize, "pagesize")
	d.Ctx.Input.Bind(&pageparam.Pagenum, "pagenum")
	logs.Debug("pagesize is %#v, pagenum is %#v", pageparam.Pagesize, pageparam.Pagenum)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}
	// 获取设备数据
	var deviceService services.DeviceService

	totalNum, totalPages, devices, err := deviceService.GetDeviceByNilBusiness(pageparam.Pagesize, pageparam.Pagenum)
	if err != nil {
		d.Response(400, "查找不到设备")
		return
	}

	retList := DeviceList{
		TotalNum:   totalNum,
		TotalPages: totalPages,
		DeviceList: &devices,
	}
	d.Response(200, "获取设备列表成功", retList)
}

// 更新设备的业务ID
func (d *DeviceController) UpdateForBusiness() {
	devUpBusiness := DevUpBusiness{}
	err := json.Unmarshal(d.Ctx.Input.RequestBody, &devUpBusiness)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		d.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&devUpBusiness)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}

	// 获取设备数据
	var deviceService services.DeviceService

	err = deviceService.UpdateForBusiness(devUpBusiness.AssetsNum, devUpBusiness.BusinessId)
	if err != nil {
		d.Response(400, "更新设备失败")
		return
	}

	d.Response(200, "更新设备成功")
}

// 通过业务ID，获取设备列表
func (d *DeviceController) ListForIndex() {
	deviceOfBusiness := DeviceOfBusiness{}
	err := json.Unmarshal(d.Ctx.Input.RequestBody, &deviceOfBusiness)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		d.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&deviceOfBusiness)
	if err != nil {
		// handler error
		d.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		d.Response(400, "输入参数错误")
		return
	}

	// 获取设备数据
	var deviceService services.DeviceService

	devices, err := deviceService.GetDeviceByClass(deviceOfBusiness.BusinessId)
	if err != nil {
		d.Response(400, "查找不到设备")
		return
	}

	d.Response(200, "获取设备列表成功", devices)
}
