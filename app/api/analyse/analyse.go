// 联通物联网卡流量数据分析API
//
// 对从联通物联网平台同步回来的数据进行统计与分析。
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta
package analyse

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	serviceAnalyse "github.com/xinjiayu/SimServerUnicom/app/service/analyse"
	"github.com/xinjiayu/SimServerUnicom/library/response"
)

// Controller API管理对象
type Controller struct{}

// TwoDaysSimCardFlow 获取最近两天流量的sim卡列表
func (c *Controller) TwoDaysSimCardFlow(r *ghttp.Request) {
	rsData := serviceAnalyse.TwoFlow()
	numStr := gconv.String(len(rsData))
	response.Json(r, 0, "SIM卡流量数据条数："+numStr, rsData)
}

// AllSimFlowList 获取计费周期内所有sim卡用量信息
// 参数：planName 计划名称,默可以为空，输出全部数据
func (c *Controller) AllSimFlowList(r *ghttp.Request) {
	planName := r.GetString("planName")
	rsData := serviceAnalyse.SimList(planName)
	numStr := gconv.String(len(rsData))
	response.Json(r, 0, "SIM卡流量数据条数："+numStr, rsData)
}

// PlanInfoCountList 获取计费套餐计划的统计信息
func (c *Controller) PlanInfoCountList(r *ghttp.Request) {
	rsData := serviceAnalyse.PlanInfoList()
	numStr := gconv.String(len(rsData))
	response.Json(r, 0, "SIM卡流量套餐计划数："+numStr, rsData)
}

// MonthSimFlowByIccid 获取指定sim卡最近两个月的流量
// 参数：iccid ，sim卡的iccid号
func (c *Controller) MonthSimFlowByIccid(r *ghttp.Request) {
	iccid := r.GetString("iccid")
	rsData := serviceAnalyse.MonthSimFlowListByIccid(gconv.String(iccid))
	response.Json(r, 0, iccid+": SIM卡，流量信息：", rsData)
}

// MonthSimFlowCount 获取所有sim卡最近两个月的流量
func (c *Controller) MonthSimFlowCount(r *ghttp.Request) {
	rsData := serviceAnalyse.MonthSimFlowList()
	response.Json(r, 0, ": SIM月度每日流量信息：", rsData)
}

// Notice 获取通知信息
// 参数：event_type，字符类型，默认为空，将显示本周期内所有的通知信息。当前事件类型支持：PAST24H_DATA_USAGE_EXCEEDED 24小时内流量超过指定量的通知、CTD_USAGE 周期使用超过指定的通知
func (c *Controller) Notice(r *ghttp.Request) {
	eventType := r.GetString("event_type")
	rsData := serviceAnalyse.GetNotice(eventType)
	if eventType==""{
		eventType = "全部"
	}
	response.Json(r, 0, "通知信息，类型："+eventType, rsData)

}
