package thirdpartaccess

import (
	"fmt"
	"net/url"
	"oyster-iot/services/httpclient"
)

// Manager 代表一个客户端
type Manager struct {
	ApiHost string
	client  *httpclient.Client
}

// New 初始化 Client.
func NewManager(apiHost string) *Manager {
	client := httpclient.DefaultClient
	return &Manager{
		ApiHost: apiHost,
		client:  &client,
	}
}

func setQuery(q url.Values, key string, v interface{}) {
	q.Set(key, fmt.Sprint(v))
}

func (manager *Manager) url(format string, args ...interface{}) string {
	return manager.ApiHost + fmt.Sprintf(format, args...)
}
