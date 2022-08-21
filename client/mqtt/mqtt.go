package mqtt

import (
	"errors"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

var Subtopic map[string]string = make(map[string]string) //订阅主题
var Pubtopic map[string]string = make(map[string]string) //发布主题

var running bool
var _client mqtt.Client

func InitData() {
	Subtopic["oper"] = viper.GetString("mqtt.topicToSubscribe")
	Pubtopic["property"] = viper.GetString("mqtt.topicToPublish")
	fmt.Println(Subtopic)
	fmt.Println(Pubtopic)
}

// 生成连接平台MQTT的信息
func createMqttInfo() (string, string, string, string) {
	//获取broker和登录用户名
	mqttuser := viper.GetString("mqtt.user")
	mqttpwd := viper.GetString("mqtt.passwd")
	broker := viper.GetString("mqtt.broker")
	clientId := viper.GetString("mqtt.clientid")

	return broker, clientId, mqttuser, mqttpwd
}

func ConnectAndListen(cb func(m mqtt.Message)) {
	Connect(cb)
	Listen(Subtopic)
}

func Connect(msgProc func(m mqtt.Message)) (err error) {
	broker, clientId, mqttuser, mqttencpwd := createMqttInfo()

	running = false
	if _client == nil {
		var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
			fmt.Printf("Connect lost: %v", err)
		}
		opts := mqtt.NewClientOptions()
		opts.SetUsername(mqttuser)
		opts.SetPassword(mqttencpwd)
		opts.SetClientID(clientId)
		opts.AddBroker(broker)
		opts.SetAutoReconnect(true)
		opts.OnConnectionLost = connectLostHandler
		opts.SetOnConnectHandler(func(c mqtt.Client) {
			if !running {
				fmt.Println("MQTT CONNECT SUCCESS -- ", broker)
			}
			running = true
		})
		opts.SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
			msgProc(m)
		})
		_client = mqtt.NewClient(opts)
		if token := _client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())

		}

		for _, value := range Subtopic {
			if token := _client.Subscribe(value, 0, nil); token.Wait() &&
				token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}
		}
	}
	return
}

func Listen(subtopic map[string]string) {
	for _, value := range subtopic {
		if token := _client.Subscribe(value, 0, nil); token.Wait() &&
			token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}
}

//发送消息
func Send(payload []byte, pubtopic string) (err error) {
	var clientErr = errors.New("_client is error")
	if _client == nil {
		return clientErr
	}
	token := _client.Publish(pubtopic, 1, false, string(payload))
	if token.Error() != nil {
		fmt.Println(token.Error())
	}
	return token.Error()
}

func Close() {
	if _client != nil {
		_client.Disconnect(3000)
	}
}
