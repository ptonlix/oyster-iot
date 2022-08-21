package main

import (
	"flag"
	"fmt"
	"oyster-iot/client/datahandler"
	"oyster-iot/client/mqtt"

	"github.com/spf13/viper"
)

func init() {
	//读取配置配置文件
	loadConfig()
}

//读取配置文件
func loadConfig() {
	fmt.Println("read devaccess config")
	envConfigFile := flag.String("config", "./config.ini", "The path of config file")
	viper.SetConfigFile(*envConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("FAILURE", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			fmt.Println("devaccess config Not Found!")
		} else {
			// 配置文件被找到，但产生了另外的错误
			fmt.Println("devaccess config Found, but ERROR !")
		}
	}
}

func main() {
	fmt.Println("Oyster Client Start...")
	mqtt.InitData()
	mqtt.ConnectAndListen(datahandler.MqttMsgCallback)
	datahandler.DataHandler()
}
