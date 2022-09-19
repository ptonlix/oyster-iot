package models

import "time"

type VideoSpace struct {
	Id                  int       `json:"id"                    orm:"auto;description(设备ID)"`
	UserId              int       `json:"user_id"               orm:"null;index;description(用户ID)"`
	Username            string    `json:"username"              orm:"null;description(用户名称)"`
	Platform            string    `json:"platform"              orm:"null;description(流媒体平台)"`
	SpaceId             string    `json:"space_id"              orm:"index;description(空间ID)"`
	SpaceName           string    `json:"space_name"            orm:"description(空间名称)"`
	AccessType          string    `json:"access_type"           orm:"description(接入类型)"`
	Desc                string    `json:"desc"                  orm:"type(text);description(空间描述)"`
	Disabled            bool      `json:"disabled"              orm:"description(流是否被启用,mfalse:启用,true:禁用)"`
	SpaceCreatedAt      int64     `json:"space_create_at"       orm:"description(空间创建时间)"`
	SpaceUpdatedAt      int64     `json:"space_update_at"       orm:"description(空间修改时间)"`
	DevicesCount        int64     `json:"devices_count"         orm:"description(空间设备数量)"`
	StreamCount         int64     `json:"stream_count"          orm:"description(空间视频流数量)"`
	OnlineStreamCount   int64     `json:"online_stream_count"   orm:"description(空间在线视频流数量)"`
	DisabledStreamCount int64     `json:"disabled_stream_count" orm:"description(空间禁用视频流数量)"`
	OnDemandPull        bool      `json:"on_demand_pull"        orm:"description(按需拉流开关，默认关闭)"`
	Createdat           time.Time `json:"created_at"            orm:"auto_now_add;type(datetime);description(创建时间)"`
	Updatedat           time.Time `json:"updated_at"            orm:"auto_now;type(datetime);description(更新时间)"`
}

type VideoDevice struct {
	Id           int    `json:"id"             orm:"auto;description(设备ID)"`
	DeviceName   string `json:"device_name"    orm:"description(设备名称)"`
	CloudId      string `json:"cloud_id"       orm:"description(云端设备ID设备国标ID)"`
	Type         string `json:"type"           orm:"size(64);description(设备类型:平台 IPC)"`
	BusinessId   int    `json:"busiUess_id"    orm:"null;index;description(设备关联的业务ID)"`
	PlatformCode string `json:"platform_code"  orm:"size(64);description(平台类型:腾讯或七牛)"`
	Username     string `json:"username"       orm:"size(64);description(摄像头登录用户名,一般与国标ID一致)"`
	Password     string `json:"password"       orm:"size(64);description(摄像头登录密码)"`
	Desc         string `json:"desc"           orm:"type(text);description(设备描述)"`
}

type VideoSpaceAndUser struct {
	SpaceId   string `json:"space_id"              orm:"index;description(空间ID)"`
	SpaceName string `json:"space_name"            orm:"description(空间名称)"`
	Usernames string `json:"usernames"              orm:"null;description(用户名称)"`
}
