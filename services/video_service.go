package services

import (
	"errors"

	"github.com/beego/beego/v2/core/logs"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/qvs"
)

// 视频监控SDK接口
type VideoSDK interface {
	//创建单例
	Create()
	/* 	获取监控设备列表
	参数	说明
	offset	在全部Device中的偏移量
	line	一次返回多少条
	prefix	可以通过gbid前缀进行检索查询
	state	按设备状态查询，offline: 离线，online: 在线，notReg: 未注册，locked: 锁定
	qtype	按设备类型查询，0:全部，1:摄像头，2:平台
	*/
	GetVideoDevicesList(offset, line int, prefix, state string, qType int) (*[]VideoDevice, int64, error)

	GetVideoDeviceInfo(gbId string) (*VideoDevice, error)         //获取单个设备信息
	GetVideoDevicePlayUrl(streamId string) (*VideoPlayUrl, error) //获取设备的流ID
	/* 获取回放录制列表
	参数	说明
	streamId	流id
	start	查询开始时间(unix时间戳，单位为秒)
	end	查询结束时间(unix时间戳，单位为秒)
	marker	上次返回的marker值
	line	一次返回多少条
	format	查询录制格式，""为查询所有格式视频
	*/
	GetVideoDeviceRecordList(streamId string, start, end int, marker string, line int, format string) (*VideoRecordList, error)
	// 获取空间信息
	GetVideoSpaceInfo(nsId string) (interface{}, error)
	// 查询视频流信息
	GetVideoStreamInfo(streamId string) (interface{}, error)
	// 开始录制
	StartVideoRecord(streamId string) error
	// 停止录制
	StopVideoRecord(streamId string) error
}

type ChannelInfo struct {
}
type VideoDevice struct {
	NamespaceId     string      `json:"namespace_id"`      //所属的空间ID
	Name            string      `json:"name"`              //设备名称
	GBId            string      `json:"gb_id"`             //设备国标ID
	Type            int         `json:"type"`              //设备类型，取值1（摄像头）, 2（平台）
	Username        string      `json:"username"`          //用户名
	Password        string      `json:"password"`          //密码
	PullIfRegister  bool        `json:"pull_if_register"`  //注册成功后启动拉流, 默认关闭
	Desc            string      `json:"desc"`              //设备描述
	NamespaceName   string      `json:"namespace_name"`    //所属的空间名称
	State           string      `json:"state"`             //状态（offline: 离线, online: 在线, notReg: 未注册, locked: 锁定）
	Channels        int         `json:"channels"`          //设备通道数
	Vendor          string      `json:"vendor"`            //厂商
	CreatedAt       int64       `json:"created_at"`        //创建时间，Unix时间戳，秒
	UpdatedAt       int64       `json:"updated_at"`        //更新时间，Unix时间戳，秒
	LastRegisterAt  int64       `json:"last_register_at"`  //上一次注册时间，Unix时间戳，秒
	LastKeepaliveAt int64       `json:"last_keepalive_at"` //上一次心跳时间，Unix时间戳，秒
	ChannelInfo     interface{} `json:"channel_info"`
}

type VideoPlayUrl struct {
	Rtmp string `json:"rtmp"` // rtmp播放地址
	Flv  string `json:"flv"`  // flv播放地址
	Hls  string `json:"hls"`  // hls播放地址
}

type VideoRecordList struct {
	Marker string      `json:"marker"`
	List   interface{} `json:"list"`
}

// 封装SDK(SDK)
type VideoService struct {
	Platform string   //TENCENT(腾讯),QINIU(七牛)
	SDK      VideoSDK //视频监控SDK
}

type QiniuSDK struct {
	Ak      string
	Sk      string
	NsId    string //所属的空间ID
	manager *qvs.Manager
}

func NewVideoService(platform string, Option ...string) (*VideoService, error) {
	vs := &VideoService{}
	vs.Platform = platform
	if vs.Platform == "QINIU" {
		qiniuSDK := &QiniuSDK{
			Ak:   Option[0],
			Sk:   Option[1],
			NsId: Option[2], // "3nm4x17o751vq"
		}
		vs.SDK = qiniuSDK
	} else {
		return nil, errors.New("Not Support Platform, New VideoService Failed!")
	}
	return vs, nil
}

func (q *QiniuSDK) Create() {
	mac := auth.New(q.Ak, q.Sk)
	q.manager = qvs.NewManager(mac, nil)
}

