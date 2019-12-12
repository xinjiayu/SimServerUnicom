package datamodel

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/xinjiayu/SimServerUnicom/app/model/basemodel"
	"github.com/xinjiayu/SimServerUnicom/app/model/unicommodel"
	"github.com/xinjiayu/SimServerUnicom/library/utils"
)

type SimUnicom struct {
	Id                int    `json:"id"`
	Iccid             string `json:"iccid"`
	Imsi              string `json:"imsi"`
	Msisdn            string `json:"msisdn"`
	Imei              string `json:"imei"`
	Status            string `json:"status"`
	RatePlan          string `json:"rateplan"`
	CommunicationPlan string `json:"communicationplan"`
	Customer          string `json:"customer"`
	EndConsumerID     string `json:"endconsumerid"`
	DateActivated     string `json:"dateactivated"`
	DateAdded         string `json:"dateadded"`
	DateUpdated       string `json:"dateupdated"`
	DateShipped       string `json:"dateshipped"`
	AccountID         string `json:"accountid"`
	FixedIPAddress    string `json:"fixedipaddress"`
	AccountCustom1    string `json:"accountcustom1"`
	AccountCustom2    string `json:"accountcustom2"`
	AccountCustom3    string `json:"accountcustom3"`
	AccountCustom4    string `json:"accountcustom4"`
	AccountCustom5    string `json:"accountcustom5"`
	AccountCustom6    string `json:"accountcustom6"`
	AccountCustom7    string `json:"accountcustom7"`
	AccountCustom8    string `json:"accountcustom8"`
	AccountCustom9    string `json:"accountcustom9"`
	AccountCustom10   string `json:"accountcustom10"`
	SimNotes          string `json:"simnotes"`
	DeviceID          string `json:"deviceid"`
	ModemID           string `json:"modemid"`
	GlobalSimType     string `json:"globalsimtype"`
	CtdDataUsage      int64  `json:"ctddatausage"`

	basemodel.BaseModel
}

const Plan01 = "831WLW016555_MON-FLEX_1024M_SP"
const Plan02 = "831WLW016555_MON-FLEX_2048M_SP"
const Plan03 = "831WLW016555_MON-FLEX_3072M_SP"

func (model SimUnicom) Get() SimUnicom {
	if model.Id <= 0 {
		glog.Error(model.TableName() + " get id error")
		return SimUnicom{}
	}

	var resData SimUnicom
	err := model.dbModel().Where(" id = ?", model.Id).Struct(&resData)
	if err != nil {
		glog.Error(model.TableName()+" get one error", err)
		return SimUnicom{}
	}

	return resData
}

func (model SimUnicom) GetByIccid(iccid string) SimUnicom {
	if iccid == "" {
		return SimUnicom{}
	}

	var resData SimUnicom
	err := model.dbModel().Where(" iccid = ?", iccid).Struct(&resData)
	if err != nil {
		glog.Error(model.TableName()+" get one error", err)
		return SimUnicom{}
	}

	return resData
}

//计费周期内流量使用列表
func (model SimUnicom) FlowList(rateplan string) []SimUnicom {
	var resData []SimUnicom
	var err error
	if rateplan != "" {
		err = model.dbModel().Where(" rateplan = ?", rateplan).OrderBy("ctddatausage desc").Structs(&resData)
	} else {
		err = model.dbModel().Where(" 1 = ?", 1).OrderBy("ctddatausage desc").Structs(&resData)
	}
	if err != nil {
		glog.Error(model.TableName()+" list error", err)
		return []SimUnicom{}
	}
	return resData
}

