package datamodel

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/xinjiayu/SimServerUnicom/app/model/basemodel"
)

type SimFlow struct {
	Id         int     `json:"id"`
	Iccid      string  `json:"iccid"`
	Month      string  `json:"month"`
	Provider   string  `json:"provider"`
	Carrier    string  `json:"carrier"`
	D1         int64 `json:"d1"`
	D2         int64 `json:"d2"`
	D3         int64 `json:"d3"`
	D4         int64 `json:"d4"`
	D5         int64 `json:"d5"`
	D6         int64 `json:"d6"`
	D7         int64 `json:"d7"`
	D8         int64 `json:"d8"`
	D9         int64 `json:"d9"`
	D10        int64 `json:"d10"`
	D11        int64 `json:"d11"`
	D12        int64 `json:"d12"`
	D13        int64 `json:"d13"`
	D14        int64 `json:"d14"`
	D15        int64 `json:"d15"`
	D16        int64 `json:"d16"`
	D17        int64 `json:"d17"`
	D18        int64 `json:"d18"`
	D19        int64 `json:"d19"`
	D20        int64 `json:"d20"`
	D21        int64 `json:"d21"`
	D22        int64 `json:"d22"`
	D23        int64 `json:"d23"`
	D24        int64 `json:"d24"`
	D25        int64 `json:"d25"`
	D26        int64 `json:"d26"`
	D27        int64 `json:"d27"`
	D28        int64 `json:"d28"`
	D29        int64 `json:"d29"`
	D30        int64 `json:"d30"`
	D31        int64 `json:"d31"`

	basemodel.BaseModel
}


func (model SimFlow) Get() SimFlow {
	if model.Id <= 0 {
		glog.Error(model.TableName() + " get id error")
		return SimFlow{}
	}

	var resData SimFlow
	err := model.dbModel().Where(" id = ?", model.Id).Struct(&resData)
	if err != nil {
		glog.Error(model.TableName()+" get one error", err)
		return SimFlow{}
	}

	return resData
}
func (model SimFlow) GetByOne(where interface{}, args ...interface{}) SimFlow {
	var resData SimFlow
	r, err := model.dbModel().Where(where,args).One()
	if err != nil {
		return SimFlow{}
	}
	r.Struct(&resData)
	return resData
}

func (model SimFlow) FlowList(where interface{}, args ...interface{}) []SimFlow {
	var resData []SimFlow
	err := model.dbModel().Where(where,args).Structs(&resData)
	if err != nil {
		//glog.Error(model.TableName()+" list error", err)
		return []SimFlow{}
	}
	return resData
}

func (model SimFlow) GetSumFlowByOne(where interface{}, args ...interface{}) SimFlow {
	var resData SimFlow
	r, err := model.dbModel().Fields(model.sumColumns()).Where(where,args).One()
	if err != nil {
		return SimFlow{}
	}
	r.Struct(&resData)
	return resData
}

//Save 存储sim卡流量数据到月份表
func (model SimFlow) Save(iccid, flow, provider, carrler, rate_plan, status string, ctd_session_count int) error {

	if iccid == "" {
		return gerror.New("iccid为空值！")
	}
	//获取当前的月份、第几天
	mtime := gtime.Now()
	year := mtime.Format("Y")
	day := mtime.Format("j")
	month := mtime.Format("n")
	//记录或是修改的时间
	ctime := gtime.Now().Format("U")

	r, _ := model.dbModel().Where("iccid=? and year=? and month=?", iccid, year, month).One()
	if r != nil {
		//修改已有记录
		sc := SimFlow{}
		if err := r.Struct(&sc); err == nil {
			data := g.Map{"d" + day: flow, "updatetime": ctime, "rate_plan": rate_plan, "status": status, "ctd_session_count": ctd_session_count}
			//glog.Info("data:", data)
			model.dbModel().Data(data).Where("id=?", sc.Id).Update()

		} else {
			glog.Error(err)
		}
	} else {
		//新增记录
		newData := g.Map{"iccid": iccid, "year": year, "month": month, "d" + day: flow, "provider": provider, "carrier": carrler, "createtime": ctime, "rate_plan": rate_plan, "status": status, "ctd_session_count": ctd_session_count}
		_, err := model.dbModel().Data(newData).Save()
		if err != nil {
			glog.Error(err)
		}

	}

	return nil
}


func (model SimFlow) PkVal() int {
	return model.Id
}

func (model SimFlow) TableName() string {
	table := g.Config().Get("datatable.flow_table")
	return table.(string)
}

func (model SimFlow) dbModel(alias ...string) *gdb.Model {
	var tmpAlias string
	if len(alias) > 0 {
		tmpAlias = " " + alias[0]
	}
	tableModel := g.DB().Table(model.TableName() + tmpAlias).Safe()
	return tableModel
}

func (model SimFlow) sumColumns() string {
	sqlColumns := "year,month,carrier,sum(d1) as d1,sum(d2) as d2,sum(d3) as d3,sum(d4) as d4,sum(d5) as d5,sum(d6) as d6,sum(d7) as d7,sum(d8) as d8,sum(d9) as d9,sum(d10) as d10,sum(d11) as d11,sum(d12) as d12,sum(d13) as d13,sum(d14) as d14,sum(d15) as d15,sum(d16) as d16,sum(d17) as d17,sum(d18) as d18,sum(d19) as d19,sum(d20) as d20,sum(d21) as d21,sum(d22) as d22,sum(d23) as d23,sum(d24) as d24,sum(d25) as d25,sum(d26) as d26,sum(d27) as d27,sum(d28) as d28,sum(d29) as d29,sum(d30) as d30,sum(d31) as d31"
	return sqlColumns
}