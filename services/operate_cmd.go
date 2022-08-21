package services

import (
	"encoding/json"
	"oyster-iot/devaccess/modules/mqtt"

	"github.com/beego/beego/v2/core/logs"
)

const (
	RefreshTemp = "refresh_temp"
	RefreshSalt = "refresh_salt"
)

type DevCmd struct {
	AssetsNum string      `json:"assets_num"`
	Token     string      `json:"token,omitempty"`
	Cmd       string      `json:"cmd"`
	Data      interface{} `json:"data,omitempty"`
}

// 下发操作指令，发送命令到硬件端执行
func (t *DevCmd) OperateCmd(tcmd *DevCmd) error {
	tj, err := json.Marshal(tcmd)
	if err != nil {
		logs.Warn("Operate CMD Send Failed! err: ", err.Error())
		return err
	}
	if err := mqtt.Send(tj); err != nil {
		logs.Warn("Operate CMD Send Failed! err: ", err.Error())
	}
	return err
}
