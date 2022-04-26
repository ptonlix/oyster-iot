package devaccess

import (
	"flag"
	"log"
	devmq "oyster-iot/devaccess/modules/mqtt"
	"oyster-iot/services"

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
	broker := viper.GetString("mqtt.broker")
	clientid := viper.GetString("mqtt.clientid")
	user := viper.GetString("mqtt.user")
	pass := viper.GetString("mqtt.passwd")
	devmq.Listen(broker, user, pass, clientid, func(m mqtt.Message) {
		devDataS.MQTTMsgProc(m.Payload())
	})
}
