package controllers

import "oyster-iot/services"

type SysController struct {
	BaseController
}

// 获取服务器运行的指标数据
func (s *SysController) GetSysinfo() {

	var sys services.SysinfoService

	data, err := sys.GetInfoByPsutil()
	if err != nil {
		s.Response(500, "获取系统服务指标数据失败")
		return
	}

	s.Response(200, "获取系统服务指标数据成功", &data)
}
