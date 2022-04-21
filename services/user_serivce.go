package services

import (
	"log"
	"oyster-iot/init/mysql"
	"oyster-iot/models"

	"github.com/beego/beego/v2/client/orm"
)

type UserService struct {
}

// 通过通过用户名 获取用户信息
func (*UserService) GetUserByUsername(username string) (*models.Users, error) {
	user := models.Users{Username: username}

	err := mysql.Mydb.Read(&user, "Username")

	if err == orm.ErrNoRows {
		log.Printf("Username %s: Cannot find user!\n", username)
	} else if err != nil {
		log.Println(err)
	}

	return &user, err
}

// 通过通过用户名 获取用户信息
func (*UserService) GetUserById(id int) (*models.Users, error) {
	user := models.Users{Id: id}

	err := mysql.Mydb.Read(&user)

	if err == orm.ErrNoRows {
		log.Printf("User Id %d: Cannot find user!\n", id)
	} else if err != nil {
		log.Println(err)
	}

	return &user, err
}
