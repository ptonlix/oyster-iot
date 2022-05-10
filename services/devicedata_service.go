package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"oyster-iot/init/mysql"
	"oyster-iot/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

func init() {
	//生成测试数据
TestCreateTempData()
}

type DeviceData struct {
}

// 插入设备上报的数据
func (*DeviceData) Insert(device *models.Device, msg string) error {
	deviceData := models.DeviceData{
		DevAssetsNum: device.AssetsNum,
		DevType:      device.Type,
		Msg:          msg,
	}

	id, err := mysql.Mydb.Insert(&deviceData)
	if err != nil {
		logs.Warn("Device Data insert Failed!", err.Error())
		return err
	}
	logs.Info("Device Data insert Success! id:", id)
	return nil
}

func TestCreateTempData() {

	ti := time.Now()
	addTime := time.Date(ti.Year(), ti.Month(), ti.Day(), 0, 0, 0, 0, ti.Location())
	type temp struct {
		Temperature float64 `json:"temperature"`
	}
	//插入温度数据
	deviceDatas := []models.DeviceData{
		{
			DevAssetsNum: "6a32e884-0fd9-955f-5739-ebb78f5ba685",
			DevType:      TEMPDEVICE,
		},
		{
			DevAssetsNum: "ed1406cf-9ffa-700e-df31-2124ca90157e",
			DevType:      TEMPDEVICE,
		},
		{
			DevAssetsNum: "cc86f1ba-001a-e27a-ba18-21336e9d3f62",
			DevType:      TEMPDEVICE,
		},
	}

	for i := 0; i < 24; i++ {
		for _, v := range deviceDatas {
			rand.Seed(time.Now().UnixNano())
			value, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", 10+rand.Float64()*20), 64)

			tmp := temp{value}

			tempjson, _ := json.Marshal(tmp)

			v.Ts = addTime.Add(time.Duration(i) * time.Hour)
			v.Msg = string(tempjson)

			id, err := mysql.Mydb.Insert(&v)
			if err != nil {
				logs.Warn("Device Data insert Failed!", err.Error())
			}
			logs.Info("Device Data insert Success! id:", id)
		}
	}
}
