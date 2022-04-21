package services

import (
	"log"
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
)

type DeviceService struct {
}

// 通过Token 获取设备信息
func (*DeviceService) GetDeviceByTokenID(token string) (*models.Device, error) {
	device := models.Device{Token: token}

	err := mysql.Mydb.Read(&device, "Token")

	if err == orm.ErrNoRows {
		log.Printf("Token %s: Cannot find device!\n", token)
	} else if err != nil {
		log.Println(err)
	}

	return &device, err
}
