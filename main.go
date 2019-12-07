package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	_ "github.com/xinjiayu/SimServerUnicom/boot"

)

func main() {
	glog.Info("SimServerUnicom Version:", "V1.0.1.201911291930")
	g.Server().Run()
}

