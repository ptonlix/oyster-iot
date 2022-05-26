package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/spf13/viper"
)

const timeout = 5 * time.Second

type EmqExpService struct {
}

var (
	targetsV2 = map[string]string{
		"monitoring_metrics": "/api/v2/monitoring/metrics/%s",
		"monitoring_stats":   "/api/v2/monitoring/stats/%s",
		"monitoring_nodes":   "/api/v2/monitoring/nodes/%s",
		"management_nodes":   "/api/v2/management/nodes/%s",
	}
	//scraping endpoints for EMQ v3 api version
	targetsV3 = map[string]string{
		"nodes_metrics": "/api/v3/nodes/%s/metrics/",
		"nodes_stats":   "/api/v3/nodes/%s/stats/",
		"nodes":         "/api/v3/nodes/%s",
	}
	targetsV4 = map[string]string{
		"nodes_metrics": "/api/v4/nodes/%s/metrics/",
		"nodes_stats":   "/api/v4/nodes/%s/stats/",
		"nodes":         "/api/v4/nodes/%s",
	}
)

func (*EmqExpService) GetEmqMetrics() (map[string]interface{}, error) {
	// 读取配置文件信息
	emqhost := viper.GetString("emqexp.host")
	emqnode := viper.GetString("emqexp.node")
	emquser := viper.GetString("emqexp.user")
	emqpasswd := viper.GetString("emqexp.passwd")
	emqapiVersion := viper.GetString("emqexp.appVersion")
	logs.Debug("EMQX Info :", emqhost, emqnode, emquser, emqpasswd, emqapiVersion)

	client := NewClient(emqhost, emqnode, emqapiVersion, emquser, emqpasswd)

	data, err := client.Fetch()
	if err != nil {
		logs.Warn("Get EMQX Mertrics Failed!")
		return nil, err
	}
	return data, nil
}

type emqResponse struct {
	Code   float64                `json:"code,omitempty"`
	Result map[string]interface{} `json:"result,omitempty"` //api v2 json key
	Data   map[string]interface{} `json:"data,omitempty"`
}

//Client manages communication with emq api
type Client struct {
	hc         *http.Client
	host       string
	node       string
	apiVersion string
	targets    map[string]string
	username   string
	password   string
}

func NewClient(host, node, apiVersion, username, password string) *Client {

	c := &Client{
		hc:         &http.Client{Timeout: timeout},
		host:       host,
		node:       node,
		username:   username,
		password:   password,
		targets:    targetsV4,
		apiVersion: apiVersion,
	}

	switch apiVersion {
	case "v2":
		c.targets = targetsV2
	case "v3":
		c.targets = targetsV3
	case "v4":
		c.targets = targetsV4
	}

	return c
}

//newRequest creates a new http request, setting the relevant headers
func (c *Client) newRequest(path string) (req *http.Request, err error) {

	u := c.host + fmt.Sprintf(path, c.node)

	if !strings.Contains(u, "://") {
		u = fmt.Sprintf("http://%s", u)
	}

	logs.Debug("Fetching from " + u)

	req, err = http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		logs.Debug("Failed to create http request: " + err.Error())
		return req, fmt.Errorf("Failed to create http request: %v", err)
	}

	//set request headers
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")
	return
}

//get preforms an http GET call to the provided path and returns the response
func (c *Client) get(path string) (map[string]interface{}, error) {

	req, err := c.newRequest(path)
	if err != nil {
		return nil, err
	}

	er := &emqResponse{}
	data := make(map[string]interface{})

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to get metrics: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received status code not ok %s, got %d", req.URL, res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(er); err != nil {
		return nil, fmt.Errorf("Error in json decoder %v", err)
	}

	if er.Code != 0 {
		return nil, fmt.Errorf("Recvied code != 0 from EMQ %f", er.Code)
	}

	//Print the returned response data for debuging
	logs.Debug("%#v", *er)

	if c.apiVersion == "v2" {
		data = er.Result
	} else {
		data = er.Data
	}
	return data, nil
}

