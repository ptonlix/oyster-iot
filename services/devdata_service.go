package services

import (
	"encoding/json"
	"log"
	"oyster-iot/models"
)

type DevdataSevice struct {
}

type mqttPayload struct {
	Token string                 `json:"token"`
	Msg   map[string]interface{} `json:"msg"`
}

func (d *DevdataSevice) MQTTMsgProc(msgbody []byte) (err error) {

	payload := &mqttPayload{}
	if err = json.Unmarshal(msgbody, payload); err != nil {
		log.Println("Msg Consumer: Cannot unmarshal msg payload to JSON:", err)
		return
	}
	if len(payload.Token) == 0 {
		log.Println("Msg Consumer: Payload token missing")
		return ErrMQTTToken
	}
	if len(payload.Msg) == 0 {
		log.Println("Msg Consumer: Payload token missing")
		return ErrMQTTMsg
	}
	log.Printf("Token is %s, Msg is %v", payload.Token, payload.Msg)

	//1.查询设备表，判断Token的合法性
	var deviceService DeviceService
	var device *models.Device
	if device, err = deviceService.GetDeviceByTokenID(payload.Token); err != nil {
		log.Println("Msg Consumer: Cannot find device!\n")
		return err
	}

	//2.插入设备上报的信息
	var deviceData DeviceData
	childJson, _ := json.Marshal(payload.Msg)
	childString := string(childJson)
	deviceData.Insert(device, childString)

	return nil
}
