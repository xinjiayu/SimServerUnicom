package analyse

//最近两天的数据结构
type TwoDaysFlow struct {
	Iccid          string `json:"iccid"`
	TodayUsage     int64  `json:"today_usage"`
	YesterdayUsage int64  `json:"yesterday_usage"`
}

//计划内卡的信息结构
type PlanSimCardInfo struct {
	Iccid    string
	Flow     int64
	PlanName string
}

//计费套餐计划信息结构
type PlanInfo struct {
	PlanName      string
	AllFlow       int64 //流量池总流量
	UseFlow       int64 //已用流量
	AveSimUseFlow int64 //每个卡已用平均流量
	AveDayFlow    int64 //日均流量
	SurplusFlow   int64 //剩余流量
	OutFlow       int64 //超出流量
	ExpectFlow    int64 //预计需要的流量
	SurplusDayNum int64 //剩余天数
	Num           int
}

//指定卡两个月流量统计信息结构
type TwoDaysBaySimCardFlow struct {
	DayName     []int   `json:"day_name"`
	Nowmonth    []int64 `json:"nowmonth"`
	Beforemonth []int64 `json:"beforemonth"`
}

//两个月所有流量的统计信息结构
type TwoMonthFlowCount struct {
	DayName           []int   `json:"day_name"`
	Nowmonth          []int64 `json:"nowmonth"`
	NowmonthSimnum    []int64 `json:"nowmonth_simnum"`
	NowmonthAvg       []int64 `json:"nowmonth_avg"`
	Beforemonth       []int64 `json:"beforemonth"`
	BeforeMonthSimnum []int64 `json:"before_month_simnum"`
	BeforeMonthAvg    []int64 `json:"before_month_avg"`
}
