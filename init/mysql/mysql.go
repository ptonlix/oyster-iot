package mysql

import (
	"fmt"
	"oyster-iot/models"

	"github.com/beego/beego/v2/core/logs"

	bcrypt "oyster-iot/utils"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/client/orm/filter/bean"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

var Mydb orm.Ormer

func init() {
	mysqluser, _ := beego.AppConfig.String("mysqluser")
	mysqlpasswd, _ := beego.AppConfig.String("mysqlpasswd")
	mysqladdr, _ := beego.AppConfig.String("mysqladdr")
	mysqlport, _ := beego.AppConfig.Int("mysqlport")
	mysqldb, _ := beego.AppConfig.String("mysqldb")
	mysqlMaxIdle, _ := beego.AppConfig.Int("mysqlMaxIdle")
	mysqlMaxConn, _ := beego.AppConfig.Int("mysqlMaxConn")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&loc=Asia%%2FShanghai",
		mysqluser,
		mysqlpasswd,
		mysqladdr,
		mysqlport,
		mysqldb,
	)
	fmt.Println(dataSource)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		logs.Error("Connect Mysql DataBase ERROR!")
	}
	logs.Info("Connect Mysql DataBase Success!")
	orm.SetMaxIdleConns("default", mysqlMaxIdle)
	orm.SetMaxOpenConns("default", mysqlMaxConn)
	orm.RegisterModel(new(models.Device), new(models.DeviceData), new(models.Users), new(models.Business))
	orm.RunSyncdb("default", false, true) //第二个参数是是否强制建表，true会删除数据库数据重新建表
	//orm.RunCommand() //命令模式 /main orm 显示帮助
	Mydb = orm.NewOrm()

	addDefaultUser()
}

func addDefaultUser() {
	builder := bean.NewDefaultValueFilterChainBuilder(nil, true, true)
	orm.AddGlobalFilterChain(builder.FilterChain)
	_, _ = Mydb.Insert(&models.Users{
		Id:        1,
		Username:  "admin",
		Password:  bcrypt.HashAndSalt([]byte("123456")),
		Enabled:   true,
		Email:     "260431910@qq.com",
		Firstname: "Chen",
		Lastname:  "Fudong",
		Mobile:    "13510605710",
		Remark:    "管理员账号",
		IsAdmin:   true,
		Wxopenid:  "test",
		Wxunionid: "test",
	})
}
