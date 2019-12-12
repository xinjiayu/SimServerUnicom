package operate

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/database/gdb"
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
	PlanInfo        map[string]analyse.PlanInfo
	PlanListSimList *gmap.AnyAnyMap
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
func (ac *AutoChangePlan) toSetupPlan(planNum int) (map[string]datamodel.SimUnicom, int) {
	//需求调整的
	var simPlan02List = make(map[string]datamodel.SimUnicom)
	//是否继续进行调节
	continuePlan := 0
	aListData := ac.PlanListSimList.Get(getPlanName(planNum))
	if aListData == nil {
		return nil, 0
	}
	aList := gconv.Map(aListData)

	planInfo := ac.PlanInfo[getPlanName(planNum)]

	if planInfo.OutFlow > 0 {
		var newListAllFlow int64 = 0
		var newListAllNum int64 = 0

		switch planNum {
		case 1:
			for k1, v1 := range aList {
				simInfo := datamodel.SimUnicom{}
				gconv.Struct(v1, &simInfo)

				f1 := simInfo.CtdDataUsage / utils.MB1
				if f1 > 3072 {
					simInfo.RatePlan = plan02
					simPlan02List[simInfo.Iccid] = simInfo
					newListAllNum++
					newListAllFlow = newListAllFlow + simInfo.CtdDataUsage
					delete(aList, k1)
				}

			}
			p := ac.PlanInfo[getPlanName(planNum)]
			if p.OutFlow > 0 {
				continuePlan = 2
			}
			glog.Info("计划一：", getPlanName(1), p)

		case 2:
			for k1, v1 := range aList {
				simInfo := datamodel.SimUnicom{}
				gconv.Struct(v1, &simInfo)
				f1 := simInfo.CtdDataUsage / utils.MB1
				if f1 > 4096 {
					simInfo.RatePlan = plan03
					simPlan02List[simInfo.Iccid] = simInfo
					newListAllNum++
					newListAllFlow = newListAllFlow + simInfo.CtdDataUsage
					delete(aList, k1)
				}

			}
			p := ac.PlanInfo[getPlanName(planNum)]
			if p.OutFlow > 0 {
				continuePlan = 3
			}
			glog.Info("计划二：", getPlanName(2), p)

		case 3:
			var simnum1 int64 = 1
			outFlowNumTmp1 := planInfo.OutFlow / utils.MB1
			outFlowNum1 := outFlowNumTmp1 / 3
			if outFlowNum1 > 0 {
				simnum1 = outFlowNum1 + 1
			}

			var cnt int64 = 0
			//a0, _ := ac.PlanListSimList[plan01]
			a1 := ac.PlanListSimList.Get(plan01)

			if a1 == nil {
				return nil, 0
			}
			aList := gconv.Map(a1)
			for k1, v1 := range aList {
				simInfo := datamodel.SimUnicom{}
				gconv.Struct(v1, &simInfo)
				if simInfo.CtdDataUsage/utils.MB1 < 500 {
					if cnt < simnum1 {
						simInfo.RatePlan = plan03
						simPlan02List[simInfo.Iccid] = simInfo
						delete(aList, k1)

					}

					cnt++
				}

			}
			p := ac.PlanInfo[getPlanName(planNum)]
			if p.OutFlow > 0 {
				continuePlan = 3
			}
			glog.Info("计划三有超出：", getPlanName(3), p)

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

func (ac *AutoChangePlan) AutoSetupPlan() (error, int) {
	simCardList, _ := datamodel.SimUnicom{}.GetUnicomSimInfoList()
	ac.CountPlanFlow(simCardList)
	glog.Info("全部卡数：", len(simCardList))

	var changeSimcardList = make(map[string]string)

	for a := 1; a < 4; a++ {
		aList, c := ac.toSetupPlan(a)
		for _, v := range aList {
			changeSimcardList[v.Iccid] = v.RatePlan
		}
		if c > 0 {
			bList, _ := ac.toSetupPlan(c)
			for _, v1 := range bList {
				changeSimcardList[v1.Iccid] = v1.RatePlan
			}

		}

	}
	for cid, plan := range changeSimcardList {
		glog.Info(cid, plan)
		go startToChangePlan(cid, plan)
	}

	glog.Info("sim卡计划变更数：", len(changeSimcardList))

	return nil, len(changeSimcardList)
}

func toChangePlan(iccid, ratePlan string) *unicommodel.PutResultData {
	//延时处理
	apiurl := g.Config().Get("unicom.api_url")
	APIURL := apiurl.(string)
	getURL := APIURL + "devices/" + iccid
	searchStr := "{\"ratePlan\":\"" + ratePlan + "\"}"
	glog.Info("提交", searchStr)
	dataModel := new(unicommodel.PutResultData)
	collect.PutAPIData(getURL, searchStr, dataModel)
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

//计算流量池的真实大小：流量池中，当前计算周期新激活卡，是从卡激活的日期到计算结束日按每天流量累计的。
func (ac *AutoChangePlan) CountPlanFlow(simList gdb.Result) {
	ac.PlanInfo = make(map[string]analyse.PlanInfo)
	//ac.PlanListSimList = make(map[string][]datamodel.SimUnicom)

	var simCordList1 = make(map[string]datamodel.SimUnicom)
	var simCordList2 = make(map[string]datamodel.SimUnicom)
	var simCordList3 = make(map[string]datamodel.SimUnicom)

	ac.PlanListSimList = gmap.NewAnyAnyMap()

	var flowPoolSize = make(map[string]int64) //池子总量
	var useFlow = make(map[string]int64)      //池子使用量
	var simNum = make(map[string]int)         //sim卡的数量

	monthNum := gconv.Int64(gtime.Now().Format("n"))
	yearNum := gconv.Int64(gtime.Now().Format("Y"))
	toDayNum := gconv.Int64(gtime.Now().Format("j"))

	startYear := yearNum       //计费周期开始的年
	startMonth := monthNum - 1 //计划周期开始的月

	if monthNum == 1 {
		if toDayNum <= 26 {
			startYear = startYear - 1
			startMonth = 12
		}
	}

	endYear := yearNum   //计费周期结束的年
	endMonth := monthNum //计划周期结束的月
	if toDayNum > 26 {
		endMonth = endMonth + 1
		if monthNum == 12 {
			endYear = endYear + 1
			endMonth = 1
		}
	}

	//计费周期开始日期
	startDayStr := gconv.String(startYear) + "-" + gconv.String(startMonth) + "-27" //计周期开始日期
	nowDayStr := gconv.String(gtime.Now().Format("Y-m-d"))                          //今天的日期
	endDayStr := gconv.String(endYear) + "-" + gconv.String(endMonth) + "-26"       //计费周期结束的日期

	//计算计费周期还剩余的天数
	t1, _ := time.Parse("2006-01-02", endDayStr)
	t2, _ := time.Parse("2006-01-02", nowDayStr)
	surplusFlowDayNumTmp := utils.TimeSub(t1, t2)
	surplusDayNum := gconv.Int64(surplusFlowDayNumTmp)

	//计算已使用的天数
	a, _ := time.Parse("2006-01-02", nowDayStr)
	b, _ := time.Parse("2006-01-02", startDayStr)
	useDayNumtmp := utils.TimeSub(a, b)
	//已使用的天数
	useDayNum := gconv.Int64(useDayNumtmp)

	var G1DayFlow int64 = utils.G1 / 30
	var G2DayFlow int64 = utils.G1 * 2 / 30
	var G3DayFlow int64 = utils.G1 * 3 / 30

	for _, v := range simList {
		simInfo := datamodel.SimUnicom{}
		v.Struct(&simInfo)
		//计算当前周期内激活卡的池子流量
		var chargingDayNum int64 = 30
		tmpTime := gconv.Int64(simInfo.DateActivated)
		simCardActivatedYear := gconv.Int64(gtime.NewFromTimeStamp(tmpTime).Format("Y"))
		simCardActivatedMonth := gconv.Int64(gtime.NewFromTimeStamp(tmpTime).Format("n"))
		if simCardActivatedYear == yearNum && simCardActivatedMonth == monthNum {
			simCardActivatedDate := gtime.NewFromTimeStamp(tmpTime).Format("Y-m-d")
			//计算剩余的天数
			a2, _ := time.Parse("2006-01-02", endDayStr)
			b2, _ := time.Parse("2006-01-02", simCardActivatedDate)
			chargingDayNumTmp := utils.TimeSub(a2, b2)
			chargingDayNum = gconv.Int64(chargingDayNumTmp)
		}

		switch simInfo.RatePlan {
		case plan01:
			flowPoolSize[plan01] = flowPoolSize[plan01] + chargingDayNum*G1DayFlow
			useFlow[plan01] = useFlow[plan01] + simInfo.CtdDataUsage
			simCordList1[simInfo.Iccid] = simInfo
			ac.PlanListSimList.Set(plan01, simCordList1)
			simNum[plan01]++

		case plan02:
			flowPoolSize[plan02] = flowPoolSize[plan02] + chargingDayNum*G2DayFlow
			useFlow[plan02] = useFlow[plan02] + simInfo.CtdDataUsage
			simCordList2[simInfo.Iccid] = simInfo
			ac.PlanListSimList.Set(plan02, simCordList2)
			simNum[plan02]++

		case plan03:
			flowPoolSize[plan03] = flowPoolSize[plan03] + chargingDayNum*G3DayFlow
			useFlow[plan03] = useFlow[plan03] + simInfo.CtdDataUsage
			simCordList3[simInfo.Iccid] = simInfo
			ac.PlanListSimList.Set(plan03, simCordList3)
			simNum[plan03]++

		}
	}

	//1G套餐统计 ===============================
	planInfo1 := analyse.PlanInfo{}
	planInfo1.PlanName = plan01
	planInfo1.AllFlow = flowPoolSize[plan01] / utils.MB1 //池子总量
	planInfo1.UseFlow = useFlow[plan01] / utils.MB1      //使用量
	//计算剩余流量
	surplusFlow := flowPoolSize[plan01] - useFlow[plan01]
	var outFlow int64 = 0
	if surplusFlow < 0 {
		outFlow = utils.Abs(surplusFlow)
		surplusFlow = 0
	}
	planInfo1.SurplusFlow = surplusFlow / utils.MB1                  //余量
	planInfo1.OutFlow = outFlow / utils.MB1                          //超出量
	planInfo1.AveDayFlow = planInfo1.UseFlow / useDayNum / utils.MB1 //每天的使用流量
	planInfo1.Num = simNum[plan01]                                   //sim卡数量

	planInfo1.AveSimUseFlow = 0
	if simNum[plan01] > 0 {
		planInfo1.AveSimUseFlow = planInfo1.AveDayFlow / gconv.Int64(simNum[plan01]) / utils.MB1 //计算每天每卡的平均量
	}

	planInfo1.SurplusDayNum = surplusDayNum //剩余的天数
	planInfo1.ExpectFlow = planInfo1.AveDayFlow * surplusDayNum / utils.MB1
	ac.PlanInfo[plan01] = planInfo1

	//2G套餐统计 ===============================
	planInfo2 := analyse.PlanInfo{}
	planInfo2.PlanName = plan02
	planInfo2.AllFlow = flowPoolSize[plan02] / utils.MB1 //池子总量
	planInfo2.UseFlow = useFlow[plan02] / utils.MB1      //使用量
	//计算剩余流量
	surplusFlow2 := flowPoolSize[plan02] - useFlow[plan02]
	var outFlow2 int64 = 0
	if surplusFlow2 < 0 {
		outFlow2 = utils.Abs(surplusFlow2)
		surplusFlow2 = 0
	}
	planInfo2.SurplusFlow = surplusFlow2 / utils.MB1                 //余量
	planInfo2.OutFlow = outFlow2 / utils.MB1                         //超出量
	planInfo2.AveDayFlow = planInfo2.UseFlow / useDayNum / utils.MB1 //每天的使用流量
	planInfo2.Num = simNum[plan02]                                   //sim卡数量

	planInfo2.AveSimUseFlow = 0
	if simNum[plan02] > 0 {
		planInfo2.AveSimUseFlow = planInfo2.AveDayFlow / gconv.Int64(simNum[plan02]) / utils.MB1 //计算每天每卡的平均量
	}

	planInfo2.SurplusDayNum = surplusDayNum //剩余的天数
	planInfo2.ExpectFlow = planInfo2.AveDayFlow * surplusDayNum / utils.MB1

	ac.PlanInfo[plan02] = planInfo2

	//3G套餐统计 ===============================
	planInfo3 := analyse.PlanInfo{}
	planInfo3.PlanName = plan03
	planInfo3.AllFlow = flowPoolSize[plan03] / utils.MB1 //池子总量
	planInfo3.UseFlow = useFlow[plan03] / utils.MB1      //使用量
	//计算剩余流量
	surplusFlow3 := flowPoolSize[plan03] - useFlow[plan03]
	var outFlow3 int64 = 0
	if surplusFlow3 < 0 {
		outFlow3 = utils.Abs(surplusFlow3)
		surplusFlow3 = 0
	}
	planInfo3.SurplusFlow = surplusFlow3 / utils.MB1                 //余量
	planInfo3.OutFlow = outFlow3 / utils.MB1                         //超出量
	planInfo3.AveDayFlow = planInfo3.UseFlow / useDayNum / utils.MB1 //每天的使用流量
	planInfo3.Num = simNum[plan03]                                   //sim卡数量

	planInfo3.AveSimUseFlow = 0
	if simNum[plan03] > 0 {
		planInfo3.AveSimUseFlow = planInfo3.AveDayFlow / gconv.Int64(simNum[plan03]) / utils.MB1 //计算每天每卡的平均量
	}

	planInfo3.SurplusDayNum = surplusDayNum //剩余的天数
	planInfo3.ExpectFlow = planInfo3.AveDayFlow * surplusDayNum / utils.MB1

	ac.PlanInfo[plan03] = planInfo3
}
