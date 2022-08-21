package datahandler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"oyster-iot/client/mqtt"
	"strconv"
	"time"

	mq "github.com/eclipse/paho.mqtt.golang"

	"github.com/spf13/viper"
)

// 模拟设备发送数据
type PropertyData struct {
	AssetsNum string      `json:"assets_num"`
	Token     string      `json:"token"`
	Msg       interface{} `json:"msg"`
}

// 模拟设备发送数据
type TempData struct {
	Temperature float64 `json:"temperature"`
}

type SaltData struct {
	Salinity float64 `json:"salinity"`
}

// 接受数据格式
type PayloadData struct {
	AssetsNum string `json:"assets_num"`
	Token     string `json:"token"`
	Cmd       string `json:"cmd"`
}

func MqttMsgCallback(m mq.Message) {
	fmt.Println("recv topic: ", m.Topic())
	fmt.Println("recv data: ", string(m.Payload()))
	recvData(m.Payload())
}

func recvData(payload []byte) {
	P := PayloadData{}
	err := json.Unmarshal(payload, &P)
	if err != nil {
		fmt.Println("json Unmarshal error!")
		return
	}
	if P.Cmd == "refresh_temp" {
		recvTemp(&P)
	} else if P.Cmd == "refresh_salt" {
		recvSalt(&P)
	}

}

func recvTemp(p *PayloadData) {
	tempdata := TempData{Temperature: randFloats(10, 30)}
	sendData(p.AssetsNum, p.Token, tempdata)
}

func recvSalt(p *PayloadData) {
	saltdata := SaltData{Salinity: randFloats(28, 35)}
	sendData(p.AssetsNum, p.Token, saltdata)
}

func sendData(assetsNum, token string, payload interface{}) {
	rand.Seed(time.Now().UnixNano())

	data := PropertyData{assetsNum, token, payload}

	senddata, err := json.Marshal(data)
	if err != nil {
		fmt.Println(" Data  Format ERROR : ", err)
		return
	}
	fmt.Println("senddata :", string(senddata))
	err = mqtt.Send(senddata, mqtt.Pubtopic["property"])
	if err != nil {
		fmt.Println(" Send Data ERROR : ", err)
		return
	}
	return
}

func DataHandler() {

	tempClientNum, err := strconv.ParseInt(viper.GetString("tempclientNum.count"), 10, 64)
	if err != nil {
		fmt.Println("Param tempClientNum Error! APP exit")
		os.Exit(1)
	}
	saltClientNum, err := strconv.ParseInt(viper.GetString("saltclientNum.count"), 10, 64)
	if err != nil {
		fmt.Println("Param saltclientNum Error! APP exit")
		os.Exit(1)
	}
	fmt.Println("temp Client Num : ", tempClientNum)

	var tempClient = make(map[string]string)

	for i := 1; int64(i) <= tempClientNum; i++ {
		tempAssertsNum := viper.GetString(fmt.Sprintf("tempclient%d.assetsNum", i))
		tempToken := viper.GetString(fmt.Sprintf("tempclient%d.token", i))
		tempClient[tempAssertsNum] = tempToken

		fmt.Println("AssertNum : ", tempAssertsNum)
		fmt.Println("tempToken : ", tempToken)
	}

	fmt.Println("salt Client Num : ", saltClientNum)
	var saltClient = make(map[string]string)

	for i := 1; int64(i) <= saltClientNum; i++ {
		saltAssertsNum := viper.GetString(fmt.Sprintf("saltclient%d.assetsNum", i))
		saltToken := viper.GetString(fmt.Sprintf("saltclient%d.token", i))
		saltClient[saltAssertsNum] = saltToken

		fmt.Println("AssertNum : ", saltAssertsNum)
		fmt.Println("tempToken : ", saltToken)
	}

	sendtime, err := strconv.ParseInt(viper.GetString("SendTimer.time"), 10, 64)
	if err != nil {
		fmt.Println("Param SendTimer Error! APP exit")
		os.Exit(1)
	}
	fmt.Printf("SendTime : %d s\n", sendtime)
	timeDuration := time.Second * time.Duration(sendtime)

	TempTimer := time.NewTimer(timeDuration) // 启动定时器
	SaltTimer := time.NewTimer(timeDuration) // 启动定时器

	for {
		select {
		case <-TempTimer.C:
			for assetsNum, token := range tempClient {
				tempdata := TempData{Temperature: randFloats(20, 30)}
				sendData(assetsNum, token, tempdata)
			}
			TempTimer.Reset(timeDuration) // 每次使用完后需要人为重置下

		case <-SaltTimer.C:
			for assetsNum, token := range saltClient {
				saltdata := SaltData{Salinity: randFloats(28, 35)}
				sendData(assetsNum, token, saltdata)
			}
			SaltTimer.Reset(timeDuration) // 每次使用完后需要人为重置下
		}
	}
	// 不再使用了，结束它
	TempTimer.Stop()
	SaltTimer.Stop()
}

func randFloats(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	res := min + rand.Float64()*(max-min)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", res), 64)
	return value
}
