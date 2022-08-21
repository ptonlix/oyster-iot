package services

import (
	"encoding/json"
	"fmt"
	"oyster-iot/init/mysql"
	"oyster-iot/models"
	"strconv"
	"time"

	"oyster-iot/init/constants"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type TempService struct {
}

type TempOne struct {
	Temp float32 `json:"temperature"`
	Time string  `json:"ts"`
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
	// 如果InfluxDB在线,则通过InflIxDB读取数据
	if constants.InfluxConf.Online {
		return t.getTempOneFromInfluxDB(devAssetsNum)
	} else {
		return t.getTempOneFromSql(devAssetsNum)
	}
}

// 获取近24小时的温度情况
func (t *TempService) GetTempInDay(devAssetsNum []string) (*[]DevTempDay, error) {
	// 如果InfluxDB在线,则通过InflIxDB读取数据
	if constants.InfluxConf.Online {
		return t.getTempInDayFromInfluxDB(devAssetsNum)
	} else {
		return t.getTempInDayFromSql(devAssetsNum)
	}
}

// 从mysql获取当前温度传感器的温度
func (t *TempService) getTempOneFromSql(devAssetsNum string) (*TempOne, error) {
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

// 从mysql获取近24小时的温度情况
func (t *TempService) getTempInDayFromSql(devAssetsNum []string) (*[]DevTempDay, error) {
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

// 从influxDB获取最新的温度数据
func (t *TempService) getTempOneFromInfluxDB(devAssetsNum string) (*TempOne, error) {
	devDataS := DevdataSevice{}
	devDataS.Host = constants.InfluxConf.Host
	devDataS.Token = constants.InfluxConf.Token
	devDataS.Org = constants.InfluxConf.Org
	devDataS.Bucket = constants.InfluxConf.Bucket
	queryStr := fmt.Sprintf(`
	from(bucket:"oyster") 
		|> range(start: -1d) 
		|> filter(fn: (r) => r._measurement == "temperature" and r.dev_assets_num == "%s" and r._field == "temperature") 
		|> last()`, devAssetsNum)
	logs.Debug("queryStr : ", queryStr)
	data, err := devDataS.GetDataFromInfluxDB(queryStr)
	if err != nil {
		logs.Warn("Get Dev:%#v  temperature Failed! err:%#v", devAssetsNum, err)
		return nil, err
	}

	var temp TempOne
	for i := 0; i < len(*data); i++ {
		if value, ok := (*data)[i]["_value"].(float64); ok {
			temp.Temp = float32(value)
		} else {
			logs.Warn("Get Dev:%#v  temperature Failed! err: interface error", devAssetsNum)
		}
		if value, ok := (*data)[i]["_time"].(time.Time); ok {
			temp.Time = value.Local().Format("2006-01-02 15:04:05")
		} else {
			logs.Warn("Get Dev:%#v  temperature Failed! err: interface error", devAssetsNum)
		}
	}
	return &temp, nil
}

// 从influxDB获取今天的温度数据
func (t *TempService) getTempInDayFromInfluxDB(devAssetsNum []string) (*[]DevTempDay, error) {
	devDataS := DevdataSevice{}
	devDataS.Host = constants.InfluxConf.Host
	devDataS.Token = constants.InfluxConf.Token
	devDataS.Org = constants.InfluxConf.Org
	devDataS.Bucket = constants.InfluxConf.Bucket

	stoptime := time.Now()
	//logs.Debug(stoptime.UTC().Format(time.RFC3339))
	startTime := time.Date(stoptime.Year(), stoptime.Month(), stoptime.Day(), 0, 0, 0, 0, stoptime.Location())

	var alldev []DevTempDay

	for _, v := range devAssetsNum {
		queryStr := fmt.Sprintf(`
		from(bucket:"oyster") 
			|> range(start:%s, stop:%s) 
			|> filter(fn: (r) => r._measurement == "temperature" and r.dev_assets_num == "%s" and r._field == "temperature") 
			`, startTime.UTC().Format(time.RFC3339), stoptime.UTC().Format(time.RFC3339), v)
		logs.Debug("queryStr : ", queryStr)
		data, err := devDataS.GetDataFromInfluxDB(queryStr)
		if err != nil {
			logs.Warn("Get Dev:%#v  temperature Failed! err:%#v", devAssetsNum, err)
			return nil, err
		}

		templist := [24]interface{}{} //每日温度记录列表

		for i := 0; i < len(*data); i++ {
			tempValue := 0.
			if value, ok := (*data)[i]["_value"].(float64); ok {
				tempValue = value
			} else {
				logs.Warn("Get Dev:%#v  temperature Failed! err: interface error", devAssetsNum)
			}
			if value, ok := (*data)[i]["_time"].(time.Time); ok {
				templist[value.Local().Hour()] = tempValue
			} else {
				logs.Warn("Get Dev:%#v  temperature Failed! err: interface error", devAssetsNum)
			}
		}

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

		one := DevTempDay{DevName: device.DeviceName, Templist: templist} //数组是值拷贝
		alldev = append(alldev, one)

	}
	return &alldev, nil
}
