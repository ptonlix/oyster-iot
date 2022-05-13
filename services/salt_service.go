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

type SaltService struct {
}

type SaltOne struct {
	Salt float32 `json:"salinity"`
	Time string  `json:"ts"`
}

type SaltMsg struct {
	Salinity float32 `json:salinity`
}
type SaltTempDay struct {
	DevName  string          `json:"name"`
	Saltlist [24]interface{} `json:"data"`
}

// 获取当前温度传感器的温度
func (t *SaltService) GetSaltOne(devAssetsNum string) (*SaltOne, error) {
	saltData := models.DeviceData{}
	err := mysql.Mydb.Raw("SELECT msg, ts FROM device_data WHERE dev_assets_num = ? ORDER BY ts  desc limit 1;", devAssetsNum).QueryRow(&saltData)
	if err == orm.ErrNoRows {
		logs.Warn("Get Dev:%#v  salinity ErrNoRows!", devAssetsNum)
		return nil, nil
	} else if err != nil {
		logs.Warn("Get Dev:%#v  salinity Failed! err:%#v", devAssetsNum, err)
		return nil, err
	}

	var msg SaltMsg
	var salt SaltOne
	if err := json.Unmarshal([]byte(saltData.Msg), &msg); err != nil {
		logs.Warn("Unmarshal Msg salinity Failed!")
		return nil, err
	}
	salt.Salt = msg.Salinity
	salt.Time = saltData.Ts.Format("2006-01-02 15:04:05")

	return &salt, err
}

// 获取近24小时的盐度情况
func (t *SaltService) GetSaltInDay(devAssetsNum []string) (*[]SaltTempDay, error) {
	nowtime := time.Now().Format("2006/01/02")
	var lists []orm.ParamsList
	var alldev []SaltTempDay
	for _, v := range devAssetsNum {
		//查找该设备
		var deviceService *DeviceService
		device, err := deviceService.GetDeviceByAssetsNum(v)
		if err != nil {
			return nil, ErrDevNoFound
		}
		// 比较设备类型是否一致
		if device.Type != SALINITYDEVICE {
			return nil, ErrDevType
		}

		saltlist := [24]interface{}{} //每日盐度记录列表

		// 获取一台设备今天的温度数据
		num, err := mysql.Mydb.Raw("SELECT msg, HOUR(ts) FROM device_data WHERE dev_assets_num = ? && DATE_FORMAT(ts,'%Y/%m/%d')= ? ;", v, nowtime).ValuesList(&lists)
		if err != nil {
			logs.Warn(err)
			return nil, err
		}
		logs.Info("Get Dev:%#v  salinity number:%#v successful!", v, num)
		// 查找设备名称

		// 将数据json信息转成盐度列表信息
		for _, value := range lists {
			var msg SaltMsg

			err := json.Unmarshal([]byte(value[0].(string)), &msg)
			if err != nil {
				logs.Warn("Unmarshal Msg salinity Failed!")
				continue
			}
			t, _ := strconv.Atoi(value[1].(string))
			saltlist[t] = msg.Salinity
		}
		one := SaltTempDay{DevName: device.DeviceName, Saltlist: saltlist} //数组是值拷贝
		alldev = append(alldev, one)
	}
	logs.Warn(alldev)
	return &alldev, nil
}