func (q *QiniuSDK) GetVideoDevicesList(offset, line int, prefix, state string, qType int) (*[]VideoDevice, int64, error) {
	logs.Info(offset, line)
	devs, sum, err := q.manager.ListDevice(q.NsId, offset, line, prefix, state, qType)
	if err != nil {
		logs.Error(err)
		return nil, 0, err
	}
	logs.Info(devs)
	videoDevs := []VideoDevice{}
	for _, device := range devs {
		videoDev := VideoDevice{
			NamespaceId:     device.NamespaceId,
			Name:            device.Name,
			GBId:            device.GBId,
			Type:            device.Type,
			Username:        device.Username,
			Password:        device.Password,
			PullIfRegister:  device.PullIfRegister,
			Desc:            device.Desc,
			NamespaceName:   device.NamespaceName,
			State:           device.State,
			Channels:        device.Channels,
			Vendor:          device.Vendor,
			CreatedAt:       device.CreatedAt,
			UpdatedAt:       device.UpdatedAt,
			LastRegisterAt:  device.LastRegisterAt,
			LastKeepaliveAt: device.LastKeepaliveAt,
		}
		videoDevs = append(videoDevs, videoDev)
	}
	return &videoDevs, sum, err
}

func (q *QiniuSDK) GetVideoDeviceInfo(gbId string) (*VideoDevice, error) {
	device, err := q.manager.QueryDevice(q.NsId, gbId)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	videoDev := VideoDevice{
		NamespaceId:     device.NamespaceId,
		Name:            device.Name,
		GBId:            device.GBId,
		Type:            device.Type,
		Username:        device.Username,
		Password:        device.Password,
		PullIfRegister:  device.PullIfRegister,
		Desc:            device.Desc,
		NamespaceName:   device.NamespaceName,
		State:           device.State,
		Channels:        device.Channels,
		Vendor:          device.Vendor,
		CreatedAt:       device.CreatedAt,
		UpdatedAt:       device.UpdatedAt,
		LastRegisterAt:  device.LastRegisterAt,
		LastKeepaliveAt: device.LastKeepaliveAt,
	}
	if device.Type == 2 {
		chanelInfo, err := q.manager.ListChannels(q.NsId, gbId, "")
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		videoDev.ChannelInfo = chanelInfo
	}
	return &videoDev, nil
}

func (q *QiniuSDK) GetVideoDevicePlayUrl(streamId string) (*VideoPlayUrl, error) {
	ret, err := q.manager.DynamicPublishPlayURL(q.NsId, streamId, &qvs.DynamicLiveRoute{PublishIP: "127.0.0.1", PlayIP: "127.0.0.1", UrlExpireSec: 0})
	if err != nil {
		logs.Error("err=%\n", err.Error())
		return nil, err
	}

	logs.Info("ret is %#v", ret)
	videoPlayUrl := VideoPlayUrl{
		Rtmp: ret.PlayUrls.Rtmp,
		Flv:  ret.PlayUrls.Flv,
		Hls:  ret.PlayUrls.Hls,
	}
	return &videoPlayUrl, err
}

func (q *QiniuSDK) GetVideoDeviceRecordList(streamId string, start, end int, marker string, line int, format string) (*VideoRecordList, error) {
	// 查询录制记录
	recordHistory, marker, err := q.manager.QueryStreamRecordHistories(q.NsId, streamId, start, end, marker, line, format)
	if err != nil {
		logs.Error("err=%\n", err.Error())
		return nil, err
	}
	returnData := &VideoRecordList{
		Marker: marker,
		List:   &recordHistory,
	}

	return returnData, nil
}

func (q *QiniuSDK) GetVideoSpaceInfo(nsId string) (interface{}, error) {
	// 查询录制记录
	spaceInfo, err := q.manager.QueryNamespace(nsId)
	if err != nil {
		logs.Error("err=%\n", err.Error())
		return nil, err
	}
	return spaceInfo, nil
}

func (q *QiniuSDK) GetVideoStreamInfo(streamId string) (interface{}, error) {
	streamInfo, err := q.manager.QueryStream(q.NsId, streamId)
	if err != nil {
		logs.Error("err=%\n", err.Error())
		return nil, err
	}
	return streamInfo, nil
}

func (q *QiniuSDK) StartVideoRecord(streamId string) error {
	err := q.manager.StartRecord(q.NsId, streamId)
	if err != nil {
		logs.Error("err=%\n", err.Error())
		return err
	}
	logs.Info("Start VideoRecord success.")
	return nil
}

func (q *QiniuSDK) StopVideoRecord(streamId string) error {
	err := q.manager.StopRecord(q.NsId, streamId)
	if err != nil {
		logs.Error("err=%\n", err.Error())
		return err
	}
	logs.Info("Stop VideoRecord success.")
	return nil
}
