package controllers

import (
	"encoding/json"
	"errors"
	"oyster-iot/models"
	"oyster-iot/services"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type SaltController struct {
	BaseController
}

func (t *SaltController) validAssetsInfo() (device *models.Device, err error) {
	assetsInfo := DevAssetsInfo{}
	err = json.Unmarshal(t.Ctx.Input.RequestBody, &assetsInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		t.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&assetsInfo)
	if err != nil {
		// handler error
		t.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		t.Response(400, "输入参数错误")
		return
	}

	// 通过资产编码获取温度信息
	var deviceService services.DeviceService

	if device, err = deviceService.GetDeviceByAssetsNum(assetsInfo.AssetsNum); err != nil {
		t.Response(500, "查找不到该设备")
		return
	}

	// 比较设备类型是否一致
	if device.Type != assetsInfo.DeviceType {
		t.Response(400, "设备类型不正确")
		err = errors.New("DevTypeError")
		return
	}
	return
}

// 获取最新的设备温度信息
func (t *SaltController) GetSalt() {
	device, err := t.validAssetsInfo()
	if err != nil {
		return
	}
	// 获取设备数据信息
	var saltService services.SaltService
	saltData, err := saltService.GetSaltOne(device.AssetsNum)
	if err != nil {
		t.Response(500, "获取盐度数据失败")
		return
	}

	// 数据转换展示给前端
	t.Response(200, "获取盐度数据成功", saltData)
}

// 获取近24小时的设备的温度信息
func (t *SaltController) GetSaltInDay() {
	assetsInfo := SomeDevAssets{}
	err := json.Unmarshal(t.Ctx.Input.RequestBody, &assetsInfo)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		t.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&assetsInfo)
	if err != nil {
		// handler error
		t.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		t.Response(400, "输入参数错误")
		return
	}

	// 获取设备数据信息
	var saltService services.SaltService
	saltData, err := saltService.GetSaltInDay(assetsInfo.AssetsNum)
	if err != nil {
		logs.Warn(err.Error())
		t.Response(500, "获取盐度数据失败")
		return
	}

	// 数据转换展示给前端
	t.Response(200, "获取盐度数据成功", saltData)
}

// 发送消息获取设备温度
func (t *SaltController) SendSaltCmd() {
	device, err := t.validAssetsInfo()
	if err != nil {
		return
	}

	cmd := &services.DevCmd{AssetsNum: device.AssetsNum, Token: device.Token, Cmd: services.RefreshSalt}
	var devCmd services.DevCmd
	if err := devCmd.OperateCmd(cmd); err != nil {
		t.Response(500, "发送操作设备命令失败")
		return
	}

	t.Response(200, "发送操作命令成功")
}
