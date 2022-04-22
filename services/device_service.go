package services

import (
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/core/logs"

	"github.com/beego/beego/v2/client/orm"
)

type DeviceService struct {
}

// 通过Token 获取设备信息
func (*DeviceService) GetDeviceByTokenID(token string) (*models.Device, error) {
	device := models.Device{Token: token}

	err := mysql.Mydb.Read(&device, "Token")

	if err == orm.ErrNoRows {
		logs.Warn("Token %s: Cannot find device!", token)
	} else if err != nil {
		logs.Warn(err)
	}

	return &device, err
}

// 添加设备
func (*DeviceService) Add(device *models.Device) error {

	id, err := mysql.Mydb.Insert(device)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Insert Device successful! ID:", id)
	return err
}

// 添加设备
func (*DeviceService) Update(device *models.Device) error {

	id, err := mysql.Mydb.Update(device)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Update Device successful! ID:", id)
	return err
}

// 添加设备
func (*DeviceService) Delete(device *models.Device) error {

	id, err := mysql.Mydb.Delete(device)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Delete Device successful! ID:", id)
	return err
}

// 通过资产编码查找设备
func (*DeviceService) GetDeviceByAssetsNum(assetsNum string) (*models.Device, error) {
	device := models.Device{AssetsNum: assetsNum}

	err := mysql.Mydb.Read(&device, "AssetsNum")

	if err == orm.ErrNoRows {
		logs.Warn("AssetsNum %#v: Cannot find device!\n", assetsNum)
	} else if err != nil {
		logs.Warn(err)
	}

	return &device, err
}

// 获取全部设备
func (*DeviceService) GetDevicesByPage(pageSize, pageNum int) ([]*models.Device, error) {

	var devices []*models.Device
	qs := mysql.Mydb.QueryTable(&models.Device{})

	num, err := qs.Limit(pageSize, pageSize*(pageNum-1)).All(&devices)

	if err != nil {
		logs.Warn(err)
	}

	logs.Info("Get Devices successful! Returned Rows Num: %#v", num)

	return devices, err
}
