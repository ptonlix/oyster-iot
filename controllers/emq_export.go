package controllers

import "oyster-iot/services"

type EmqExportController struct {
	BaseController
}

// 获取 EMQX 服务的指标数据
func (e *EmqExportController) GetMetrics() {
	var emqexp services.EmqExpService

	data, err := emqexp.GetEmqMetrics()
	if err != nil {
		e.Response(500, "获取EMQX服务指标数据失败")
		return
	}

	e.Response(200, "获取EMQX服务指标数据成功", data)
}
