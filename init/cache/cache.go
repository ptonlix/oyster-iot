package cache

import (
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/logs"
)

var Bm cache.Cache
var Err error

func init() {
	Bm, Err = cache.NewCache("memory", `{"interval":60}`) //使用内存做缓存
	if Err != nil {
		logs.Error("初始化cache失败", Err)
	}
}