//func (model SimUnicom) PlanCountInfo(planName string) analyse.PlanInfo {
//
//	var simList = make(map[string]analyse.PlanSimCardInfo)
//	var resData []SimUnicom
//	var err error
//	if planName != "" {
//		err = model.dbModel().Where(" rateplan = ?", planName).OrderBy("ctddatausage desc").Structs(&resData)
//	} else {
//		err = model.dbModel().Where(" 1 = ?", 1).OrderBy("ctddatausage desc").Structs(&resData)
//	}
//	if err != nil {
//		//glog.Error(model.TableName()+" list error", err)
//		return analyse.PlanInfo{}
//	}
//
//	for _, v := range resData {
//
//		planSimCardInfo := analyse.PlanSimCardInfo{}
//		planSimCardInfo.Flow = v.CtdDataUsage
//		planSimCardInfo.Iccid = v.Iccid
//		planSimCardInfo.PlanName = v.RatePlan
//		simList[v.Iccid] = planSimCardInfo
//
//	}
//
//	//计算每种计费套餐的流量池的真实大小。
//	planPoolFlow := countPlanFlow(resData)
//	planFlow := planPoolFlow[planName]
//
//	//计算计费周期
//	yearNumStr := gtime.Now().Format("Y")
//	monthNum := gtime.Now().Format("n")
//	lastMonth := gtime.Now().AddDate(0, -1, 0).Format("n") //上个月
//	nextMonth := gtime.Now().AddDate(0, +1, 0).Format("n") //上个月
//
//	//计费周期开始日期
//	startDayStr := yearNumStr + "-" + monthNum + "-27"
//	endDayStr := gconv.String(gtime.Now().Format("Y-m-d"))
//
//	if gconv.Int(gtime.Now().Format("j")) < 27 {
//		startDayStr = yearNumStr + "-" + lastMonth + "-27"
//		//到本月26日
//		endDayStr = yearNumStr + "-" + gconv.String(monthNum) + "-26"
//	} else {
//		startDayStr = yearNumStr + "-" + monthNum + "-27"
//		//到下个月26日
//		endDayStr = yearNumStr + "-" + gconv.String(nextMonth) + "-26"
//	}
//	//glog.Info("开始日期：",startDayStr)
//	//glog.Info("结束日期：",endDayStr)
//
//	dayStr := gtime.Now().AddDate(0, 0, 0).Format("Y-m-d")
//	//glog.Info("今天日期：",dayStr)
//
//	//计算已使用的天数
//	a, _ := time.Parse("2006-01-02", dayStr)
//	b, _ := time.Parse("2006-01-02", startDayStr)
//	useDayNumtmp := utils.TimeSub(a, b)
//	//已使用的天数
//	useDayNum := gconv.Int64(useDayNumtmp)
//	//glog.Info("已用天数：",useDayNum)
//
//	//计算剩余的天数
//	a2, _ := time.Parse("2006-01-02", endDayStr)
//	b2, _ := time.Parse("2006-01-02", dayStr)
//	remainderDayNumTmp := utils.TimeSub(a2, b2)
//	remainderDayNum := gconv.Int64(remainderDayNumTmp)
//
//	//glog.Info("剩余天数：",remainderDayNum)
//	var useFlow int64 = 0
//	for _, v := range simList {
//		useFlow = useFlow + v.Flow
//	}
//
//	//剩余流量
//	var surplusFlow int64 = 0
//	//超出流量
//	var outFlow int64 = 0
//	//计算剩余流量
//	surplusFlow = planFlow - useFlow
//	if surplusFlow < 0 {
//		outFlow = utils.Abs(surplusFlow)
//		surplusFlow = 0
//	}
//
//	planInfo := analyse.PlanInfo{}
//	planInfo.PlanName = planName
//	planInfo.Num = len(simList)
//	planInfo.AllFlow = planFlow / utils.MB1
//	planInfo.UseFlow = useFlow / utils.MB1
//	planInfo.AveDayFlow = useFlow / useDayNum / utils.MB1
//	planInfo.AveSimUseFlow = useFlow / gconv.Int64(len(simList)) / utils.MB1 //这个留MB适合
//	planInfo.SurplusFlow = surplusFlow / utils.MB1
//	planInfo.OutFlow = outFlow / utils.MB1
//	planInfo.RemainderDayNum = remainderDayNum
//	planInfo.ExpectFlow = useFlow / useDayNum * remainderDayNum / utils.MB1
//	return planInfo
//}

