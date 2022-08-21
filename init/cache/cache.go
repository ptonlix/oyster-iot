package cache

import (
	"encoding/json"

	"github.com/beego/beego/cache"
	_ "github.com/beego/beego/cache/redis"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

var Bm cache.Cache
var Err error

func init() {
	// Bm, Err = cache.NewCache("memory", `{"interval":60}`) //使用内存做缓存
	// if Err != nil {
	// 	logs.Error("初始化cache失败", Err)
	// }

	//  使用Redis做缓存
	redisHost, _ := beego.AppConfig.String("redishost")
	redisPort, _ := beego.AppConfig.String("redisport")
	redisDbnum, _ := beego.AppConfig.String("redisdb")
	redisPasswd, _ := beego.AppConfig.String("redispasswd")

	cacheRedisConn, _ := json.Marshal(map[string]string{
		"key":      "beecacheRedis",
		"conn":     redisHost + ":" + redisPort,
		"dbNum":    redisDbnum,
		"password": redisPasswd,
	})
	logs.Info("Redis  Config :", string(cacheRedisConn))
	Bm, Err = cache.NewCache("redis", string(cacheRedisConn))
	if Err != nil {
		logs.Error("初始化cache失败 ", Err)
	}
}
