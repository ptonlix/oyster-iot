package models

import "time"

type Business struct {
	Id        int       `json:"id"              orm:"auto;description(业务ID)"`
	Name      string    `json:"business_name"   orm:"description(业务名称)"`
	Remark    string    `json:"remark"          orm:"size(255);description(业务信息说明)"`
	Createdat time.Time `json:"created_at"      orm:"auto_now_add;type(datetime);description(创建时间)"`
	Updatedat time.Time `json:"updated_at"      orm:"auto_now;type(datetime);description(更新时间)"`
}
