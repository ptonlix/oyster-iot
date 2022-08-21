package services

import (
	"encoding/json"
	"fmt"
	"oyster-iot/init/constants"
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
type DevSaltDay struct {
	DevName  string          `json:"name"`
	Saltlist [24]interface{} `json:"data"`
}

func (t *SaltService) GetSaltOne(devAssetsNum string) (*SaltOne, error) {
	// 如果InfluxDB在线,则通过InflIxDB读取数据
	if constants.InfluxConf.Online {
		return t.getSaltOneFromInfluxDB(devAssetsNum)
	} else {
		return t.getSaltOneFromSql(devAssetsNum)
	}
}
func (t *SaltService) GetSaltInDay(devAssetsNum []string) (*[]DevSaltDay, error) {
	// 如果InfluxDB在线,则通过InflIxDB读取数据
	if constants.InfluxConf.Online {
		return t.getSaltInDayFromInfluxDB(devAssetsNum)
	} else {
		return t.getSaltInDayFromSql(devAssetsNum)
	}
}

// 获取当前温度传感器的温度
func (t *SaltService) getSaltOneFromSql(devAssetsNum string) (*SaltOne, error) {
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
func (t *SaltService) getSaltInDayFromSql(devAssetsNum []string) (*[]DevSaltDay, error) {
	nowtime := time.Now().Format("2006/01/02")
	var lists []orm.ParamsList
	var alldev []DevSaltDay
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
		one := DevSaltDay{DevName: device.DeviceName, Saltlist: saltlist} //数组是值拷贝
		alldev = append(alldev, one)
	}
	logs.Warn(alldev)
	return &alldev, nil
}

// 获取当前温度传感器的温度
func (t *SaltService) getSaltOneFromInfluxDB(devAssetsNum string) (*SaltOne, error) {
	devDataS := DevdataSevice{}
	devDataS.Host = constants.InfluxConf.Host
	devDataS.Token = constants.InfluxConf.Token
	devDataS.Org = constants.InfluxConf.Org
	devDataS.Bucket = constants.InfluxConf.Bucket
	queryStr := fmt.Sprintf(`
	from(bucket:"oyster") 
		|> range(start: -1d) 
		|> filter(fn: (r) => r._measurement == "salinity" and r.dev_assets_num == "%s" and r._field == "salinity") 
		|> last()`, devAssetsNum)
	logs.Debug("queryStr : ", queryStr)
	data, err := devDataS.GetDataFromInfluxDB(queryStr)
	if err != nil {
		logs.Warn("Get Dev:%#v  salinity Failed! err:%#v", devAssetsNum, err)
		return nil, err
	}

	var salt SaltOne
	for i := 0; i < len(*data); i++ {
		if value, ok := (*data)[i]["_value"].(float64); ok {
			salt.Salt = float32(value)
		} else {
			logs.Warn("Get Dev:%#v  salinity Failed! err: interface error", devAssetsNum)
		}
		if value, ok := (*data)[i]["_time"].(time.Time); ok {
			salt.Time = value.Local().Format("2006-01-02 15:04:05")
		} else {
			logs.Warn("Get Dev:%#v  salinity Failed! err: interface error", devAssetsNum)
		}
	}
	return &salt, nil
}

func (t *SaltService) getSaltInDayFromInfluxDB(devAssetsNum []string) (*[]DevSaltDay, error) {
	devDataS := DevdataSevice{}
	devDataS.Host = constants.InfluxConf.Host
	devDataS.Token = constants.InfluxConf.Token
	devDataS.Org = constants.InfluxConf.Org
	devDataS.Bucket = constants.InfluxConf.Bucket

	stoptime := time.Now()
	//logs.Debug(stoptime.UTC().Format(time.RFC3339))
	startTime := time.Date(stoptime.Year(), stoptime.Month(), stoptime.Day(), 0, 0, 0, 0, stoptime.Location())

	var alldev []DevSaltDay

	for _, v := range devAssetsNum {
		queryStr := fmt.Sprintf(`
		from(bucket:"oyster") 
			|> range(start:%s, stop:%s) 
			|> filter(fn: (r) => r._measurement == "salinity" and r.dev_assets_num == "%s" and r._field == "salinity") 
			`, startTime.UTC().Format(time.RFC3339), stoptime.UTC().Format(time.RFC3339), v)
		logs.Debug("queryStr : ", queryStr)
		data, err := devDataS.GetDataFromInfluxDB(queryStr)
		if err != nil {
			logs.Warn("Get Dev:%#v  salinity Failed! err:%#v", devAssetsNum, err)
			return nil, err
		}

		templist := [24]interface{}{} //每日温度记录列表

		for i := 0; i < len(*data); i++ {
			saltValue := 0.
			if value, ok := (*data)[i]["_value"].(float64); ok {
				saltValue = value
			} else {
				logs.Warn("Get Dev:%#v  salinity Failed! err: interface error", devAssetsNum)
			}
			if value, ok := (*data)[i]["_time"].(time.Time); ok {
				templist[value.Local().Hour()] = saltValue
			} else {
				logs.Warn("Get Dev:%#v  salinity Failed! err: interface error", devAssetsNum)
			}
		}

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

		one := DevSaltDay{DevName: device.DeviceName, Saltlist: templist} //数组是值拷贝
		alldev = append(alldev, one)

	}
	return &alldev, nil
}
