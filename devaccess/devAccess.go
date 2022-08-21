package devaccess

import (
	"flag"
	"log"
	devmq "oyster-iot/devaccess/modules/mqtt"
	"oyster-iot/services"
	"oyster-iot/utils"

	"github.com/beego/beego/v2/core/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

func init() {
	loadConfig()
	listenMQTT()
}

func loadConfig() {
	logs.Info("read devaccess config")
	envConfigFile := flag.String("config", "./devaccess/config.ini", "The path of device access layer config file")
	flag.Parse()
	viper.SetConfigFile(*envConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("FAILURE", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			logs.Warn("devaccess config Not Found!")
		} else {
			// 配置文件被找到，但产生了另外的错误
			logs.Error("devaccess config Found, but ERROR !")
		}
	}
}

func listenMQTT() {
	var devDataS services.DevdataSevice
	// 初始化InfluxDB配置
	devDataS.Host = viper.GetString("influxDB.influxhost")
	devDataS.Token = viper.GetString("influxDB.influxtoken")
	devDataS.Org = viper.GetString("influxDB.influxorg")
	devDataS.Bucket = viper.GetString("influxDB.influxbucket")

	logs.Info("InfluxDB Config: Host: %s | Token: %s | Organization: %s | Bucket: %s|",
		devDataS.Host, devDataS.Token, devDataS.Org, devDataS.Bucket)
	devDataS.Online = utils.DetectInfluxDBOnline(devDataS.Host, devDataS.Token, devDataS.Org)
	broker := viper.GetString("mqtt.broker")
	clientid := viper.GetString("mqtt.clientid")
	user := viper.GetString("mqtt.user")
	pass := viper.GetString("mqtt.passwd")
	devmq.Listen(broker, user, pass, clientid, func(m mqtt.Message) {
		devDataS.MQTTMsgProc(m.Payload())
	})

}
