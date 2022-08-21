package services

import (
	"errors"
	"oyster-iot/init/constants"
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/qiniu/go-sdk/v7/qvs"
)

//视频监控,本地数据库操作
type VideoDbService struct {
}

func (*VideoDbService) Add(platform, nsId, username string, userId int) error {
	// 查询七牛云 QVS空间数据
	// 创建视频监控服务
	videoS, err := NewVideoService(platform, constants.QiniuAK, constants.QiniuSK, nsId) // 3nm4x17o751vq
	if err != nil {
		logs.Error(err)
		return err
	}
	videoS.SDK.Create()
	spaceInfo, err := videoS.SDK.GetVideoSpaceInfo(nsId)
	if err != nil {
		logs.Error(err)
		return err
	}
	var videoSpace *models.VideoSpace
	switch spaceInfo.(type) {
	case *qvs.NameSpace:
		videoSpace = &models.VideoSpace{
			UserId:              userId,
			Username:            username,
			Platform:            platform,
			SpaceId:             spaceInfo.(*qvs.NameSpace).ID,
			SpaceName:           spaceInfo.(*qvs.NameSpace).Name,
			AccessType:          spaceInfo.(*qvs.NameSpace).AccessType,
			Desc:                spaceInfo.(*qvs.NameSpace).Desc,
			Disabled:            spaceInfo.(*qvs.NameSpace).Disabled,
			SpaceCreatedAt:      spaceInfo.(*qvs.NameSpace).CreatedAt,
			SpaceUpdatedAt:      spaceInfo.(*qvs.NameSpace).UpdatedAt,
			DevicesCount:        spaceInfo.(*qvs.NameSpace).DevicesCount,
			StreamCount:         spaceInfo.(*qvs.NameSpace).StreamCount,
			OnlineStreamCount:   spaceInfo.(*qvs.NameSpace).OnlineStreamCount,
			DisabledStreamCount: spaceInfo.(*qvs.NameSpace).DisabledStreamCount,
			OnDemandPull:        spaceInfo.(*qvs.NameSpace).OnDemandPull,
		}
	default:
		return errors.New("spaceInfo interface convert failed!")
	}

	// 插入数据库
	id, err := mysql.Mydb.Insert(videoSpace)
	if err != nil {
		logs.Warn(err)
		return err
	}
	logs.Info("Insert Video Space successful! ID:", id)
	return err

}

// 更新一条视频空间数据
func (*VideoDbService) Update(videoSpace *models.VideoSpace) error {
	id, err := mysql.Mydb.Update(videoSpace)
	if err != nil {
		logs.Warn(err)
		return err
	}
	logs.Info("Update VideoSpace successful! ID:", id)
	return err
}

// 删除一条视频空间数据
func (*VideoDbService) Delete(videoSpace *models.VideoSpace) error {
	id, err := mysql.Mydb.Delete(videoSpace)
	if err != nil {
		logs.Warn(err)
		return err
	}

	logs.Info("Delete VideoSpace successful! ID:", id)
	return err
}

// 读取视频监控空间数据
func (*VideoDbService) Read(Id int) (*models.VideoSpace, error) {
	videoSpace := &models.VideoSpace{
		Id: Id,
	}
	err := mysql.Mydb.Read(videoSpace)
	if err == orm.ErrNoRows {
		logs.Info("Not Found VideoSpace")
		return nil, err
	} else if err == orm.ErrMissPK {
		logs.Warn("ErrMissPK")
		return nil, err
	} else {
		return videoSpace, nil
	}
}

// 读取视频监控空间数据
func (*VideoDbService) ReadByUserId(UserId int) (*models.VideoSpace, error) {
	videoSpace := &models.VideoSpace{
		UserId: UserId,
	}
	err := mysql.Mydb.Read(videoSpace, "UserId")
	if err == orm.ErrNoRows {
		logs.Info("Not Found VideoSpace")
		return nil, err
	} else if err == orm.ErrMissPK {
		logs.Warn("ErrMissPK")
		return nil, err
	} else {
		return videoSpace, nil
	}
}

// 判断该UserId是否已绑定空间
func (*VideoDbService) JudgeUser(userId int) bool {
	videoSpace := &models.VideoSpace{
		UserId: userId,
	}
	err := mysql.Mydb.Read(videoSpace)
	if err == orm.ErrNoRows {
		logs.Info("Not Found VideoSpace")
		return true
	} else if err == orm.ErrMissPK {
		logs.Warn("ErrMissPK")
		return false
	} else {
		logs.Warn("Found VideoSpace of UserId", userId)
		return false
	}
}

// 可以支持关键字搜索查询
func (*VideoDbService) GetSpacesByPageAndKey(pageSize, pageNum int, keyword string) (int, int, []*models.VideoSpace, error) {

	videoSpace := []*models.VideoSpace{}

	totalRecord, err := mysql.Mydb.Raw("SELECT * FROM video_space WHERE concat(ifnull(username,''), ifnull(space_name,'')) like ?", "%"+keyword+"%").QueryRows(&videoSpace)
	if err == orm.ErrNoRows {
		logs.Warn("Get Keyword:%#v  ErrNoRows!", keyword)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get Keyword:%#v  Failed! err:%#v", keyword, err)
		return 0, 0, nil, err
	}
	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	num, err := mysql.Mydb.Raw("SELECT * FROM video WHERE concat(ifnull(username,''), ifnull(space_name,'')) like ? LIMIT ? OFFSET  ?", "%"+keyword+"%", pageSize, pageSize*(pageNum-1)).QueryRows(&videoSpace)

	if err == orm.ErrNoRows {
		logs.Warn("Get Keyword:%#v  ErrNoRows!", keyword)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get Keyword:%#v  Failed! err:%#v", keyword, err)
		return 0, 0, nil, err
	}

	logs.Info("Get Devices successful! Totalcount: %v TotalPages: %v Returned Rows Num: %#v", totalRecord, totalPageNum, num)
	return int(totalRecord), totalPageNum, videoSpace, err
}
