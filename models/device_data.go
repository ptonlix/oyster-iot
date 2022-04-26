package models

import "time"

type DeviceData struct {
	Id           int       `json:"id"       orm:"pk;auto"`
	DevAssetsNum string    `json:"dev_assets_num"   orm:"index;description(设备资产编码)"`
	DevType      string    `json:"dev_type" orm:"size(64)"`
	Msg          string    `json:"msg"      orm:"type(text);description(设备上报的信息)"`
	Ts           time.Time `json:"ts"       orm:"auto_now_add;type(datetime);description(保存信息时间戳)"`
}
