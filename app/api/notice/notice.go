package notice

import (
	"encoding/json"
	"github.com/xinjiayu/SimServerUnicom/app/model/datamodel"
	"github.com/xinjiayu/SimServerUnicom/app/model/unicommodel"
	"strconv"

	"github.com/gogf/gf/encoding/gxml"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
)


// Controller API管理对象
type Controller struct{}

// DataReceive 通知数据接收
func (c *Controller) DataReceive(r *ghttp.Request) {
	eventId := r.GetString("eventId")
	eventType := r.GetString("eventType")
	timestamp := r.GetString("timestamp")
	//signature := r.GetString("signature")
	data := r.GetString("data")

	glog.Info("接收到的通知数据：",data)

	dataxml, err := gxml.ToJson([]byte(data))
	if err != nil {
		glog.Error(err.Error())
	}

	//包含过去 24 小时的用量超过指定数量的设备的 ICCID 和数据用量。
	//"PAST24H_DATA_USAGE_EXCEEDED":
	resultsData := new(unicommodel.ResultsData)

	Err := json.Unmarshal(dataxml, &resultsData)
	if Err != nil {
		glog.Error(Err)
	}
	iccid := ""
	dataUsage := 0
	switch eventType {
	case "PAST24H_DATA_USAGE_EXCEEDED":

		dataUsageStr := resultsData.Past24HDataUsage.DataUsage
		dataUsage, _ = strconv.Atoi(dataUsageStr)
		iccid = resultsData.Past24HDataUsage.Iccid

	case "CTD_USAGE":
		dataUsageStr := resultsData.CtdUsage.DataUsage
		dataUsage, _ = strconv.Atoi(dataUsageStr)
		iccid = resultsData.CtdUsage.Iccid
	}

	ttime := gtime.NewFromStr(timestamp).Second()
	strInt64 := strconv.FormatInt(ttime, 10)
	tt, _ := strconv.Atoi(strInt64)

	data1 := ""
	data2 := ""
	datamodel.UnicomNotice{}.Save(eventId, eventType, iccid, data1, data2, tt, dataUsage)
}
