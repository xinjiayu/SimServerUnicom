package operate

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	serviceOperate "github.com/xinjiayu/SimServerUnicom/app/service/operate"
	"github.com/xinjiayu/SimServerUnicom/library/response"
)

// Controller API管理对象
type Controller struct{}

// ChangePlan 跟据sim卡流量池使用情况进行自动设置平衡，平衡的顺序为1G池超出时，自动变更部分sim卡到2G池中，
//当2G池流量超出时自动将部分sim卡变更到3G池。当3G池也超出时将1G池中部分sim卡变更到3G池
func (c *Controller) ChangePlan(r *ghttp.Request) {

	as := new(serviceOperate.AutoChangePlan)
	err, num := as.AutoSetupPlan()
	if err != nil {
		response.Json(r, -1, err.Error(), "")

	}
	response.Json(r, 0, "SIM卡计费套餐变动数："+gconv.String(num), "")

}

// ChangeInitPlan 自动初始化所有sim卡的流量套餐变更为1G套餐
func (c *Controller) ChangeInitPlan(r *ghttp.Request) {

	as := new(serviceOperate.AutoChangePlan)
	if as.GetNot01PlanNum() > 0 {
		err, num := as.AutoSetupPlanInit()
		if err != nil {
			response.Json(r, -1, err.Error(), "")

		}
		response.Json(r, 0, "SIM卡计费套餐变动数："+gconv.String(num), "")
	}

}

