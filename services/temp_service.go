package services

import (
	"encoding/json"
	"oyster-iot/devaccess/modules/mqtt"
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type TempService struct {
}

type TempOne struct {
	Temp int    `json:"temperature"`
	Time string `json:"ts"`
}
type TempInDay struct {
	Temps []TempOne `json:"temp_in_day"`
}

type TempMsg struct {
	Temperature int `json:temperature`
}

type TempCmd struct {
	Token string      `json:"token,omitempty"`
	Cmd   string      `json:"cmd"`
	Data  interface{} `json:"data,omitempty"`
}

// 下发操作指令，发送命令到硬件端执行
func (t *TempService) OperateCmd(tcmd *TempCmd) error {
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

// 获取当前温度传感器的温度
func (t *TempService) GetTempOne(devAssetsNum string) (*TempOne, error) {
	tempData := models.DeviceData{}
	err := mysql.Mydb.Raw("SELECT msg, ts FROM device_data WHERE dev_assets_num = ? ORDER BY ts  desc limit 1;", devAssetsNum).QueryRow(&tempData)
	if err == orm.ErrNoRows {
		logs.Warn("Get Dev:%#v  temperature ErrNoRows!", devAssetsNum)
		return nil, nil
	} else if err != nil {
		logs.Warn("Get Dev:%#v  temperature Failed! err:%#v", devAssetsNum, err)
		return nil, err
	}

	var msg TempMsg
	var temp TempOne
	if err := json.Unmarshal([]byte(tempData.Msg), &msg); err != nil {
		logs.Warn("Unmarshal Msg temperature Failed!")
		return nil, err
	}
	temp.Temp = msg.Temperature
	temp.Time = tempData.Ts.Format("2006-01-02 15:04:05")

	return &temp, err
}

// 获取近24小时的温度情况
func (t *TempService) GetTempInDay(devAssetsNum string) (*TempInDay, error) {
	var tempDay TempInDay
	var lists []orm.ParamsList
	num, err := mysql.Mydb.Raw("SELECT msg, ts FROM device_data WHERE dev_assets_num = ? && ts >= (NOW() - interval 24 hour);", devAssetsNum).ValuesList(&lists)
	if err != nil {
		logs.Warn(err)
		return nil, err
	}
	logs.Info("Get Dev:%#v  temperature number:%#v successful!", devAssetsNum, num)
	// 将数据json信息转成温度列表信息
	for _, value := range lists {
		var msg TempMsg
		var temp TempOne
		err := json.Unmarshal([]byte(value[0].(string)), &msg)
		if err != nil {
			logs.Warn("Unmarshal Msg temperature Failed!")
			continue
		}
		temp.Temp = msg.Temperature
		temp.Time = value[1].(string)

		tempDay.Temps = append(tempDay.Temps, temp)
	}
	logs.Warn(tempDay.Temps)
	return &tempDay, err
}
