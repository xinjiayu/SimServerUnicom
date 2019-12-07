package operate

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/xinjiayu/SimServerUnicom/app/model/analyse"
	"github.com/xinjiayu/SimServerUnicom/app/model/datamodel"
	"github.com/xinjiayu/SimServerUnicom/app/model/unicommodel"
	"github.com/xinjiayu/SimServerUnicom/app/service/collect"
	"github.com/xinjiayu/SimServerUnicom/library/utils"
	"time"
)

const plan01 = "831WLW016555_MON-FLEX_1024M_SP"
const plan02 = "831WLW016555_MON-FLEX_2048M_SP"
const plan03 = "831WLW016555_MON-FLEX_3072M_SP"

type AutoChangePlan struct {
	SimCardLIst   []analyse.PlanSimCardInfo
	SimPlan01List map[string]analyse.PlanSimCardInfo
	SimPlan02List map[string]analyse.PlanSimCardInfo
	SimPlan03List map[string]analyse.PlanSimCardInfo
	PlanList      map[string]analyse.PlanInfo
}


//AutoSetupInit 自动初始化sim卡的设置 ，为1g套餐
func (ac *AutoChangePlan) AutoSetupPlanInit() (error, int) {
	num := 0
	simCardList, err := datamodel.SimUnicom{}.GetUnicomSimInfoList()
	for _, v := range simCardList {
		simInfo := new(unicommodel.SimInfo)
		v.Struct(simInfo)
		if simInfo.RatePlan != plan01 {
			glog.Info(simInfo.Iccid, plan01)
			go startToChangePlan(simInfo.Iccid, plan01)
			num++
		}
	}
	return err, num
}

func (ac *AutoChangePlan) GetNot01PlanNum() int {
	num := 0
	simCardList, _ := datamodel.SimUnicom{}.GetUnicomSimInfoList()
	for _, v := range simCardList {
		simInfo := new(unicommodel.SimInfo)
		v.Struct(simInfo)
		if simInfo.RatePlan != plan01 {
			num++
		}
	}
	return num
}

