package collect

import (
	"encoding/json"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/util/gconv"
	"github.com/qianlnk/pgbar"
	"github.com/xinjiayu/SimServerUnicom/app/model/datamodel"
	"github.com/xinjiayu/SimServerUnicom/app/model/unicommodel"
	"strconv"
	"sync"
)

//获取指定sim卡的流量，并存入运营数据库
func GetUsages(iccid string) interface{} {
	getURL := ApiURL + "devices/" + iccid + "/ctdUsages"
	dataModel := new(unicommodel.CtdUsages)
	err := getAPIData(getURL, dataModel)
	if err != nil {
		return nil
	}

	//存入数据库
	flowNum := dataModel.CtdDataUsage
	flow := gconv.String(flowNum)
	go datamodel.SimFlow{}.Save(dataModel.Iccid, flow, "cqunicom", "unicom", dataModel.RatePlan, dataModel.Status, dataModel.CtdSessionCount)
	return dataModel

}

// CardList 通过联通平台的devices接口获取设备卡列表
func GetCardList() unicommodel.CardList {
	pageSize := "50"
	pageNumber := "1"
	status := "ACTIVATED"
	var page int = 0
	var simList unicommodel.CardList
	for {
		page++
		//glog.Info(page, "=========页数=================")
		pageNumber = strconv.Itoa(page)
		getURL := ApiURL + "devices?modifiedSince=2016-04-18T17%3A31%3A34%2B00%3A00&pageSize=" + pageSize + "&pageNumber=" + pageNumber + "&status=" + status
		dataModel := new(unicommodel.CardList)
		getAPIData(getURL, dataModel)
		for _, v := range dataModel.Devices {
			simList.Devices = append(simList.Devices, v)
		}
		if dataModel.LastPage {
			break
		}
	}
	go getSimData(simList) //直接起个线程跳走，加快返回籹据的速度
	return simList
}

func getSimData(dd unicommodel.CardList) {
	pool := grpool.New(3)
	wg := sync.WaitGroup{}

	var vlist unicommodel.CardList
	b4 := pgbar.NewBar(0, "获取卡流量与信息", len(dd.Devices))
	for _, v := range dd.Devices {
		wg.Add(1)
		vd := v //采用临时变量的形式来传递当前变量v的值
		pool.Add(func() {
			getURL := ApiURL + "devices/" + vd.Iccid + "/ctdUsages"
			dataModel := new(unicommodel.CtdUsages)
			body, err := getAPIDataBody(getURL)
			if err == nil {
				Err := json.Unmarshal(body, &dataModel)
				if Err != nil {
					glog.Error(Err)
				}

				flowNum := dataModel.CtdDataUsage
				flow := gconv.String(flowNum)
				if dataModel.Iccid != "" {
					//存入数据库
					datamodel.SimFlow{}.Save(dataModel.Iccid, flow, "cqunicom", "unicom", dataModel.RatePlan, dataModel.Status, dataModel.CtdSessionCount)
					//glog.Info("sim卡详细信息：",dataModel)
					getUnicomSimInfo(vd.Iccid, flowNum)
				} else {
					vlist.Devices = append(vlist.Devices, vd)
				}
			}
			b4.Add()
			wg.Done()
		})
	}
	wg.Wait()

	//如果有获取失败的sim卡，继续获取
	if len(vlist.Devices) > 0 {
		getSimData(vlist)
		//glog.Info(vlist)
		vlistNum := len(vlist.Devices)
		glog.Infof("====未成功获取籹据的卡数：%d \n", vlistNum)

	}
	glog.Infof("sim卡数据获取结束！")

}

//获取sim卡的详细信息
func getUnicomSimInfo(iccid string, flow int64) {
	getURL := ApiURL + "devices/" + iccid
	dataModel := new(unicommodel.SimInfo)
	err := getAPIData(getURL, dataModel)
	if err != nil {
		getAPIData(getURL, dataModel)
	}
	dataModel.CtdDataUsage = flow
	//存入数据库
	datamodel.SimUnicom{}.SaveUnicomSimInfo(*dataModel)

}