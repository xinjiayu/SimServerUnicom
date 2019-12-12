package analyse

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/xinjiayu/SimServerUnicom/app/model/analyse"
	"github.com/xinjiayu/SimServerUnicom/app/model/datamodel"
	"github.com/xinjiayu/SimServerUnicom/app/service/operate"
	"github.com/xinjiayu/SimServerUnicom/library/utils"
)

//分析所有卡最近两天的流量情况
func TwoFlow() []analyse.TwoDaysFlow {

	fud := utils.GetFlowUseDate()

	carrier := "unicom" //sim卡运营商

	var twoDaysFlowList []analyse.TwoDaysFlow
	//本月流量
	monthFlowList := datamodel.SimFlow{}.FlowList("year=? and month=? and carrier=?", fud.Year, fud.Month, carrier)
	//如果当前是1月份，后面取上一年的数据

	year := fud.Year
	if fud.Month == "1" {
		year = fud.LastYear
	}
	//上个月流量
	lastMonthFlowList := datamodel.SimFlow{}.FlowList("year=? and month=? and carrier=?", year, fud.LastMonth, carrier)
	var lastMonthFlowSimList = make(map[string]datamodel.SimFlow)
	for _, v := range lastMonthFlowList {
		lastMonthFlowSimList[v.Iccid] = v
	}

	for _, v := range monthFlowList {

		twoDaysFlow := analyse.TwoDaysFlow{}
		twoDaysFlow.Iccid = v.Iccid

		flowData := gconv.Map(v)
		lastMonthFlowSim := gconv.Map(lastMonthFlowSimList[v.Iccid])

		todayNum := gconv.Int(fud.Today)

		var d0 int64 = 0
		var d1 int64 = 0
		var d2 int64 = 0
		d0 = gconv.Int64(flowData["d"+fud.Today])
		d1 = gconv.Int64(flowData["d"+fud.Yesterday])
		d2 = gconv.Int64(flowData["d"+fud.BeforeYesterday])

		switch todayNum {
		case 1:
			d1 = gconv.Int64(lastMonthFlowSim["d"+fud.Yesterday])
			d2 = gconv.Int64(lastMonthFlowSim["d"+fud.BeforeYesterday])
		case 2:
			d2 = gconv.Int64(lastMonthFlowSim["d"+fud.BeforeYesterday])
		}

		//今天的流量
		twoDaysFlow.TodayUsage = (d0 - d1) / utils.MB1
		if fud.Today == "27" {
			twoDaysFlow.TodayUsage = d0 / utils.MB1 //如果是27号，不计算
		}

		if twoDaysFlow.TodayUsage < 0 {
			twoDaysFlow.TodayUsage = 0
		}
		twoDaysFlow.YesterdayUsage = (d1 - d2) / utils.MB1 //昨天的流量
		if twoDaysFlow.YesterdayUsage < 0 {
			twoDaysFlow.YesterdayUsage = 0
		}

		twoDaysFlowList = append(twoDaysFlowList, twoDaysFlow)
	}
	return twoDaysFlowList
}

//分析所有卡的流量列表，可以指定计费套餐
func SimList(panName string) []analyse.PlanSimCardInfo {

	monthFlowList := datamodel.SimUnicom{}.FlowList(panName)
	var SimFlowList []analyse.PlanSimCardInfo
	for _, v := range monthFlowList {
		simFlow := analyse.PlanSimCardInfo{}

		simFlow.Iccid = v.Iccid
		simFlow.Flow = v.CtdDataUsage
		simFlow.PlanName = v.RatePlan
		SimFlowList = append(SimFlowList, simFlow)
	}

	return SimFlowList

}

//分析计费套餐使用情况
func PlanInfoList() []analyse.PlanInfo {

	simCardList, _ := datamodel.SimUnicom{}.GetUnicomSimInfoList()
	var op = new(operate.AutoChangePlan)
	op.CountPlanFlow(simCardList)

	var planInfoList []analyse.PlanInfo
	planInfo1 := op.PlanInfo[datamodel.Plan01]
	if planInfo1.PlanName != "" {
		planInfoList = append(planInfoList, planInfo1)

	}
	planInfo2 := op.PlanInfo[datamodel.Plan02]
	if planInfo2.PlanName != "" {
		planInfoList = append(planInfoList, planInfo2)

	}
	planInfo3 := op.PlanInfo[datamodel.Plan03]
	if planInfo3.PlanName != "" {
		planInfoList = append(planInfoList, planInfo3)

	}

	return planInfoList
}

