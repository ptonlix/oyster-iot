package services

import (
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type UserService struct {
}

// 通过用户名 获取用户信息
func (*UserService) GetUserByUsername(username string) (*models.Users, error) {
	user := models.Users{Username: username}

	err := mysql.Mydb.Read(&user, "Username")

	if err == orm.ErrNoRows {
		logs.Warn("Username %s: Cannot find user!\n", username)
	} else if err != nil {
		logs.Warn(err)
	}

	return &user, err
}

// 通过用户ID 获取用户信息
func (*UserService) GetUserById(id int) (*models.Users, error) {
	user := models.Users{Id: id}

	err := mysql.Mydb.Read(&user)

	if err == orm.ErrNoRows {
		logs.Warn("User Id %d: Cannot find user!\n", id)
	} else if err != nil {
		logs.Warn(err)
	}

	return &user, err
}

// 添加用户
func (*UserService) Add(u *models.Users) error {

	id, err := mysql.Mydb.Insert(u)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Insert User successful! ID:", id)
	return err

}

// 修改用户
func (*UserService) Update(u *models.Users) error {

	id, err := mysql.Mydb.Update(u)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Update User successful! ID:", id)
	return err
}

// 修改用户
func (*UserService) Delete(u *models.Users) error {
	// 删除设备, 同时删除设备上报的数据
	to, err := mysql.Mydb.Begin()
	if err != nil {
		logs.Error("start the transaction failed")
		return err
	}
	defer func() {
		if err != nil {
			err = to.Rollback()
			if err != nil {
				logs.Error("roll back transaction failed", err)
			}
		} else {
			err = to.Commit()
			if err != nil {
				logs.Error("commit transaction failed.", err)
			}
		}
	}()

	//删除设备数据 删除业务数据 删除设备上报的数据
	qb := to.QueryTable(&models.Business{}).Filter("user_id", u.Id)
	delid, err := qb.Delete()
	if err != nil {
		logs.Warn(err)
		return err
	}

	logs.Info("Delete Business Data num: ", delid)
	var devices []*models.Device
	qs := to.QueryTable(&models.Device{}).Filter("user_id", u.Id)
	_, err = qs.All(&devices)
	if err != nil {
		logs.Error(err)
		return err
	}

	for _, v := range devices {
		deviceData := models.DeviceData{
			DevAssetsNum: v.AssetsNum,
			DevType:      v.Type,
		}
		qs := to.QueryTable(deviceData).Filter("dev_assets_num", v.AssetsNum)
		delid, err := qs.Delete()
		if err != nil {
			logs.Warn(err)
			return err
		}
		logs.Info("Delete Device Data num: ", delid)
	}

	devnum, err := qs.Delete()
	if err != nil {
		logs.Warn(err)
		return err
	}

	logs.Info("Delete Device Data num: ", devnum)

	id, err := to.Delete(u)
	if err != nil {
		logs.Warn(err)
	}

	logs.Info("Delete User successful! ID:", id)

	return err
}

// 可以支持关键字搜索查询
func (*UserService) GetUserByPageAndKey(pageSize, pageNum int, keyword string) (int, int, []*models.Users, error) {

	UserData := []*models.Users{}

	totalRecord, err := mysql.Mydb.Raw("SELECT * FROM users WHERE concat(ifnull(username,''), ifnull(firstname,''), ifnull(lastname, '')) like ?", "%"+keyword+"%").QueryRows(&UserData)
	if err == orm.ErrNoRows {
		logs.Warn("Get Keyword:%#v  ErrNoRows!", keyword)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get Keyword:%#v  Failed! err:%#v", keyword, err)
		return 0, 0, nil, err
	}
	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	num, err := mysql.Mydb.Raw("SELECT * FROM users WHERE concat(ifnull(username,''), ifnull(firstname,''), ifnull(lastname, '')) like ? LIMIT ? OFFSET  ?", "%"+keyword+"%", pageSize, pageSize*(pageNum-1)).QueryRows(&UserData)

	if err == orm.ErrNoRows {
		logs.Warn("Get Keyword:%#v  ErrNoRows!", keyword)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get Keyword:%#v  Failed! err:%#v", keyword, err)
		return 0, 0, nil, err
	}

	logs.Info("Get Devices successful! Totalcount: %v TotalPages: %v Returned Rows Num: %#v", totalRecord, totalPageNum, num)
	return int(totalRecord), totalPageNum, UserData, err
}
