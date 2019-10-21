package SplunkRestFul

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/oliveagle/jsonpath"
	"net/http"
	"strconv"
)

type SReault struct {
	SessionKey string `json:"sessionKey"`
	Sid string `json:"sid"`
}

type SplunkQuery struct {
	baseUrl string
	userName string
	password string
	sessionKey string
}

type JobParam struct{
	Earliest 		string `earliest`
	Latest 			string `latest`
	SerialNumber    string `serialnumber`
}


func NewJobParam(serial string) *JobParam {
	return &JobParam{SerialNumber: serial}
}

func NewJobParamEx(serial string,earliest string,lastest string) *JobParam {
	return &JobParam{SerialNumber: serial,Earliest:earliest,Latest:	lastest}
}

func(j *JobParam)ToString()string {
	query := `index="ysb-syslog" `
	if len(j.SerialNumber) > 0{
		query += j.SerialNumber
	}

	if len(j.Earliest) > 0{
		query += " earliest=" + j.Earliest
	}

	if len(j.Latest) > 0{
		query += " latest=" + j.Latest
	}
	query += " |reverse"
	return query
}

func NewSplunkQuery(userName string, password string) *SplunkQuery {
	var Splunk = &SplunkQuery{baseUrl: `https://log.bigfintax.com:18001/services`, userName: userName, password: password}
	Splunk.sessionKey = Splunk.GetKey()
	return  Splunk
}

func(this *SplunkQuery) GetKey() string{
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	req := httplib.Post(this.baseUrl+ `/auth/login`).SetTransport(tr)
	req.Param("username",this.userName)
	req.Param("password",this.password)
	req.Param("output_mode","json")
	var  Result SReault
	_ = req.ToJSON(&Result)
	fmt.Println(Result.SessionKey)
	return Result.SessionKey
}

func(this *SplunkQuery) SubmitJob(query string) string{
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	req := httplib.Post(this.baseUrl+ `/search/jobs`).SetTransport(tr)
	req.Param("search","search "+ query)
	req.Param("output_mode","json")
	req.Header("Authorization","Splunk " + this.sessionKey)
	var  Result SReault
	_ = req.ToJSON(&Result)
	fmt.Println(Result.Sid)
	return Result.Sid
}

type LogRecord struct {
	Raw		string  `json:"_raw"`
	Host	string  `json:"host"`
}

type LogResult struct{
	Offset	int64 			`json:"init_offset"`
	Results []LogRecord		`json:"results" `
}

func(this *SplunkQuery) QueryResults(sid string) []string{
	var Records []string
	tr := &http.Transport{TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},}
	req := httplib.Post(this.baseUrl+ `/search/jobs/` + sid).SetTransport(tr)
	req.Param("output_mode","json")
	req.Header("Authorization","Splunk " + this.sessionKey)
	result, _ := req.String()
	var jsonResult interface{}
	_ = json.Unmarshal([]byte(result), &jsonResult)
	content, err := jsonpath.JsonPathLookup(jsonResult, "$.entry[0].content")
	if err != nil{
		fmt.Println(err)
		return Records
	}

	_content := content.(map[string]interface{})
	var isDone = _content["isDone"].(bool)
	var eventCount = 0
	for {
		if isDone{
			eventCount = int(_content["eventCount"].(float64))
			break
		}
		req := httplib.Post(this.baseUrl+ `/search/jobs/` + sid).SetTransport(tr)
		req.Param("output_mode","json")
		req.Header("Authorization","Splunk " + this.sessionKey)
		result, _ := req.String()
		// fmt.Println(result)
		var jsonResult interface{}
		_ = json.Unmarshal([]byte(result), &jsonResult)
		content, _ := jsonpath.JsonPathLookup(jsonResult, "$.entry[0].content")
		_content = content.(map[string]interface{})
		isDone = _content["isDone"].(bool)
	}

	var iIndex = 0
	for {
		if iIndex > eventCount  {
			break
		}
		reqGet := httplib.Get(this.baseUrl+ `/search/jobs/` + sid+ "/results?output_mode=json&f=_raw&offset="+ strconv.Itoa(iIndex)).SetTransport(tr)
		reqGet.Header("Authorization","Splunk " + this.sessionKey)
		var logSet LogResult
		_ = reqGet.ToJSON(&logSet)
		// fmt.Println(logSet.Offset)
		for _, v := range logSet.Results {
			Records = append(Records, v.Raw)
		}
		iIndex = iIndex + 100
	}
	return  Records
}