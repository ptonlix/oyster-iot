package services

import (
	"encoding/json"
	"oyster-iot/models"

	"github.com/beego/beego/v2/core/logs"
)

type InfluxConfig struct {
	Host    string
	Token   string
	Org     string
	Bucket  string
	Online  bool
	InfluxS *InfluxServer
}
type DevdataSevice struct {
	InfluxConfig
}

type mqttPayload struct {
	AssetsNum string                 `json:"assets_num"`
	Token     string                 `json:"token"`
	Msg       map[string]interface{} `json:"msg"`
}

func (d *DevdataSevice) MQTTMsgProc(msgbody []byte) (err error) {

	payload := &mqttPayload{}
	if err = json.Unmarshal(msgbody, payload); err != nil {
		logs.Error("Msg Consumer: Cannot unmarshal msg payload to JSON:", err)
		return
	}
	if len(payload.Token) == 0 {
		logs.Warn("Msg Consumer: Payload token missing")
		return ErrMQTTToken
	}
	if len(payload.Msg) == 0 {
		logs.Warn("Msg Consumer: Payload token missing")
		return ErrMQTTMsg
	}
	logs.Info("Token is %s, Msg is %#v", payload.Token, payload.Msg)

	//1.查询设备表，判断Token的合法性
	var deviceService DeviceService
	var device *models.Device
	if device, err = deviceService.GetDeviceByTokenID(payload.Token, payload.AssetsNum); err != nil {
		logs.Warn("Msg Consumer: Cannot find device!")
		return err
	}

	//2.插入设备上报的信息, influxDB不在线才插入到mysql
	if !d.Online {
		var deviceData DeviceData
		childJson, _ := json.Marshal(payload.Msg)
		childString := string(childJson)
		deviceData.Insert(device, childString)
	}
	//3.插入influxDB
	switch device.Type {
	case TEMPDEVICE:
		d.InsertTempDataToInfluxDB(device, payload.Msg)
	case SALINITYDEVICE:
		d.InsertSaltDataToInfluxDB(device, payload.Msg)
	}

	return nil
}

// 设备上报的数据插入influxDB
func (d *DevdataSevice) InsertTempDataToInfluxDB(device *models.Device, fields map[string]interface{}) error {
	d.InfluxS = NewInfluxService(d.Host, d.Token, d.Org, d.Bucket)
	defer d.InfluxS.Close()

	//生成Tag数据
	tags := map[string]string{"dev_assets_num": device.AssetsNum}

	if err := d.InfluxS.WriteData("temperature", tags, fields); err != nil {
		logs.Warn("Insert Temperature Data To InfluxDB, Failed : ", err.Error())
		return err
	}
	return nil
}

// 设备上报的数据插入influxDB
func (d *DevdataSevice) InsertSaltDataToInfluxDB(device *models.Device, fields map[string]interface{}) error {
	d.InfluxS = NewInfluxService(d.Host, d.Token, d.Org, d.Bucket)
	defer d.InfluxS.Close()
	//生成Tag数据
	tags := map[string]string{"dev_assets_num": device.AssetsNum}

	if err := d.InfluxS.WriteData("salinity", tags, fields); err != nil {
		logs.Warn("Insert Salinity Data To InfluxDB, Failed : ", err.Error())
		return err
	}
	return nil
}

/*
query: Flux语法字符串
`from(bucket:"oyster")|> range(start: -1h) |> filter(fn: (r) => r._measurement == "temperature")`
*/
func (d *DevdataSevice) GetDataFromInfluxDB(query string) (*[]map[string]interface{}, error) {
	d.InfluxS = NewInfluxService(d.Host, d.Token, d.Org, d.Bucket)
	defer d.InfluxS.Close()

	data, err := d.InfluxS.ReadData(query)
	if err != nil {
		logs.Warn("Insert Salinity Data To InfluxDB, Failed : ", err.Error())
		return nil, err
	}
	return data, nil
}
