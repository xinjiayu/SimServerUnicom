package datamodel

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/xinjiayu/SimServerUnicom/app/model/basemodel"
)

type UnicomNotice struct {
	Id        int     `json:"id"`
	Iccid     string  `json:"iccid"`
	EventId   string  `json:"event_id"`
	EventType string  `json:"event_type"`
	Timestamp string  `json:"timestamp"`
	DataUsage float64 `json:"datausage"`
	Data1     float64 `json:"data1"`
	Data2     float64 `json:"data1"`

	basemodel.BaseModel

}
// List 获取通知列表数据
func (model UnicomNotice) List(where interface{}, args ...interface{}) []UnicomNotice {
	var resData []UnicomNotice
	err := model.dbModel().OrderBy("timestamp desc").Where(where,args).Structs(&resData)
	if err != nil {
		return []UnicomNotice{}
	}
	return resData
}

// Save 存储数据
func  (model UnicomNotice) Save(event_id, event_type, iccid, data1, data2 string, timestamp, data_usage int) error {
	newData := g.Map{"event_id": event_id, "event_type": event_type, "iccid": iccid, "data1": data1, "data2": data2, "timestamp": timestamp, "data_usage": data_usage}
	glog.Info("Data: ", newData)
	r, _ := model.dbModel().Where("event_id=?", event_id).One()
	if r != nil {
		//修改已有记录
		un := UnicomNotice{}
		if err := r.Struct(&un); err == nil {
			model.dbModel().Data(newData).Where("id=?", un.Id).Update()

		} else {
			glog.Error(err)
		}
	} else {
		//新增记录
		_, err := model.dbModel().Data(newData).Save()
		if err != nil {
			glog.Error(err)
		}

	}

	return nil
}



func (model UnicomNotice) PkVal() int {
	return model.Id
}

func (model UnicomNotice) TableName() string {
	table := g.Config().Get("datatable.notice_table")
	return table.(string)
}

func (model UnicomNotice) dbModel(alias ...string) *gdb.Model {
	var tmpAlias string
	if len(alias) > 0 {
		tmpAlias = " " + alias[0]
	}
	tableModel := g.DB().Table(model.TableName() + tmpAlias).Safe()
	return tableModel
}