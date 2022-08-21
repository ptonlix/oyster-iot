package routers

import (
	"encoding/json"
	"fmt"
	"oyster-iot/controllers"
	"oyster-iot/init/cache"
	"strings"

	jwt "oyster-iot/utils"

	oysterlog "oyster-iot/init/log"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {

	AuthMiddle()
	OperLogMiddle()

	api := web.NewNamespace("/api",
		// 登录
		web.NSRouter("/auth/login", &controllers.AuthController{}, "*:Login"),
		web.NSRouter("/auth/logout", &controllers.AuthController{}, "*:Logout"),
		web.NSRouter("/auth/refresh", &controllers.AuthController{}, "*:Refresh"),
		web.NSRouter("/auth/getuserinfo", &controllers.AuthController{}, "*:Getuserinfo"),
		// web.NSRouter("/auth/me", &controllers.AuthController{}, "*:Me"),TenantGetById
		// web.NSRouter("/auth/register", &controllers.AuthController{}, "*:Register"), GetManageUserinfo
		web.NSRouter("/auth/loginmanage", &controllers.AuthController{}, "*:LoginManage"),
		web.NSRouter("/auth/logoutmanage", &controllers.AuthController{}, "*:LogoutManage"),
		web.NSRouter("/auth/getmanageuserinfo", &controllers.AuthController{}, "*:GetManageUserinfo"),

		web.NSRouter("/user/add", &controllers.UserController{}, "*:Add"),
		web.NSRouter("/user/edit", &controllers.UserController{}, "*:Edit"),
		web.NSRouter("/user/delete", &controllers.UserController{}, "*:Delete"),
		web.NSRouter("/user/list", &controllers.UserController{}, "*:List"),
		web.NSRouter("/user/resetpassword", &controllers.UserController{}, "*:ResetPassword"),

		// 设备管理
		web.NSRouter("/device/edit", &controllers.DeviceController{}, "*:Edit"),
		web.NSRouter("/device/add", &controllers.DeviceController{}, "*:Add"),
		web.NSRouter("/device/delete", &controllers.DeviceController{}, "*:Delete"),
		web.NSRouter("/device/list", &controllers.DeviceController{}, "*:List"),
		web.NSRouter("/device/listforbusiness", &controllers.DeviceController{}, "*:ListForBusiness"),
		web.NSRouter("/device/listfornilbusiness", &controllers.DeviceController{}, "*:ListForNilBusiness"),
		web.NSRouter("/device/updateforbusiness", &controllers.DeviceController{}, "*:UpdateForBusiness"),
		web.NSRouter("/device/listforindex", &controllers.DeviceController{}, "*:ListForIndex"),
		web.NSRouter("/device/listforallnum", &controllers.DeviceController{}, "*:ListForAllNum"),

		// 业务管理
		web.NSNamespace("/business",
			web.NSRouter("/add", &controllers.BusinessController{}, "*:Add"),
			web.NSRouter("/edit", &controllers.BusinessController{}, "*:Edit"),
			web.NSRouter("/delete", &controllers.BusinessController{}, "*:Delete"),
			web.NSRouter("/list", &controllers.BusinessController{}, "*:List"),
			web.NSRouter("/listforallnum", &controllers.BusinessController{}, "*:ListForAllNum"),
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
			web.NSRouter("/operloglist", &controllers.OperLogController{}, "*:List"),
			web.NSRouter("/operlogsearch", &controllers.OperLogController{}, "*:Search"),
			web.NSRouter("/operlogdelete", &controllers.OperLogController{}, "*:Delete"),
			web.NSRouter("/sysloglist", &controllers.SysLogController{}, "*:GetSysLogList"),
			web.NSRouter("/syslogcontext", &controllers.SysLogController{}, "*:GetSysLogFileContext"),
		),
		// 视频监控系统
		web.NSNamespace("video/",
			web.NSRouter("/list", &controllers.VideoController{}, "*:List"),
			web.NSRouter("/getdevice", &controllers.VideoController{}, "*:GetVideoDevice"),
			web.NSRouter("/geturls", &controllers.VideoController{}, "*:GetVideoPlayUrlList"),
			web.NSRouter("/getrecordlist", &controllers.VideoController{}, "*:GetVideoRecordList"),
			web.NSRouter("/getstream", &controllers.VideoController{}, "*:GetVideoStream"),
			web.NSRouter("/startrecord", &controllers.VideoController{}, "*:StartVideoRecord"),
			web.NSRouter("/stoprecord", &controllers.VideoController{}, "*:StopVideoRecord"),
			web.NSNamespace("space/",
				web.NSRouter("/add", &controllers.VideoController{}, "*:AddSpace"),
				web.NSRouter("/edit", &controllers.VideoController{}, "*:EditSpace"),
				web.NSRouter("/delete", &controllers.VideoController{}, "*:DeleteSpace"),
				web.NSRouter("/list", &controllers.VideoController{}, "*:ListSpace"),
				web.NSRouter("/getbyuser", &controllers.VideoController{}, "*:GetSpaceByUser"),
			),
		),
	)
	web.AddNamespace(api)
}

// AuthMiddle 中间件
func AuthMiddle() {
	//不需要验证的url
	noLogin := map[string]interface{}{
		"api/auth/login":       0,
		"api/auth/refresh":     0,
		"api/auth/register":    1,
		"api/auth/loginmanage": 0,
	}
	var filterLogin = func(ctx *context.Context) {
		url := strings.TrimLeft(ctx.Input.URL(), "/")
		if !isAuthExceptUrl(strings.ToLower(url), noLogin) {
			//获取TOKEN
			if len(ctx.Request.Header["Authorization"]) == 0 {
				// ctx.Redirect(302, "/login.html")
				ctx.ResponseWriter.WriteHeader(401)
				ctx.WriteString("login out")
				return
			}
			authorization := ctx.Request.Header["Authorization"][0]
			if len(authorization) <= len(jwt.JWTType)+1 {
				// 异常
				ctx.ResponseWriter.WriteHeader(401)
				ctx.WriteString("login out")
				return
			}
			userToken := authorization[len(jwt.JWTType)+1:]
			_, err := jwt.ParseCliamsToken(userToken)
			if err != nil {
				// 异常
				ctx.ResponseWriter.WriteHeader(401)
				ctx.WriteString("login out")
				return
			}
			s := cache.Bm.IsExist(userToken)
			if !s {
				ctx.ResponseWriter.WriteHeader(401)
				ctx.WriteString("login out")
				return
			}
		}
	}
	web.InsertFilter("/api/*", web.BeforeRouter, filterLogin)
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

// 记录操作日志 中间件
func OperLogMiddle() {
	var FilterUser = func(ctx *context.Context) {
		URL := ctx.Input.URL() // 获取请求的URL
		username := "system"
		usertype := "system"
		//ctx.Input.Context.Request.ParseForm()
		requestBody := ""
		if ctx.Input.Context.Request.Method == "GET" {
			urlParam := ctx.Input.Context.Request.URL.Query()
			b, err := json.Marshal(urlParam)
			if err != nil {
				logs.Info("json.Marshal failed:", err)
				return
			}
			requestBody = string(b)
		} else {
			requestBody = string(ctx.Input.RequestBody)
		}
		if len(ctx.Request.Header["Authorization"]) != 0 {
			authorization := ctx.Request.Header["Authorization"][0]
			if len(authorization) > len(jwt.JWTType)+1 {
				userToken := authorization[len(jwt.JWTType)+1:]
				userInfo, err := jwt.ParseCliamsToken(userToken)
				if err == nil {
					username = userInfo.Username
					usertype = userInfo.Usertype
				}
			}
		}
		ip := ""
		if len(ctx.Request.Header["X-Forwarded-For"]) != 0 {
			ip = ctx.Request.Header["X-Forwarded-For"][0]
		}
		op := oysterlog.OperLogContext{
			Ip:          ip,
			Address:     "-", //TODO 后续增加redis缓存再做
			URL:         URL,
			Username:    username,
			UserType:    usertype,
			RequestBody: requestBody,
			Flag:        ctx.Output.Context.ResponseWriter.ResponseWriter.Header().Get(oysterlog.LogHeaderFlag),
		}
		op.ParseUserAgent(ctx.Request.Header["User-Agent"][0]) //解析浏览器和操作系统
		op.FilterOperLog()
	}
	//设置 returnOnOutput 的值(默认 true), 即如果有输出是否跳过其他过滤器，默认只要有输出就不再执行其他过滤器，即执行完controller之后不会执行后面的过滤器
	web.InsertFilter("/*", web.AfterExec, FilterUser, web.WithReturnOnOutput(false))

}
