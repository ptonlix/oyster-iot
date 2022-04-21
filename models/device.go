package models

type Device struct {
	Id        int    `json:"id"         orm:"pk;auto;description(设备ID)"`
	Token     string `json:"token"      orm:";description(设备验证token)"`
	Protocol  string `json:"protocol"   orm:";description(设备采用的协议 MQTT TCP ...)"`
	Publish   string `json:"publish"    orm:";description(设备发布消息的主题)"`
	Subscribe string `json:"subscribe"  orm:";description(设备订阅的主题)"`
	Type      string `json:"type"       orm:"size(64);description(设备类型)"`
}
