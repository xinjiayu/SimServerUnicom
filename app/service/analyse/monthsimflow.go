package analyse

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/xinjiayu/SimServerUnicom/app/model/analyse"
	"github.com/xinjiayu/SimServerUnicom/app/model/datamodel"
	"github.com/xinjiayu/SimServerUnicom/library/utils"
)

func MonthSimFlowList() analyse.TwoMonthFlowCount {

	fud := utils.GetFlowUseDate()
	carrier := "unicom" //sim卡运营商

	//如果当前是1月份，后面取上一年的数据
	lastYear := fud.Year
	if fud.Month == "1" {
		lastYear = fud.LastYear
	}


	//接口返回数据的信息结构
	var twoMonthFlowCount analyse.TwoMonthFlowCount

	for i := 1; i < 32; i++ {
		//输出每天的数字
		twoMonthFlowCount.DayName = append(twoMonthFlowCount.DayName, i)
	}
	//上上个月的流量数据
	BeforeLastMonthFlow := datamodel.SimFlow{}.FlowList("year=? and month=? and carrier=?", lastYear, fud.BeforeLastMonth, carrier)
	BeforeLastMonthEndDayFlowData :=make(map[string]int64)
	for _, v := range BeforeLastMonthFlow {
		data := gconv.Map(v)
		//获取上上个月最后一天的数据
		BeforeLastMonthEndDayFlowData[v.Iccid] = gconv.Int64(data["d"+fud.BeforeLastMonthDays])
	}

	//上个月的流量数据 ===============================
	lastMonthFlow := datamodel.SimFlow{}.FlowList("year=? and month=? and carrier=?", lastYear, fud.LastMonth, carrier)
	lastMonthEndDayFlowData :=make(map[string]int64)

	lastMonthFlowData := make([]int64,gconv.Int(fud.LastMonthDays)) //定义上个月每一天的流量数据表
	lastMonthSimcardNum := make([]int64,gconv.Int(fud.LastMonthDays)) //定义上个月有效卡数据
	lastMonthAvg := make([]int64,gconv.Int(fud.LastMonthDays)) //定义上个月的每天每卡平均流量数据表
	for _, v := range lastMonthFlow {
		data := gconv.Map(v)
		lastMonthEndDayFlowData[v.Iccid] =  gconv.Int64(data["d"+fud.LastMonthDays])		//获取上个月最后一天的数据
		beforeFlowData := BeforeLastMonthEndDayFlowData[v.Iccid] //获取上上个月最后一天的数据
		dayData := countDaysFlow(beforeFlowData, data)
		for i:=0;i<=gconv.Int(fud.LastMonthDays)-1 ;i++  {
			lastMonthFlowData[i]=lastMonthFlowData[i]+dayData[i]
			if dayData[i] > 0 {
				lastMonthSimcardNum[i] = lastMonthSimcardNum[i] + 1
			}
		}

	}

	for i:=0;i<=gconv.Int(fud.LastMonthDays)-1 ;i++  {
		lastMonthAvg[i] = 0
		if lastMonthFlowData[i] >0 &&  lastMonthSimcardNum[i] >0 {
			lastMonthAvg[i] = lastMonthFlowData[i]/lastMonthSimcardNum[i]
		}
	}

	twoMonthFlowCount.Beforemonth = lastMonthFlowData
	twoMonthFlowCount.BeforeMonthSimnum = lastMonthSimcardNum
	twoMonthFlowCount.BeforeMonthAvg = lastMonthAvg



	//本月的流量数据 ================================
	monthFlow := datamodel.SimFlow{}.FlowList("year=? and month=? and carrier=?", fud.Year, fud.Month, carrier)
	monthFlowData := make([]int64,gconv.Int(fud.Today)) //定义到今天的流量数据表
	monthSimcardNum := make([]int64,gconv.Int(fud.Today)) //定义到今天的每一天的有效卡数据
	monthAvg := make([]int64,gconv.Int(fud.Today)) //定义到今天的每天每卡平均流量数据表
	for _, v := range monthFlow {
		data := gconv.Map(v)
		beforeFlowData := lastMonthEndDayFlowData[v.Iccid] //获取上个月最后一天的数据
		dayData := countDaysFlow(beforeFlowData, data)
		for i:=0;i<=gconv.Int(fud.Today)-1 ;i++  {
			monthFlowData[i]=monthFlowData[i]+dayData[i]
			if dayData[i] > 0 {
				monthSimcardNum[i] = monthSimcardNum[i] + 1
			}
		}

	}
	for i:=0;i<=gconv.Int(fud.Today)-1 ;i++  {
		monthAvg[i] = 0
		if monthFlowData[i] >0 &&  monthSimcardNum[i] >0 {
			monthAvg[i] = monthFlowData[i]/monthSimcardNum[i]
		}
	}

	twoMonthFlowCount.Nowmonth = monthFlowData
	twoMonthFlowCount.NowmonthSimnum = monthSimcardNum
	twoMonthFlowCount.NowmonthAvg = monthAvg

	return twoMonthFlowCount
}