//AutoSetupPlan 自动设置sim卡套餐
func toSetupPlan(planNum int) (map[string]analyse.PlanSimCardInfo, int) {

	//不需要调整的
	var simPlan01List = make(map[string]analyse.PlanSimCardInfo)
	//需求调整的
	var simPlan02List = make(map[string]analyse.PlanSimCardInfo)
	//是否继续进行调节
	continuePlan := 0
	a1, _ := datamodel.SimUnicom{}.GetUnicomSimInfoListByPlan(getPlanName(planNum))
	if a1 == nil {
		return nil, 0
	}
	glog.Info("全部卡数：", len(a1))
	for _, v := range a1 {
		simInfo := new(unicommodel.SimInfo)
		planSimInfo := analyse.PlanSimCardInfo{}

		v.Struct(simInfo)
		planSimInfo.Iccid = simInfo.Iccid
		planSimInfo.Flow = simInfo.CtdDataUsage
		planSimInfo.PlanName = simInfo.RatePlan

		simPlan01List[simInfo.Iccid] = planSimInfo

	}

	planInfo := getListCountInfo(simPlan01List, 1)

	//glog.Info("=====调整前=====")
	if planInfo.OutFlow > 0 {
		var newListAllFlow int64 = 0
		var newListAllNum int64 = 0

		//glog.Info("*******调整后********")
		//glog.Info(getPlanName(planNum), getListCountInfo(simPlan01List, gconv.Int64(planNum)))

		switch planNum {
		case 1:
			for k1, v1 := range simPlan01List {
				f1 := v1.Flow /utils.MB1
				if f1 > 1024 && f1 < 3072 {
					v1.PlanName = plan02
					simPlan02List[v1.Iccid] = v1
					newListAllNum++
					newListAllFlow = newListAllFlow + v1.Flow
					delete(simPlan01List, k1)
				}

			}

			p2 := getListCountInfo(simPlan02List, 2)
			if p2.OutFlow > 0 {
				continuePlan = 2
			}
			glog.Info(getPlanName(2), p2)

		case 2:
			for k1, v1 := range simPlan01List {
				f1 := v1.Flow / utils.MB1
				if f1 > 2048 && f1 < 4096 {
					v1.PlanName = plan03
					simPlan02List[v1.Iccid] = v1
					newListAllNum++
					newListAllFlow = newListAllFlow + v1.Flow
					delete(simPlan01List, k1)
				}

			}
			p2 := getListCountInfo(simPlan02List, 2)
			if p2.OutFlow > 0 {
				continuePlan = 3
			}
			glog.Info(getPlanName(2), p2)

		case 3:
			var simnum1 int64 = 1
			outFlowNumTmp1 := planInfo.OutFlow / utils.MB1
			outFlowNum1 := outFlowNumTmp1 / 3
			if outFlowNum1 > 0 {
				simnum1 = outFlowNum1 + 1
			}

			var i0 int64 = 0
			a0, _ := datamodel.SimUnicom{}.GetUnicomSimInfoListByPlan(getPlanName(1))
			for _, v1 := range a0 {
				simInfo2 := new(unicommodel.SimInfo)
				v1.Struct(simInfo2)
				if simInfo2.CtdDataUsage < 0 {
					if i0 < simnum1 {
						planSimInfo := analyse.PlanSimCardInfo{}
						planSimInfo.Iccid = simInfo2.Iccid
						planSimInfo.PlanName = plan03
						planSimInfo.Flow = simInfo2.CtdDataUsage
						simPlan02List[simInfo2.Iccid] = planSimInfo
					}

					i0++
				}

			}
			p2 := getListCountInfo(simPlan02List, 3)
			glog.Info(getPlanName(3), p2)

		}

		p1 := getListCountInfo(simPlan01List, gconv.Int64(planNum))
		if p1.OutFlow > 0 {
			continuePlan = planNum
		}

	}

	return simPlan02List, continuePlan
}

func getPlanName(num int) string {

	switch num {
	case 1:
		return plan01

	case 2:
		return plan02

	case 3:
		return plan03

	}
	return ""
}

func getListCountInfo(simList map[string]analyse.PlanSimCardInfo, planNum int64) analyse.PlanInfo {

	//计算计费周期
	yearNumStr := gtime.Now().Format("Y")
	monthNum := gtime.Now().Format("n")
	lastMonth := gtime.Now().AddDate(0, -1, 0).Format("n") //上个月
	nextMonth := gtime.Now().AddDate(0, +1, 0).Format("n") //上个月
	//计费周期开始日期
	startDayStr := yearNumStr + "-" + monthNum + "-27" //开始日期
	endDayStr := gconv.String(gtime.Now().Format("Y-m-d")) //结束日期

	if gconv.Int(gtime.Now().Format("j")) < 27 {
		startDayStr = yearNumStr + "-" + lastMonth + "-27"
		//到本月26日
		endDayStr = yearNumStr + "-" + gconv.String(monthNum) + "-26"
	} else {
		startDayStr = yearNumStr + "-" + monthNum + "-27"
		//到下个月26日
		endDayStr = yearNumStr + "-" + gconv.String(nextMonth) + "-26"
	}

	dayStr := gtime.Now().AddDate(0, 0, 0).Format("Y-m-d")

	//计算已使用的天数
	a, _ := time.Parse("2006-01-02", dayStr)
	b, _ := time.Parse("2006-01-02", startDayStr)
	useDayNumtmp := utils.TimeSub(a, b)
	//已使用的天数
	useDayNum := gconv.Int64(useDayNumtmp)

	//计算剩余的天数
	a2, _ := time.Parse("2006-01-02", endDayStr)
	b2, _ := time.Parse("2006-01-02", dayStr)
	remainderDayNumTmp := utils.TimeSub(a2, b2)
	//剩余天数
	remainderDayNum := gconv.Int64(remainderDayNumTmp)

	var useFlow int64 = 0
	for _, v := range simList {
		useFlow = useFlow + v.Flow
	}

	////计算总流量
	planFlow := gconv.Int64(len(simList)) * planNum * utils.MB1
	//剩余流量
	var surplusFlow int64 = 0
	//超出流量
	var outFlow int64 = 0
	//计算剩余流量
	surplusFlow = planFlow - useFlow
	if surplusFlow < 0 {
		outFlow = utils.Abs(surplusFlow)
		surplusFlow = 0
	}

	planInfo := analyse.PlanInfo{}
	planInfo.PlanName = getPlanName(gconv.Int(planNum))
	planInfo.Num = len(simList)
	planInfo.AllFlow = planFlow / utils.MB1
	planInfo.UseFlow = useFlow / utils.MB1
	planInfo.AveDayFlow = useFlow / useDayNum / utils.MB1
	planInfo.AveSimUseFlow = useFlow / gconv.Int64(len(simList)) / utils.MB1
	planInfo.SurplusFlow = surplusFlow / utils.MB1
	planInfo.OutFlow = outFlow / utils.MB1
	planInfo.RemainderDayNum = remainderDayNum
	planInfo.ExpectFlow = useFlow / useDayNum * remainderDayNum / utils.MB1
	return planInfo
}

