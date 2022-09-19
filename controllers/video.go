package controllers

import (
	"encoding/json"
	"oyster-iot/init/constants"
	"oyster-iot/models"
	"oyster-iot/services"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

// 视频监控
type VideoController struct {
	BaseController
}

/*
参数	说明
nsId	空间id
offset	在全部Device中的偏移量
line	一次返回多少条
prefix	可以通过gbid前缀进行检索查询
state	按设备状态查询，offline: 离线，online: 在线，notReg: 未注册，locked: 锁定
qtype	按设备类型查询，0:全部，1:摄像头，2:平台
*/
type VideoList struct {
	NsId   string `json:"ns_id"  valid:"Required;MaxSize(100)"`
	Offset int    `json:"offset" valid:"Max(256)"`
	Line   int    `json:"line"   valid:"Required;Max(256)"`
	Prefix string `json:"prefix" valid:"MaxSize(100)"`
	State  string `json:"state"  valid:"MaxSize(100)"`
	Qtype  int    `json:"qtype"  valid:"Range(0,2)"`
}

type VideoDevice struct {
	NsId string `json:"ns_id"  valid:"Required;MaxSize(100)"`
	GbId string `json:"gb_id"  valid:"Required;MaxSize(100)"`
}

type VideoStream struct {
	NsId     string `json:"ns_id"      valid:"Required;MaxSize(100)"`
	StreamId string `json:"stream_id"  valid:"Required;MaxSize(100)"`
}

type VideoRecordList struct {
	NsId     string `json:"ns_id"      valid:"Required;MaxSize(100)"`
	StreamId string `json:"stream_id"  valid:"Required;MaxSize(100)"`
	Start    int    `json:"start"      valid:"Required"`
	End      int    `json:"end"        valid:"Required"`
	Marker   string `json:"marker"     valid:""`
	Line     int    `json:"line"       valid:"Required;Max(256)"`
	Format   string `json:"format"     valid:"MaxSize(100)"`
}

type AddSpaceParam struct {
	Platform string `json:"platform"   valid:"Required;MaxSize(100)"`
	NsId     string `json:"ns_id"      valid:"Required;MaxSize(100)"`
	UserId   int    `json:"userid"     valid:"Required;"`
	Username string `json:"username"   valid:"Required;MaxSize(100);Match(/[a-zA-Z0-9_]{3,100}/)"`
}

type EditSpaceParam struct {
	Id       int    `json:"id"         valid:"Required;"`
	UserId   int    `json:"userid"     valid:""` //允许取消用户关联，所以可以为空
	Username string `json:"username"   valid:"MaxSize(100)"`
}

type DelSpaceParam struct {
	Id int `json:"id"         valid:"Required;"`
}

type VideoDeviceList struct {
	TotalNum   int                   `json:"totalnum"`
	TotalPages int                   `json:"totalpages"`
	DeviceList *[]*models.VideoSpace `json:"list"`
}

type VideoStreamRecordCalendar struct {
	NsId     string `json:"ns_id"      valid:"Required;MaxSize(100)"`
	StreamId string `json:"stream_id"  valid:"Required;MaxSize(100)"`
	Year     string `json:"year"       valid:"Required;"`
	Month    string `json:"month       valid:"Required;"`
}

func (vc *VideoController) AddSpace() {
	addSpaceParam := AddSpaceParam{}
	err := json.Unmarshal(vc.Ctx.Input.RequestBody, &addSpaceParam)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		vc.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&addSpaceParam)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 新增空间
	var videoDb services.VideoDbService
	// 查询当前用户是否已绑定空间,目前一个用户一个空间
	if !videoDb.JudgeUser(addSpaceParam.UserId) {
		//已绑定空间
		vc.Response(400, "该用户已绑定视频监控空间,新增失败")
		return
	}
	if err := videoDb.Add(addSpaceParam.Platform, addSpaceParam.NsId, addSpaceParam.Username, addSpaceParam.UserId); err != nil {
		vc.Response(500, "数据库操作错误")
		return
	}

	vc.Response(200, "添加视频监控空间成功")
}

func (vc *VideoController) EditSpace() {
	editSpaceParam := EditSpaceParam{}
	err := json.Unmarshal(vc.Ctx.Input.RequestBody, &editSpaceParam)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		vc.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&editSpaceParam)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}
	logs.Debug(editSpaceParam)
	// 查询当前用户是否已绑定空间,目前一个用户一个空间
	var videoDb services.VideoDbService
	if editSpaceParam.UserId != 0 {
		if !videoDb.JudgeUser(editSpaceParam.UserId) {
			//已绑定空间
			vc.Response(400, "该用户已绑定视频监控空间,编辑失败")
			return
		}
	}

	videoSpace, err := videoDb.Read(editSpaceParam.Id)
	if err != nil {
		vc.Response(500, "数据库操作错误")
		return
	}

	videoSpace.UserId = editSpaceParam.UserId
	videoSpace.Username = editSpaceParam.Username

	err = videoDb.Update(videoSpace)
	if err != nil {
		vc.Response(500, "数据库操作错误")
		return
	}

	vc.Response(200, "编辑用户空间成功")
	return
}

