package utils

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"time"
)

//1G计算单位
const G1 int64 = 1024 * 1024 * 1024

//1MB计算单位
const MB1 int64 = 1024 * 1024

//统一存储流量统计的时候所用的日期相关信息
type FlowUseDate struct {
	Year            string
	LastYear        string
	Month           string
	LastMonth       string
	BeforeLastMonth string
	Today           string
	Yesterday       string
	BeforeYesterday string

	LastMonthDays string
	BeforeLastMonthDays string
}

//计算两个时间相差的天数
func TimeSub(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(t1.Sub(t2).Hours() / 24)
}

//转为正数
func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

//转换时间到时间戳格式
func ChangeUnixTime(strTime string) string {
	t := gtime.NewFromStr(strTime).Second()
	return gconv.String(t)
}
//计算流量统计多处使用的相关日期信息
func GetFlowUseDate() FlowUseDate {
	var fud = FlowUseDate{}
	fud.Year = gtime.Now().Format("Y")                              //当前年份
	fud.LastYear = gtime.Now().AddDate(-1, 0, 0).Format("Y")        //上一年年份
	fud.Month = gtime.Now().Format("n")                             //当前月份
	fud.LastMonth = gtime.Now().AddDate(0, -1, 0).Format("n")       //上个月
	fud.BeforeLastMonth = gtime.Now().AddDate(0, -2, 0).Format("n") //上上个月

	fud.Today = gtime.Now().Format("j")                             //今天
	fud.Yesterday = gtime.Now().AddDate(0, 0, -1).Format("j")       //昨天
	fud.BeforeYesterday = gtime.Now().AddDate(0, 0, -2).Format("j") //前天

	fud.LastMonthDays = gtime.Now().AddDate(0, -1, 0).Format("t")   //上个月的天数
	fud.LastMonthDays = gtime.Now().AddDate(0, -2, 0).Format("t") //上上个月的天数
	fud.BeforeLastMonthDays = gtime.Now().AddDate(0, -2, 0).Format("t") //上上个月的天数
	return fud
}
