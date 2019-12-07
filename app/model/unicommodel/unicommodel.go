//从联通服务接口获取的数据结构
package unicommodel

type CardList struct {
	PageNumber int `json:"pageNumber"`
	Devices    []struct {
		Iccid             string `json:"iccid"`
		Status            string `json:"status"`
		RatePlan          string `json:"ratePlan"`
		CommunicationPlan string `json:"communicationPlan"`
	} `json:"devices"`
	LastPage bool `json:"lastPage"`
}

type Device struct {
	Iccid             string `json:"iccid"`
	Status            string `json:"status"`
	RatePlan          string `json:"ratePlan"`
	CommunicationPlan string `json:"communicationPlan"`
}

type CardListUsages []struct {
	CtdUsages
}

type CtdUsages struct {
	Iccid                string `json:"iccid"`
	Imsi                 string `json:"imsi"`
	Msisdn               string `json:"msisdn"`
	Imei                 string `json:"imei"`
	Status               string `json:"status"`
	RatePlan             string `json:"ratePlan"`
	CommunicationPlan    string `json:"communicationPlan"`
	CtdDataUsage         int64  `json:"ctdDataUsage"`
	CtdSMSUsage          int    `json:"ctdSMSUsage"`
	CtdVoiceUsage        int    `json:"ctdVoiceUsage"`
	CtdSessionCount      int    `json:"ctdSessionCount"`
	OverageLimitReached  bool   `json:"overageLimitReached"`
	OverageLimitOverride string `json:"overageLimitOverride"`
}

type SimInfo struct {
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
}

type ResultsData struct {
	Past24HDataUsage struct {
		Xmlns     string `json:"-xmlns"`
		DataUsage string `json:"dataUsage"`
		Iccid     string `json:"iccid"`
	} `json:"Past24HDataUsage"`
	CtdUsage struct {
		Xmlns     string `json:"-xmlns"`
		DataUsage string `json:"dataUsage"`
		Iccid     string `json:"iccid"`
	} `json:"CtdUsage"`
}

type ObjValue struct {
	Xmlns     string `json:"-xmlns"`
	DataUsage string `json:"dataUsage"`
	Iccid     string `json:"iccid"`
}


type PutResultData struct {
	Iccid        string `json:"iccid"`
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}
