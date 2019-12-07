package collect

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	collectService "github.com/xinjiayu/SimServerUnicom/app/service/collect"
	"github.com/xinjiayu/SimServerUnicom/library/response"
)


// Controller API管理对象
type Controller struct{}

// CtdUsages 获取指定sim卡的流量，返回指定设备的周期累计用量信息。
func (c *Controller) CtdUsages(r *ghttp.Request) {
	iccid := r.GetString("iccid")
	dataModel :=collectService.GetUsages(iccid)
	response.Json(r, 0, "", dataModel)

}

// CardList 通过联通平台的devices接口获取设备卡列表
func (c *Controller) CardList(r *ghttp.Request) {

	glog.Info("\n开始获取SIM卡列表数据...")
	simList := collectService.GetCardList()
	number := len(simList.Devices)
	numstr := gconv.String(number)
	glog.Infof("SIM卡共%d张卡", number)
	response.Json(r, 0, "SIM卡流量数据条数："+numstr, simList.Devices)

}