func (vc *VideoController) DeleteSpace() {
	delSpaceParam := DelSpaceParam{}
	err := json.Unmarshal(vc.Ctx.Input.RequestBody, &delSpaceParam)
	if err != nil {
		logs.Warn("Json Unmarshal Failed!", err.Error())
		vc.Response(500, "系统内部错误")
		return
	}
	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&delSpaceParam)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 删除
	var videoDb services.VideoDbService
	videoSpace := &models.VideoSpace{Id: delSpaceParam.Id}
	if err := videoDb.Delete(videoSpace); err != nil {
		vc.Response(500, "数据库操作错误")
		return
	}

	vc.Response(200, "删除视频监控空间成功")
	return
}

func (vc *VideoController) ListSpace() {
	//获取URL参数
	pageparam := PageParam{}
	vc.Ctx.Input.Bind(&pageparam.Pagesize, "pagesize")
	vc.Ctx.Input.Bind(&pageparam.Pagenum, "pagenum")
	vc.Ctx.Input.Bind(&pageparam.Keyword, "keyword")
	logs.Debug("pagesize is %#v, pagenum is %#v keyword is %#v", pageparam.Pagesize, pageparam.Pagenum, pageparam.Keyword)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 获取设备数据
	var videoDb services.VideoDbService

	totalNum, totalPages, videoSpace, err := videoDb.GetSpacesByPageAndKey(pageparam.Pagesize, pageparam.Pagenum, pageparam.Keyword)
	if err != nil {
		vc.Response(400, "查找不到监控视频空间")
		return
	}

	retList := VideoDeviceList{
		TotalNum:   totalNum,
		TotalPages: totalPages,
		DeviceList: &videoSpace,
	}
	vc.Response(200, "获取监控视频空间列表成功", retList)

}

// 查询空间和用户名
func (vc *VideoController) ListSpaceAndUser() {
	// 获取设备数据
	var videoDb services.VideoDbService

	videoSpaceUsers, err := videoDb.GetSpacesAndUser()
	if err != nil {
		vc.Response(400, "查找不到监控设备")
		return
	}

	vc.Response(200, "获取监控设备列表成功", videoSpaceUsers)
}

/*下面为小程序提供的API接口*/

func (vc *VideoController) GetSpaceByUser() {
	userId, _, err := vc.GetUserInfo()
	if err != nil {
		vc.Response(500, "查询当前用户失败")
		return
	}
	//查询当前用户下的监控空间信息
	var videoDb services.VideoDbService
	var SpaceInfo *models.VideoSpace
	if SpaceInfo, err = videoDb.ReadByUserId(userId); err != nil {
		vc.Response(200, "当前用户未开通视频监控业务,请联系管理员开通", nil)
		return
	}

	vc.Response(200, "获取当前用户视频监控业务成功", SpaceInfo)
	return
}

func (vc *VideoController) List() {
	//获取URL参数
	pageparam := VideoList{}
	vc.Ctx.Input.Bind(&pageparam.Offset, "offset")
	vc.Ctx.Input.Bind(&pageparam.Line, "line")
	vc.Ctx.Input.Bind(&pageparam.NsId, "ns_id")
	vc.Ctx.Input.Bind(&pageparam.Prefix, "prefix")
	vc.Ctx.Input.Bind(&pageparam.State, "state")
	vc.Ctx.Input.Bind(&pageparam.Qtype, "qtype")
	logs.Debug("pageparam is %#v", pageparam)

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&pageparam)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, pageparam.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	videoS.SDK.Create()
	videoDevs, sum, err := videoS.SDK.GetVideoDevicesList(pageparam.Offset, pageparam.Line, pageparam.Prefix, pageparam.State, pageparam.Qtype)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}

	type VideoDeviceList struct {
		TotalNum int64       `json:"totalnum"`
		List     interface{} `json:"list"`
	}
	videoDeviceList := VideoDeviceList{
		TotalNum: sum,
		List:     videoDevs,
	}

	vc.Response(200, "获取视频监控设备列表成功", &videoDeviceList)
}

func (vc *VideoController) GetVideoDevice() {
	videoInfo := VideoDevice{}
	vc.Ctx.Input.Bind(&videoInfo.NsId, "ns_id")
	vc.Ctx.Input.Bind(&videoInfo.GbId, "gb_id")

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&videoInfo)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}
	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, videoInfo.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误")
		return
	}
	videoS.SDK.Create()
	videoDev, err := videoS.SDK.GetVideoDeviceInfo(videoInfo.GbId)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误")
		return
	}
	vc.Response(200, "获取视频监控设备信息成功", &videoDev)
}

