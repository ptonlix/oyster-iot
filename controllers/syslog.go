package controllers

import (
	"oyster-iot/services"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type SysLogController struct {
	BaseController
}

type SyslogFileParam struct {
	FileName string `json:"filename" valid:"MaxSize(64)"`
}

// 获取系统日志文件列表
func (s *SysLogController) GetSysLogList() {
	var syslogService services.SysLogService

	filenames, err := syslogService.GetSyslogFileList()
	if err != nil {
		s.Response(400, "查找不到日志")
		return
	}

	s.Response(200, "获取日志列表成功", filenames)
}

// 获取系统日志文件内容
func (s *SysLogController) GetSysLogFileContext() {
	//获取URL参数
	syslogparam := SyslogFileParam{}
	s.Ctx.Input.Bind(&syslogparam.FileName, "filename")
	logs.Debug("filename is %#v", syslogparam.FileName)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&syslogparam)
	if err != nil {
		// handler error
		s.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		s.Response(400, "输入参数错误")
		return
	}

	var syslogService services.SysLogService

	filenames, err := syslogService.GetSyslogFile(syslogparam.FileName)
	if err != nil {
		s.Response(400, "查找不到日志")
		return
	}

	s.Response(200, "获取系统日志成功", string(*filenames))
}
