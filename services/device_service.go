package services

import (
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/core/logs"

	"github.com/beego/beego/v2/client/orm"
)

//支持的设备类型
const (
	TEMPDEVICE     = "Temp"     //温度探测设备
	SALINITYDEVICE = "Salinity" //盐度探测设备
	CAMERADEVICE   = "Camera"   //相机监控设备
)

type DevicesOfBusiness struct {
	Temp        []*models.Device `json:"temp"`
	TempNum     int64            `json:"temp_num"`
	Salinity    []*models.Device `json:"salinity"`
	SalinityNum int64            `json:"salinity_num"`
}
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

// 修改设备
func (*DeviceService) Update(device *models.Device) error {

	id, err := mysql.Mydb.Update(device)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Update Device successful! ID:", id)
	return err
}

// 删除设备
func (*DeviceService) Delete(device *models.Device) (err error) {

	// 删除设备, 同时删除设备上报的数据
	to, err := mysql.Mydb.Begin()
	if err != nil {
		logs.Error("start the transaction failed")
		return err
	}
	defer func() {
		if err != nil {
			err = to.Rollback()
			if err != nil {
				logs.Error("roll back transaction failed", err)
			}
		} else {
			err = to.Commit()
			if err != nil {
				logs.Error("commit transaction failed.", err)
			}
		}
	}()
	//删除设备上报的数据

	deviceData := models.DeviceData{
		DevAssetsNum: device.AssetsNum,
		DevType:      device.Type,
	}
	qs := to.QueryTable(deviceData).Filter("dev_assets_num", device.AssetsNum)
	delid, err := qs.Delete()
	if err != nil {
		logs.Warn(err)
		return
	}
	logs.Info("Delete Device Data num: ", delid)

	id, err := to.Delete(device)
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
		return nil, err
	} else if err != nil {
		logs.Warn(err)
	}

	return &device, err
}

// 获取全部设备
func (*DeviceService) GetDevicesByPage(pageSize, pageNum int) (int, int, []*models.Device, error) {

	var devices []*models.Device
	qs := mysql.Mydb.QueryTable(&models.Device{})

	totalRecord, err := qs.Count()
	if err != nil {
		logs.Warn(err)

	}

	num, err := qs.Limit(pageSize, pageSize*(pageNum-1)).All(&devices)

	if err != nil {
		logs.Warn(err)
	}

	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	logs.Info("Get Devices successful! Totalcount: %v TotalPages: %v Returned Rows Num: %#v", totalRecord, totalPageNum, num)

	return int(totalRecord), totalPageNum, devices, err
}

// 通过业务ID获取设备
func (*DeviceService) GetDeviceByBusiness(businessId int) ([]*models.Device, error) {
	var devices []*models.Device
	qs := mysql.Mydb.QueryTable(&models.Device{})

	num, err := qs.Filter("business_id", businessId).All(&devices)
	if err == orm.ErrNoRows {
		logs.Warn("businessId %#v: Cannot find device!\n", businessId)
		return nil, err
	} else if err != nil {
		logs.Warn(err)
	}
	logs.Info("Get Devices successful! Returned Rows Num: %#v", num)
	return devices, err
}

// 通过业务ID获取设备并分类
func (*DeviceService) GetDeviceByClass(businessId int) (*DevicesOfBusiness, error) {
	var devs DevicesOfBusiness
	//var devices []*models.Device
	qs := mysql.Mydb.QueryTable(&models.Device{})
	//获取温度传感器
	num, err := qs.Filter("business_id", businessId).Filter("type", "Temp").All(&devs.Temp)
	if err == orm.ErrNoRows {
		logs.Warn("businessId %#v: Cannot find device!\n", businessId)
		return nil, err
	} else if err != nil {
		logs.Warn(err)
	}
	devs.TempNum = num
	logs.Info("Get Temp Devices successful! Returned Rows Num: %#v", num)

	num, err = qs.Filter("business_id", businessId).Filter("type", "Salinity").All(&devs.Salinity)
	if err == orm.ErrNoRows {
		logs.Warn("businessId %#v: Cannot find device!\n", businessId)
		return nil, err
	} else if err != nil {
		logs.Warn(err)
	}
	devs.SalinityNum = num
	logs.Info("Get Salinity Devices successful! Returned Rows Num: %#v", num)
	return &devs, err
}

// 通过未绑定业务的设备列表
func (*DeviceService) GetDeviceByNilBusiness(pageSize, pageNum int) (int, int, []*models.Device, error) {
	var devices []*models.Device
	qs := mysql.Mydb.QueryTable(&models.Device{})

	totalRecord, err := qs.Filter("business_id", 0).Count()
	if err != nil {
		logs.Warn(err)

	}

	num, err := qs.Filter("business_id", 0).Limit(pageSize, pageSize*(pageNum-1)).All(&devices)
	if err != nil {
		logs.Warn(err)
	}

	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	logs.Info("Get Devices successful! Totalcount: %v TotalPages: %v Returned Rows Num: %#v", totalRecord, totalPageNum, num)

	return int(totalRecord), totalPageNum, devices, err
}

// 修改设备关联的业务ID
func (*DeviceService) UpdateForBusiness(assertsNum []string, businessId int) (err error) {
	// 开启事务
	to, err := mysql.Mydb.Begin()
	if err != nil {
		logs.Error("start the transaction failed")
		return

	}

	defer func() {
		if err != nil {
			err = to.Rollback()
			if err != nil {
				logs.Error("roll back transaction failed", err)
			}
		} else {
			err = to.Commit()
			if err != nil {
				logs.Error("commit transaction failed.", err)
			}
		}
	}()

	for _, v := range assertsNum {
		device := models.Device{AssetsNum: v}
		err = to.Read(&device, "AssetsNum")
		if err != nil {
			logs.Error("execute transaction's select sql fail, rollback.", err)
			return
		}
		// 更新业务字段
		device.BusinessId = businessId
		_, err = to.Update(&device, "BusinessId")
		if err != nil {
			logs.Error("execute transaction's update sql fail, rollback.", err)
			return
		}
	}
	return nil
}
