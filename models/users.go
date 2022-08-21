package models

import (
	"time"
)

type Users struct {
	Id        int       `json:"id"          orm:"pk;auto"`
	Enabled   bool      `json:"enabled"`
	Email     string    `json:"email"       orm:"size(255)"`
	Username  string    `json:"username"    orm:"size(255);unique"`
	Password  string    `json:"-"           orm:"size(255)"`
	Firstname string    `json:"first_name"  orm:"size(255)"`
	Lastname  string    `json:"last_name"   orm:"size(255)"`
	Mobile    string    `json:"mobile"      orm:"size(255)"`
	Remark    string    `json:"remark"      orm:"size(255)"`
	IsAdmin   bool      `json:"is_admin"`
	Wxopenid  string    `json:"wx_openid"   orm:"size(255)"` // 微信openid
	Wxunionid string    `json:"wx_unionid"  orm:"size(255)"` // 微信unionid
	Createdat time.Time `json:"created_at"  orm:"auto_now_add;type(datetime);description(创建时间)"`
	Updatedat time.Time `json:"updated_at"  orm:"auto_now;type(datetime);description(更新时间)"`
}