func (ac *AutoChangePlan) AutoSetupPlan() (error, int) {
	num := 0
	for a := 1; a < 4; a++ {
		aList, c := toSetupPlan(a)
		glog.Info("计划改变卡的数量：", len(aList))
		for _, v := range aList {
			glog.Info("计划改变的卡：", v)
			//go startToChangePlan(v.Iccid, v.PlanName)

		}
		num = num + len(aList)
		if c > 0 {
			bList, _ := toSetupPlan(c)
			for _, v1 := range aList {
				glog.Info("计划改变的卡：",  v1)
				//go startToChangePlan(v1.Iccid, v1.PlanName)

			}
			num = num + len(bList)

		} else {
			break
		}

	}
	return nil, num
}

func (ac *AutoChangePlan) changePlanData(simList []analyse.PlanSimCardInfo, plan string, starNumFlow, endNumFlow, num int64) {
	var i int64 = 0
	for _, v := range simList {
		//这儿里执行处理，资费计划转移动作
		f := gconv.Int64(v.Flow) / utils.MB1
		if f > starNumFlow && f < endNumFlow {
			if i < num {
				glog.Info("变更卡：", v.Iccid, f)
				//go startToChangePlan(v.Iccid, plan)
				ac.SimCardLIst = append(ac.SimCardLIst, v)
			}
			i++
		}
	}
}

func toChangePlan(iccid, ratePlan string) *unicommodel.PutResultData {
	//延时处理
	apiurl := g.Config().Get("unicom.api_url")
	APIURL := apiurl.(string)
	getURL := APIURL + "devices/" + iccid
	searchStr := "{\"ratePlan\":\"" + ratePlan + "\"}"
	dataModel := new(unicommodel.PutResultData)
	collect.PutAPIData(getURL, searchStr, dataModel)
	glog.Info(dataModel)
	return dataModel
}

func startToChangePlan(iccid, ratePlan string) {
	time.Sleep(1e9)
	prd := toChangePlan(iccid, ratePlan)
	if prd.ErrorCode == "40000029" {
		startToChangePlan(iccid, ratePlan)

	}

}

//根据流量排序
func Sort(array []analyse.PlanSimCardInfo) []analyse.PlanSimCardInfo {
	for i := 0; i < len(array)-1; i++ {
		for j := 0; j < len(array)-1-i; j++ {
			//根据流量排序
			if array[j].Flow > array[j+1].Flow { // >升序  <降序
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
	return array
}