//Fetch gets all the metrics from the emq api listed in the targets map
//implements emq_exporter.Fetcher
func (c *Client) Fetch() (map[string]interface{}, error) {

	data := make(map[string]interface{})

	for name, path := range c.targets {

		res, err := c.get(path)
		if err != nil {
			logs.Error(err.Error())
			return nil, err
		}

		for k, v := range res {
			mName := fmt.Sprintf("%s_%s", name, strings.Replace(k, ".", "_", -1))
			data[mName] = v
		}
	}

	return data, nil
}

/*
		获取的数据说明：参考：https://www.emqx.io/docs/zh/v4.4/advanced/http-api.html#%E7%BB%9F%E8%AE%A1%E6%8C%87%E6%A0%87
        "nodes_connections": 2, //当前接入此节点的客户端数量 --
        "nodes_load1": "0.02", //1 分钟内的 CPU 平均负载  --
        "nodes_load15": "0.18", //5 分钟内的 CPU 平均负载 --
        "nodes_load5": "0.17", //15 分钟内的 CPU 平均负载 --
        "nodes_max_fds": 1048576, //操作系统的最大文件描述符限制
        "nodes_memory_total": "142.77M", //VM 已分配的系统内存
        "nodes_memory_used": "96.02M", //VM 已占用的内存大小 --
        "nodes_metrics_bytes.received": 16718,						//EMQX 接收的字节数 EMQ
        "nodes_metrics_bytes.sent": 13883,							//EMQX 在此连接上发送的字节数 EMQ
        "nodes_metrics_client.acl.allow": 71,
        "nodes_metrics_client.acl.cache_hit": 10,
        "nodes_metrics_client.acl.deny": 0,
        "nodes_metrics_client.auth.anonymous": 48,					//匿名登录的客户端数量
        "nodes_metrics_client.authenticate": 48,					//客户端认证次数
        "nodes_metrics_client.check_acl": 61,						//ACL 规则检查次数
        "nodes_metrics_client.connack": 48,							//发送 CONNACK 报文的次数
        "nodes_metrics_client.connect": 48,							//客户端连接次数
        "nodes_metrics_client.connected": 48,						//客户端成功连接次数
        "nodes_metrics_client.disconnected": 46,					//客户端断开连接次数
        "nodes_metrics_client.subscribe": 51,                 		//客户端订阅次数
        "nodes_metrics_client.unsubscribe": 0,						//客户端取消订阅次数
        "nodes_metrics_delivery.dropped": 0,						//发送时丢弃的消息总数
        "nodes_metrics_delivery.dropped.expired": 0,
        "nodes_metrics_delivery.dropped.no_local": 0,
        "nodes_metrics_delivery.dropped.qos0_msg": 0,
        "nodes_metrics_delivery.dropped.queue_full": 0,
        "nodes_metrics_delivery.dropped.too_large": 0,
        "nodes_metrics_messages.acked": 0,
        "nodes_metrics_messages.delayed": 0,//EMQX 存储的延迟发布的消息数量
        "nodes_metrics_messages.delivered": 28,//EMQX 内部转发到订阅进程的消息数量
        "nodes_metrics_messages.dropped": 1,//EMQX 内部转发到订阅进程前丢弃的消息总数
        "nodes_metrics_messages.dropped.await_pubrel_timeout": 0,
        "nodes_metrics_messages.dropped.no_subscribers": 1,
        "nodes_metrics_messages.forward": 0,
        "nodes_metrics_messages.publish": 20,
        "nodes_metrics_messages.qos0.received": 9,
        "nodes_metrics_messages.qos0.sent": 28,
        "nodes_metrics_messages.qos1.received": 11,
        "nodes_metrics_messages.qos1.sent": 0,
        "nodes_metrics_messages.qos2.received": 0,
        "nodes_metrics_messages.qos2.sent": 0,
        "nodes_metrics_messages.received": 20, //接收来自客户端的消息数量
        "nodes_metrics_messages.retained": 12468,
        "nodes_metrics_messages.sent": 28, //发送给客户端的消息数量
        "nodes_metrics_packets.auth.received": 0, //接收的报文数量
        "nodes_metrics_packets.auth.sent": 0,	//发送的报文数量
        "nodes_metrics_packets.connack.auth_error": 0,
        "nodes_metrics_packets.connack.error": 0,
        "nodes_metrics_packets.connack.sent": 48,
        "nodes_metrics_packets.connect.received": 48,
        "nodes_metrics_packets.disconnect.received": 0,
        "nodes_metrics_packets.disconnect.sent": 0,
        "nodes_metrics_packets.pingreq.received": 5634,
        "nodes_metrics_packets.pingresp.sent": 5634,
        "nodes_metrics_packets.puback.inuse": 0,
        "nodes_metrics_packets.puback.missed": 0,
        "nodes_metrics_packets.puback.received": 0,
        "nodes_metrics_packets.puback.sent": 11,
        "nodes_metrics_packets.pubcomp.inuse": 0,
        "nodes_metrics_packets.pubcomp.missed": 0,
        "nodes_metrics_packets.pubcomp.received": 0,
        "nodes_metrics_packets.pubcomp.sent": 0,
        "nodes_metrics_packets.publish.auth_error": 0,
        "nodes_metrics_packets.publish.dropped": 0,
        "nodes_metrics_packets.publish.error": 0,
        "nodes_metrics_packets.publish.inuse": 0,
        "nodes_metrics_packets.publish.received": 20,
        "nodes_metrics_packets.publish.sent": 28,
        "nodes_metrics_packets.pubrec.inuse": 0,
        "nodes_metrics_packets.pubrec.missed": 0,
        "nodes_metrics_packets.pubrec.received": 0,
        "nodes_metrics_packets.pubrec.sent": 0,
        "nodes_metrics_packets.pubrel.missed": 0,
        "nodes_metrics_packets.pubrel.received": 0,
        "nodes_metrics_packets.pubrel.sent": 0,
        "nodes_metrics_packets.received": 5753,
        "nodes_metrics_packets.sent": 5772,
        "nodes_metrics_packets.suback.sent": 51,
        "nodes_metrics_packets.subscribe.auth_error": 0,
        "nodes_metrics_packets.subscribe.error": 0,
        "nodes_metrics_packets.subscribe.received": 51,
        "nodes_metrics_packets.unsuback.sent": 0,
        "nodes_metrics_packets.unsubscribe.error": 0,
        "nodes_metrics_packets.unsubscribe.received": 0,
        "nodes_metrics_session.created": 48,
        "nodes_metrics_session.discarded": 0,
        "nodes_metrics_session.resumed": 0,
        "nodes_metrics_session.takeovered": 0,
        "nodes_metrics_session.terminated": 46,
        "nodes_node": "8de05ead9c87@172.17.0.3",//节点名称 --
        "nodes_node_status": "Running", //节点状态 --
        "nodes_otp_release": "24.1.5/12.1.5",
        "nodes_process_available": 2097152,
        "nodes_process_used": 441,
        "nodes_stats_channels.count": 2,
        "nodes_stats_channels.max": 2,
        "nodes_stats_connections.count": 2, //当前连接数量 EMQ
        "nodes_stats_connections.max": 2, //连接数量的历史最大值EMQ
        "nodes_stats_live_connections.count": 2,
        "nodes_stats_live_connections.max": 2,
        "nodes_stats_retained.count": 3,
        "nodes_stats_retained.max": 3,
        "nodes_stats_routes.count": 3,
        "nodes_stats_routes.max": 3,
        "nodes_stats_sessions.count": 2,
        "nodes_stats_sessions.max": 2,
        "nodes_stats_suboptions.count": 4,
        "nodes_stats_suboptions.max": 4,
        "nodes_stats_subscribers.count": 4, //当前订阅者数量 EMQ
        "nodes_stats_subscribers.max": 4, //订阅者数量的历史最大值EMQ
        "nodes_stats_subscriptions.count": 4,
        "nodes_stats_subscriptions.max": 4,
        "nodes_stats_subscriptions.shared.count": 0,
        "nodes_stats_subscriptions.shared.max": 0,
        "nodes_stats_topics.count": 3, //当前主题数量 EMQ
        "nodes_stats_topics.max": 3, //主题数量的历史最大值EMQ
        "nodes_uptime": "2 days, 22 hours, 32 minutes, 46 seconds", //EMQX 运行时间 --
        "nodes_version": "4.4.3"
*/
