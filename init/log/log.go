package log

import "github.com/beego/beego/v2/core/logs"

func init() {
	logs.SetLogger(logs.AdapterConsole, `{"level":7,"color":true}`)
}
