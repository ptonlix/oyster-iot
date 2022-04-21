package mysql

import (
	"fmt"
	"oyster-iot/models"

	"log"

	"github.com/beego/beego/v2/client/orm"
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
		log.Fatalln("Connect Mysql DataBase ERROR!")
	}
	log.Println("Connect Mysql DataBase Success!")
	orm.SetMaxIdleConns("default", mysqlMaxIdle)
	orm.SetMaxOpenConns("default", mysqlMaxConn)
	orm.RegisterModel(new(models.Device), new(models.DeviceData))
	orm.RunSyncdb("default", false, true) //第二个参数是是否强制建表，true会删除数据库数据重新建表

	Mydb = orm.NewOrm()
}
