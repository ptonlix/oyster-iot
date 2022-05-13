package routers

import (
	"fmt"
	"oyster-iot/controllers"
	"oyster-iot/init/cache"
	"strings"

	c "context"

	jwt "oyster-iot/utils"

	"github.com/beego/beego/v2/adapter"
	"github.com/beego/beego/v2/adapter/context"
	"github.com/beego/beego/v2/server/web"
)

func init() {

	//AuthMiddle()

	api := web.NewNamespace("/api",
		// 登录
		web.NSRouter("/auth/login", &controllers.AuthController{}, "*:Login"),
		web.NSRouter("/auth/logout", &controllers.AuthController{}, "*:Logout"),
		web.NSRouter("/auth/refresh", &controllers.AuthController{}, "*:Refresh"),
		// web.NSRouter("/auth/me", &controllers.AuthController{}, "*:Me"),
		// web.NSRouter("/auth/register", &controllers.AuthController{}, "*:Register"),

		// 设备管理
		web.NSRouter("/device/edit", &controllers.DeviceController{}, "*:Edit"),
		web.NSRouter("/device/add", &controllers.DeviceController{}, "*:Add"),
		web.NSRouter("/device/delete", &controllers.DeviceController{}, "*:Delete"),
		web.NSRouter("/device/list", &controllers.DeviceController{}, "*:List"),
		web.NSRouter("/device/listforbusiness", &controllers.DeviceController{}, "*:ListForBusiness"),
		web.NSRouter("/device/listfornilbusiness", &controllers.DeviceController{}, "*:ListForNilBusiness"),
		web.NSRouter("/device/updateforbusiness", &controllers.DeviceController{}, "*:UpdateForBusiness"),
		web.NSRouter("/device/listforindex", &controllers.DeviceController{}, "*:ListForIndex"),

		// 业务管理
		web.NSNamespace("/business",
			web.NSRouter("/add", &controllers.BusinessController{}, "*:Add"),
			web.NSRouter("/edit", &controllers.BusinessController{}, "*:Edit"),
			web.NSRouter("/delete", &controllers.BusinessController{}, "*:Delete"),
			web.NSRouter("/list", &controllers.BusinessController{}, "*:List"),
			// 获取设备最近一天的温度情况
			web.NSRouter("/temperature/devinday", &controllers.TempController{}, "*:GetTempInDay"),
			// 获取最新的温度信息
			web.NSRouter("/temperature/dev", &controllers.TempController{}, "*:GetTemp"),
			web.NSRouter("/temperature/sendtempcmd", &controllers.TempController{}, "*:SendTempCmd"),
			web.NSRouter("/salinity/devinday", &controllers.SaltController{}, "*:GetSaltInDay"),
			web.NSRouter("/salinity/dev", &controllers.SaltController{}, "*:GetSalt"),
			web.NSRouter("/salinity/sendsaltcmd", &controllers.SaltController{}, "*:SendSaltCmd"),
		),

		// 系统信息
		web.NSNamespace("sys/",
			web.NSRouter("/emqmetrisc", &controllers.EmqExportController{}, "*:GetMetrics"),
			web.NSRouter("/sysinfo", &controllers.SysController{}, "*:GetSysinfo"),
		),
	)
	web.AddNamespace(api)
}

// AuthMiddle 中间件
func AuthMiddle() {
	//不需要验证的url
	noLogin := map[string]interface{}{
		"api/auth/login":    0,
		"api/auth/refresh":  0,
		"api/auth/register": 1,
	}
	var filterLogin = func(ctx *context.Context) {
		url := strings.TrimLeft(ctx.Input.URL(), "/")
		if !isAuthExceptUrl(strings.ToLower(url), noLogin) {
			//获取TOKEN
			if len(ctx.Request.Header["Authorization"]) == 0 {
				ctx.Redirect(302, "/login.html")
				return
			}
			authorization := ctx.Request.Header["Authorization"][0]
			userToken := authorization[len(jwt.JWTType)+1:]
			_, err := jwt.ParseCliamsToken(userToken)
			if err != nil {
				// 异常
				ctx.Redirect(302, "/login.html")
				return
			}
			s, _ := cache.Bm.IsExist(c.TODO(), userToken)
			if !s {
				ctx.Redirect(302, "/login.html")
				return
			}
		}
	}
	adapter.InsertFilter("/api/*", adapter.BeforeRouter, filterLogin)
}

func isAuthExceptUrl(url string, m map[string]interface{}) bool {
	urlArr := strings.Split(url, "/")
	if len(urlArr) > 3 {
		url = fmt.Sprintf("%s/%s/%s", urlArr[0], urlArr[1], urlArr[2])
	}
	_, ok := m[url]
	if ok {
		return true
	}
	return false
}
