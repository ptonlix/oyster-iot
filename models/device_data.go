package models

import "time"

type DeviceData struct {
	Id      int       `json:"id"       orm:"pk;auto"`
	DevId   int       `json:"dev_id"   orm:"index;description(设备ID)"`
	DevType string    `json:"dev_type" orm:"size(64)"`
	Msg     string    `json:"msg"      orm:"type(text);description(设备上报的信息)"`
	Ts      time.Time `json:"ts"       orm:"auto_now_add;type(datetime);precision(4);description(保存信息时间戳)"`
}
