package services

import (
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
)

type UserService struct {
}

// 通过通过用户名 获取用户信息
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

// 通过通过用户名 获取用户信息
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
