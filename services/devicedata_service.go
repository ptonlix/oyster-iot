package services

import (
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/core/logs"
)

type DeviceData struct {
}

// 插入设备上报的数据
func (*DeviceData) Insert(device *models.Device, msg string) error {
	deviceData := models.DeviceData{
		DevId:   device.Id,
		DevType: device.Type,
		Msg:     msg,
	}

	id, err := mysql.Mydb.Insert(&deviceData)
	if err != nil {
		logs.Warn("Device Data insert Failed!", err.Error())
		return err
	}
	logs.Info("Device Data insert Success! id:", id)
	return nil
}
