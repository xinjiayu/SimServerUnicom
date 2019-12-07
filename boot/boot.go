package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

// 用于应用初始化。
func init() {
	c := g.Config()
	s := g.Server()
	routerInit()

	// log配置
	logPath := c.GetString("system.logPath")
	glog.SetPath(logPath)
	glog.SetStdoutPrint(true)

	// Web Server配置
	s.SetLogPath(logPath)
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)

	port := c.GetInt("system.port")
	s.SetPort(port)
}
