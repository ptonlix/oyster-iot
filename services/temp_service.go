package services

import (
	"encoding/json"
	"oyster-iot/init/mysql"
	"oyster-iot/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type TempService struct {
}

type TempOne struct {
	Temp float32 `json:"temperature"`
	Time string  `json:"ts"`
}
type TempInDay struct {
	Temps     []TempOne `json:"temp_in_day"`
	DevAssert string    `json:"assets_num"`
}

type AllTempInDay struct {
	DevTemp []TempInDay `json:"dev_temp"`
}

type TempMsg struct {
	Temperature float32 `json:temperature`
}

type DevTempDay struct {
	DevName  string          `json:"name"`
	Templist [24]interface{} `json:"data"`
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
func (t *TempService) GetTempInDay(devAssetsNum []string) (*[]DevTempDay, error) {
	nowtime := time.Now().Format("2006/01/02")
	var lists []orm.ParamsList
	var alldev []DevTempDay
	for _, v := range devAssetsNum {
		//查找该设备
		var deviceService *DeviceService
		device, err := deviceService.GetDeviceByAssetsNum(v)
		if err != nil {
			return nil, ErrDevNoFound
		}
		// 比较设备类型是否一致
		if device.Type != TEMPDEVICE {
			return nil, ErrDevType
		}

		templist := [24]interface{}{} //每日温度记录列表

		// 获取一台设备今天的温度数据
		num, err := mysql.Mydb.Raw("SELECT msg, HOUR(ts) FROM device_data WHERE dev_assets_num = ? && DATE_FORMAT(ts,'%Y/%m/%d')= ? ;", v, nowtime).ValuesList(&lists)
		if err != nil {
			logs.Warn(err)
			return nil, err
		}
		logs.Info("Get Dev:%#v  temperature number:%#v successful!", v, num)
		// 查找设备名称

		// 将数据json信息转成温度列表信息
		for _, value := range lists {
			var msg TempMsg

			err := json.Unmarshal([]byte(value[0].(string)), &msg)
			if err != nil {
				logs.Warn("Unmarshal Msg temperature Failed!")
				continue
			}
			t, _ := strconv.Atoi(value[1].(string))
			templist[t] = msg.Temperature
		}
		one := DevTempDay{DevName: device.DeviceName, Templist: templist} //数组是值拷贝
		alldev = append(alldev, one)
	}
	logs.Warn(alldev)
	return &alldev, nil
}
