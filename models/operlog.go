package models

import "time"

type Operlog struct {
	Id           int       `json:"id"            orm:"auto;description(操作日志ID)"`
	Ip           string    `json:"ip"            orm:"description(操作的IP地址)"`
	Address      string    `json:"address"       orm:"description(操作的IP地址物理地址)"`
	Operuser     string    `json:"operuser"      orm:"description(操作人)"`
	OperuserType string    `json:"operuser_type" orm:"description(操作人类型:管理平台、小程序)"`
	OsInfo       string    `json:"os_info"       orm:"description(操作系统信息)`
	BrowserInfo  string    `json:"browser_info"  orm:"description(浏览器信息)`
	Opertype     string    `json:"opertype"      orm:"description(操作类型)"`
	Url          string    `json:"url"           orm:"description(操作URL)"`
	Result       string    `json:"result"        orm:"description(操作结果)"`
	RequestBody  string    `json:"request_body"  orm:"type(text);description(请求的参数)"`
	Createdat    time.Time `json:"created_at"    orm:"auto_now_add;type(datetime);description(创建时间)"`
	Updatedat    time.Time `json:"updated_at"    orm:"auto_now;type(datetime);description(更新时间)"`
}