//计算流量池的真实大小：流量池中，当前计费周期内新激活的卡，流量是从卡激活的日期到计费结束日按每天流量累计的流量。
//func countPlanFlow(simList []SimUnicom) map[string]int64 {
//
//	var retData = make(map[string]int64, 3)
//	monthNum := gconv.Int64(gtime.Now().Format("n"))
//	yearNum := gconv.Int64(gtime.Now().Format("Y"))
//	toDayNum := gconv.Int64(gtime.Now().Format("j"))
//	endYear := yearNum
//	endMonth := monthNum
//	if toDayNum > 26 {
//		endMonth = endMonth + 1
//		if monthNum == 12 {
//			endYear = endYear + 1
//			endMonth = 1
//		}
//	}
//
//	var G1DayFlow int64 = utils.G1 / 30
//	var G2DayFlow int64 = utils.G1 * 2 / 30
//	var G3DayFlow int64 = utils.G1 * 3 / 30
//
//	for _, v := range simList {
//		//dateActivated
//		var remainderDayNum int64 = 30
//		if v.DateActivated != "" {
//			tmpTime := gconv.Int64(v.DateActivated)
//			simCardActivatedYear := gconv.Int64(gtime.NewFromTimeStamp(tmpTime).Format("Y"))
//			simCardActivatedMonth := gconv.Int64(gtime.NewFromTimeStamp(tmpTime).Format("n"))
//			if simCardActivatedYear == yearNum && simCardActivatedMonth == monthNum {
//				simCardActivatedDate := gtime.NewFromTimeStamp(tmpTime).Format("Y-m-d")
//				endDayStr := gconv.String(endYear) + "-" + gconv.String(endMonth) + "-26"
//				//计算剩余的天数
//				a2, _ := time.Parse("2006-01-02", endDayStr)
//				b2, _ := time.Parse("2006-01-02", simCardActivatedDate)
//				remainderDayNumTmp := utils.TimeSub(a2, b2)
//				remainderDayNum = gconv.Int64(remainderDayNumTmp)
//			}
//
//		}
//
//		switch v.RatePlan {
//		case Plan01:
//			retData[Plan01] = retData[Plan01] + remainderDayNum*G1DayFlow
//
//		case Plan02:
//			retData[Plan02] = retData[Plan02] + +remainderDayNum*G2DayFlow
//
//		case Plan03:
//			retData[Plan03] = retData[Plan03] + +remainderDayNum*G3DayFlow
//
//		}
//	}
//	return retData
//}

//SaveUnicomSimInfo 存储联通sim卡详细信息到数据表
func (model SimUnicom) SaveUnicomSimInfo(sim unicommodel.SimInfo) error {
	if sim.Iccid == "" {
		return gerror.New("iccid为空值！")
	}

	if sim.DateActivated != "" {
		sim.DateActivated = utils.ChangeUnixTime(sim.DateActivated)

	}
	if sim.DateAdded != "" {
		sim.DateAdded = utils.ChangeUnixTime(sim.DateAdded)

	}
	if sim.DateShipped != "" {
		sim.DateShipped = utils.ChangeUnixTime(sim.DateShipped)

	}
	if sim.DateUpdated != "" {
		sim.DateUpdated = utils.ChangeUnixTime(sim.DateUpdated)

	}

	r, _ := model.dbModel().Where("iccid=?", sim.Iccid).One()
	if r != nil {
		//修改已有记录
		sc := unicommodel.SimInfo{}
		if err := r.Struct(&sc); err == nil {
			//glog.Info("data:", sim)
			model.dbModel().Data(sim).Where("iccid=?", sc.Iccid).Update()

		} else {
			glog.Error(err)
		}
	} else {
		//新增记录
		_, err := model.dbModel().Data(sim).Save()
		if err != nil {
			glog.Error(err)
		}
	}

	return nil
}

//获取所有资费计划的卡信息
func (model SimUnicom) GetUnicomSimInfoList() (gdb.Result, error) {
	return model.dbModel().All()

}

//获取指定资费计划的卡信息
func (model SimUnicom) GetUnicomSimInfoListByPlan(plan string) (gdb.Result, error) {
	return model.dbModel().Where("rateplan=?", plan).All()

}

func (model SimUnicom) PkVal() int {
	return model.Id
}

func (model SimUnicom) TableName() string {
	table := g.Config().Get("datatable.info_table")
	return table.(string)
}

func (model SimUnicom) dbModel(alias ...string) *gdb.Model {
	var tmpAlias string
	if len(alias) > 0 {
		tmpAlias = " " + alias[0]
	}
	tableModel := g.DB().Table(model.TableName() + tmpAlias).Safe()
	return tableModel
}
func getPlanNum(panName string) int64 {

	switch panName {
	case Plan01:
		return 1

	case Plan02:
		return 2

	case Plan03:
		return 3

	}
	return 1
}
