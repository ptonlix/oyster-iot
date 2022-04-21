package cache

import (
	"log"

	"github.com/beego/beego/v2/client/cache"
)

var Bm cache.Cache
var Err error

func init() {
	Bm, Err = cache.NewCache("memory", `{"interval":60}`) //使用内存做缓存
	if Err != nil {
		log.Println("初始化cache失败", Err)
	}
}
