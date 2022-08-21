package log

import (
	"fmt"
	"oyster-iot/init/constants"
	"oyster-iot/models"
	"oyster-iot/services"

	"github.com/beego/beego/v2/core/logs"
	agent "github.com/wenlng/go-user-agent"
)

const (
	LogHeaderFlag = "operresult"
	Success       = "success"
	Failed        = "failed"
)

var URLToOPERTYPE map[string]string = make(map[string]string)

type LogFlag struct {
}
type OperLogContext struct {
	Ip          string
	Address     string
	OperType    string
	URL         string
	Username    string
	UserType    string
	OsInfo      string
	BrowserInfo string
	RequestBody string
	Flag        string
}

func init() {
	//logs.SetLogger(logs.AdapterConsole, `{"level":7,"color":true}`)
	config := fmt.Sprintf(`{"filename":"%s","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`, constants.LogFileDir+constants.LogFileName)
	logs.SetLogger(logs.AdapterFile, config)

	URLToOPERTYPE["/api/auth/login"] = "登录"
	URLToOPERTYPE["/api/auth/logout"] = "注销"
	URLToOPERTYPE["/api/auth/getuserinfo"] = "获取当前用户信息"
	URLToOPERTYPE["/api/auth/loginmanage"] = "登录管理平台"
	URLToOPERTYPE["/api/auth/logoutmanage"] = "注销管理平台"
	URLToOPERTYPE["/api/auth/getmanageuserinfo"] = "获取当前管理平台用户信息"

	URLToOPERTYPE["/api/user/add"] = "添加用户"
	URLToOPERTYPE["/api/user/edit"] = "编辑用户"
	URLToOPERTYPE["/api/user/delete"] = "删除用户"
	URLToOPERTYPE["/api/user/list"] = "获取用户列表"
	URLToOPERTYPE["/api/user/resetpassword"] = "重置用户密码"

	URLToOPERTYPE["/api/device/add"] = "添加设备"
	URLToOPERTYPE["/api/device/edit"] = "编辑设备"
	URLToOPERTYPE["/api/device/delete"] = "删除设备"
	URLToOPERTYPE["/api/device/list"] = "获取设备列表"
	URLToOPERTYPE["/api/device/listforbusiness"] = "获取当前业务下的设备列表"
	URLToOPERTYPE["/api/device/listfornilbusiness"] = "获取未绑定业务的设备列表"
	URLToOPERTYPE["/api/device/updateforbusiness"] = "批量更新设备业务ID"
	URLToOPERTYPE["/api/device/listforindex"] = "获取当前业务下的设备列表(已分类)"
	URLToOPERTYPE["/api/device/listforallnum"] = "获取设备总数信息"

	URLToOPERTYPE["/api/business/add"] = "添加业务"
	URLToOPERTYPE["/api/business/edit"] = "编辑业务"
	URLToOPERTYPE["/api/business/delete"] = "删除业务"
	URLToOPERTYPE["/api/business/list"] = "获取业务列表"
	URLToOPERTYPE["/api/business/listforallnum"] = "获取业务总数信息"

	URLToOPERTYPE["/api/business/temperature/devinday"] = "获取设备一天的温度数据"
	URLToOPERTYPE["/api/business/temperature/dev"] = "获取设备最新的温度数据"
	URLToOPERTYPE["/api/business/temperature/sendtempcmd"] = "发送获取设备温度数据命令"

	URLToOPERTYPE["/api/business/salinity/devinday"] = "获取设备一天的盐度数据"
	URLToOPERTYPE["/api/business/salinity/dev"] = "获取设备最新的盐度数据"
	URLToOPERTYPE["/api/business/salinity/sendsaltcmd"] = "发送获取设备温度数据命令"

	URLToOPERTYPE["/api/sys/emqmetrisc"] = "获取EMQX信息"
	URLToOPERTYPE["/api/sys/sysinfo"] = "获取系统平台信息"
	URLToOPERTYPE["/api/sys/operlogdelete"] = "删除操作日志记录"
}

func (o *OperLogContext) FilterOperLog() {
	// logs.Info(o)
	if v, ok := URLToOPERTYPE[o.URL]; ok {
		o.OperType = v
	} else {
		return
	}

	operLogData := &models.Operlog{
		Ip:           o.Ip,
		Address:      o.Address,
		Opertype:     o.OperType,
		Operuser:     o.Username,
		OperuserType: o.UserType,
		OsInfo:       o.OsInfo,
		BrowserInfo:  o.BrowserInfo,
		RequestBody:  o.RequestBody,
		Result:       o.Flag,
		Url:          o.URL,
	}

	var operLogService services.OperLogService
	err := operLogService.Add(operLogData)
	if err != nil {
		logs.Error("数据库操作错误")
	}
}

func (o *OperLogContext) ParseUserAgent(userAgent string) {
	o.OsInfo = agent.GetOsName(userAgent)
	o.BrowserInfo = agent.GetBrowserName(userAgent)
}