func (vc *VideoController) GetVideoPlayUrlList() {
	videoStream := VideoStream{}
	vc.Ctx.Input.Bind(&videoStream.NsId, "ns_id")
	vc.Ctx.Input.Bind(&videoStream.StreamId, "stream_id")

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&videoStream)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}
	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, videoStream.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误")
		return
	}
	videoS.SDK.Create()
	videoUrls, err := videoS.SDK.GetVideoDevicePlayUrl(videoStream.StreamId)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误")
		return
	}
	vc.Response(200, "获取视频监控设备信息成功", &videoUrls)
}

func (vc *VideoController) GetVideoRecordList() {
	videoRecords := VideoRecordList{}
	vc.Ctx.Input.Bind(&videoRecords.NsId, "ns_id")
	vc.Ctx.Input.Bind(&videoRecords.StreamId, "stream_id")
	vc.Ctx.Input.Bind(&videoRecords.Start, "start")
	vc.Ctx.Input.Bind(&videoRecords.End, "end")
	vc.Ctx.Input.Bind(&videoRecords.Marker, "marker")
	vc.Ctx.Input.Bind(&videoRecords.Line, "line")
	vc.Ctx.Input.Bind(&videoRecords.Format, "format")

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&videoRecords)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, videoRecords.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误")
		return
	}
	videoS.SDK.Create()
	videoRecordsList, err := videoS.SDK.GetVideoDeviceRecordList(videoRecords.StreamId, videoRecords.Start, videoRecords.End, videoRecords.Marker, videoRecords.Line, videoRecords.Format)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,获取视频回放记录失败")
		return
	}

	vc.Response(200, "获取视频回放列表成功", &videoRecordsList)
}

func (vc *VideoController) GetVideoStream() {
	videoStream := VideoStream{}
	vc.Ctx.Input.Bind(&videoStream.NsId, "ns_id")
	vc.Ctx.Input.Bind(&videoStream.StreamId, "stream_id")

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&videoStream)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, videoStream.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	videoS.SDK.Create()
	videoStreamInfo, err := videoS.SDK.GetVideoStreamInfo(videoStream.StreamId)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	vc.Response(200, "获取视控设备视频流信息成功", &videoStreamInfo)
}

func (vc *VideoController) StartVideoRecord() {
	videoStream := VideoStream{}
	vc.Ctx.Input.Bind(&videoStream.NsId, "ns_id")
	vc.Ctx.Input.Bind(&videoStream.StreamId, "stream_id")

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&videoStream)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, videoStream.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	videoS.SDK.Create()
	err = videoS.SDK.StartVideoRecord(videoStream.StreamId)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	vc.Response(200, "启动视频录制成功")
}

func (vc *VideoController) StopVideoRecord() {
	videoStream := VideoStream{}
	vc.Ctx.Input.Bind(&videoStream.NsId, "ns_id")
	vc.Ctx.Input.Bind(&videoStream.StreamId, "stream_id")

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&videoStream)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, videoStream.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	videoS.SDK.Create()
	err = videoS.SDK.StopVideoRecord(videoStream.StreamId)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	vc.Response(200, "停止视频录制成功")
}

func (vc *VideoController) GetRecordCalendar() {
	vSRecordCalendar := VideoStreamRecordCalendar{}
	vc.Ctx.Input.Bind(&vSRecordCalendar.NsId, "ns_id")
	vc.Ctx.Input.Bind(&vSRecordCalendar.StreamId, "stream_id")
	vc.Ctx.Input.Bind(&vSRecordCalendar.Year, "year")
	vc.Ctx.Input.Bind(&vSRecordCalendar.Month, "month")

	// 校验输入参数是否合法
	v := validation.Validation{}
	b, err := v.Valid(&vSRecordCalendar)
	if err != nil {
		// handler error
		vc.Response(500, "系统内部错误")
		return
	}
	if !b {
		// validation does not pass
		for _, err := range v.Errors {
			logs.Warn(err.Key, err.Message)
		}
		vc.Response(400, "输入参数错误")
		return
	}

	// 创建视频监控服务
	videoS, err := services.NewVideoService("QINIU", constants.QiniuAK, constants.QiniuSK, vSRecordCalendar.NsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	videoS.SDK.Create()
	daysMap, err := videoS.SDK.GetVideoDeviceRecordCalendar(vSRecordCalendar.StreamId, vSRecordCalendar.Year, vSRecordCalendar.Month)
	if err != nil {
		logs.Error(err)
		vc.Response(500, "视频监控服务错误,请联系管理员")
		return
	}
	vc.Response(200, "获取视频回放日历成功", daysMap)
}