//分析指定sim卡的月度流量数据
func MonthSimFlowListByIccid(iccid string) analyse.TwoDaysBaySimCardFlow {
	if iccid == "" {
		return analyse.TwoDaysBaySimCardFlow{}
	}
	fud := utils.GetFlowUseDate()

	carrier := "unicom" //sim卡运营商
	//本月流量
	monthFlow := datamodel.SimFlow{}.GetByOne("year=? and month=? and carrier=? and iccid=?", fud.Year, fud.Month, carrier, iccid)
	//如果当前是1月份，后面取上一年的数据
	year := fud.Year
	if fud.Month == "1" {
		year = fud.LastYear
	}
	//上个月流量
	lastMonthFlow := datamodel.SimFlow{}.GetByOne("year=? and month=? and carrier=? and iccid=?", year, fud.LastMonth, carrier, iccid)

	//上上个月流量
	BeforeLastMonthFlow := datamodel.SimFlow{}.GetByOne("year=? and month=? and carrier=? and iccid=?", year, fud.BeforeLastMonth, carrier, iccid)

	var twoDaysBaySimCardFlow analyse.TwoDaysBaySimCardFlow

	monthFlowData := gconv.Map(monthFlow)
	LastMonthFlowData := gconv.Map(lastMonthFlow)
	BeforeLastMonthFlowData := gconv.Map(BeforeLastMonthFlow)

	lastMonthFlowData := gconv.Int64(LastMonthFlowData["d"+fud.LastMonthDays])                   //上个月最后一天的数据
	beforeLastMonthFlowData := gconv.Int64(BeforeLastMonthFlowData["d"+fud.BeforeLastMonthDays]) //上上个月最后一天的数据

	for i := 1; i < 32; i++ {
		//输出每天的数字
		twoDaysBaySimCardFlow.DayName = append(twoDaysBaySimCardFlow.DayName, i)
	}

	//本月流量统计
	twoDaysBaySimCardFlow.Nowmonth = countDaysFlow(lastMonthFlowData, monthFlowData)
	//上个月流量统计
	twoDaysBaySimCardFlow.Beforemonth = countDaysFlow(beforeLastMonthFlowData, LastMonthFlowData)

	return twoDaysBaySimCardFlow

}

//分析每天使用的流量
func countDaysFlow(lastMonthFlowData int64, monthFlowData map[string]interface{}) []int64 {
	var dayFlows []int64
	for i := 1; i < 32; i++ {
		var d0 int64 = 0
		var d1 int64 = 0
		d0 = gconv.Int64(monthFlowData["d"+gconv.String(i)])
		d1 = gconv.Int64(monthFlowData["d"+gconv.String(i-1)])
		if i == 1 {
			d1 = lastMonthFlowData
		}
		//本天的流量
		dayFlow := d0 - d1
		if i == 27 {
			dayFlow = d0 //如果是27号，不计算
		}
		if dayFlow < 0 {
			dayFlow = 0
		}
		dayFlows = append(dayFlows, dayFlow/utils.MB1)
	}
	return dayFlows
}

//通过事件类型提取通知信息
func GetNotice(eventType string) []datamodel.UnicomNotice {
	fud := utils.GetFlowUseDate()
	startDate := fud.Year + "-" + fud.LastMonth + "-27"
	if gconv.Int(fud.Month) == 1 && gconv.Int(fud.Today) < 27 {
		startDate = fud.LastYear + "-" + fud.LastMonth + "-27"
	}
	if gconv.Int(fud.Today) > 26 {
		startDate = fud.Year + "-" + fud.Month + "-27"
	}
	endDate := gtime.Now().Format("Y-m-d")

	startTime := gtime.NewFromStr(startDate).Second()
	endTime := gtime.NewFromStr(endDate).Second()

	if eventType == "" {
		return datamodel.UnicomNotice{}.List("timestamp > ? and timestamp < ?", startTime, endTime)
	} else {
		return datamodel.UnicomNotice{}.List("timestamp > ? and timestamp < ? and event_type = ?", startTime, endTime, eventType)

	}
}
