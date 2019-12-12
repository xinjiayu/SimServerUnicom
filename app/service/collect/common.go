package collect

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

// API地址
var ApiURL string

func init() {
	api := g.Config().Get("unicom.api_url")
	ApiURL = api.(string)

}

/**
按照以下说明创建授权标头：
将用户名和 API 密钥组合到一个字符串中，用分号分隔各个值。例如，如果用户名是 starterkit，APIKey 是 d703f52b-1200-4318-ae0d-0f6092b2e6ab，则串联的字符串将是：
starterkit:d703f52b-1200-4318-ae0d-0f6092b2e6ab
使用 Base64(即 RFC2045-MIME)对串联的字符串编码：
c3RhcnRlcmtpdDpkNzAzZjUyYi0xMjAwLTQzMTgtYWUwZC0wZjYw
OTJiMmU2YWI=
将授权标头值设置为 Basic，后跟第 2 步中的编码字符串。确保 Basic 与编码字符串之间有一个空格：
Basic c3RhcnRlcmtpdDpkNzAzZjUyYi0xMjAwLTQzMTgtYWUwZC0wZjYwOTJiMmU2YWI=
*/

// 获取sign授权签名验证字符串
func getSign() string {
	apiKey := g.Config().Get("unicom.api_key")
	apiUser := g.Config().Get("unicom.api_user")
	sign := apiUser.(string) + ":" + apiKey.(string)
	sign = base64.StdEncoding.EncodeToString([]byte(sign))
	return sign
}

//getAPIDataBody 通过api获取数据
func getAPIDataBody(APIURL string) ([]byte, error) {
	c := ghttp.NewClient()
	sgin := getSign()
	c.SetHeader("Authorization", "Basic "+sgin)
	c.SetHeader("Accept", "application/json")
	if res, e := c.Get(APIURL); e != nil {
		return nil, e
	} else {
		defer res.Close()
		body := []byte(res.ReadAllString())
		return body, nil
	}
}

//getAPIData 通过api获取数据
func getAPIData(APIURL string, dataModel interface{}) error {

	c := ghttp.NewClient()
	sgin := getSign()
	c.SetHeader("Authorization", "Basic "+sgin)
	c.SetHeader("Accept", "application/json")
	if res, e := c.Get(APIURL); e != nil {
		//panic(e)
		return e
	} else {
		defer res.Close()
		body := []byte(res.ReadAllString())
		Err := json.Unmarshal(body, &dataModel)
		if Err != nil {
			glog.Error(Err)
			return Err
		}

	}
	return nil
}

//PutAPIData 通过API修改数据
func PutAPIData(apiUrl, content string, dataModel interface{}) {

	c := ghttp.NewClient()
	sgin := getSign()
	c.SetHeader("Authorization", "Basic "+sgin)
	c.SetHeader("Accept", "application/json")
	c.SetHeader("Content-Type", "application/json")
	if r, e := c.Put(apiUrl, content); e != nil {
		panic(e)
	} else {
		defer r.Close()
		glog.Info(r.StatusCode)
		//glog.Info(r.Request.GetBody)
		body := []byte(r.ReadAllString())
		glog.Info(apiUrl, body)
		Err := json.Unmarshal(body, &dataModel)
		if Err != nil {
			glog.Error(Err)
		}
	}
}
