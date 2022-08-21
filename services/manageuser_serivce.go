package services

import (
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type ManageUserService struct {
}

// 通过通过用户名 获取用户信息
func (*ManageUserService) GetUserByUsername(username string) (*models.ManageUsers, error) {
	user := models.ManageUsers{Username: username}

	err := mysql.Mydb.Read(&user, "Username")

	if err == orm.ErrNoRows {
		logs.Warn("Username %s: Cannot find user!\n", username)
	} else if err != nil {
		logs.Warn(err)
	}

	return &user, err
}

// 通过通过用户名 获取用户信息
func (*ManageUserService) GetUserById(id int) (*models.ManageUsers, error) {
	user := models.ManageUsers{Id: id}

	err := mysql.Mydb.Read(&user)

	if err == orm.ErrNoRows {
		logs.Warn("Manage User Id %d: Cannot find user!\n", id)
	} else if err != nil {
		logs.Warn(err)
	}

	return &user, err
}

// 添加用户
func (*ManageUserService) Add(u *models.ManageUsers) error {

	id, err := mysql.Mydb.Insert(u)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Insert Manage User successful! ID:", id)
	return err

}

// 修改用户
func (*ManageUserService) Update(u *models.ManageUsers) error {

	id, err := mysql.Mydb.Update(u)
	if err != nil {
		logs.Warn(err)
	}
	logs.Info("Update Manage User successful! ID:", id)
	return err
}

// 修改用户
func (*ManageUserService) Delete(u *models.ManageUsers) error {
	return nil
}

// 可以支持关键字搜索查询
func (*ManageUserService) GetUserByPageAndKey(pageSize, pageNum int, keyword string) (int, int, []*models.ManageUsers, error) {

	UserData := []*models.ManageUsers{}

	totalRecord, err := mysql.Mydb.Raw("SELECT * FROM manage_users WHERE concat(ifnull(username,''), ifnull(firstname,''), ifnull(lastname, '')) like ?", "%"+keyword+"%").QueryRows(&UserData)
	if err == orm.ErrNoRows {
		logs.Warn("Get Keyword:%#v  ErrNoRows!", keyword)
		return 0, 0, nil, err
	} else if err != nil {
		logs.Warn("Get Keyword:%#v  Failed! err:%#v", keyword, err)
		return 0, 0, nil, err
	}
	totalPageNum := (int(totalRecord) + pageSize - 1) / pageSize

	num, err := mysql.Mydb.Raw("SELECT * FROM manage_users WHERE concat(ifnull(username,''), ifnull(firstname,''), ifnull(lastname, '')) like ? LIMIT ? OFFSET  ?", "%"+keyword+"%", pageSize, pageSize*(pageNum-1)).QueryRows(&UserData)

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
