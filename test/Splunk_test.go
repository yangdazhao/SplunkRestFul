package test

import (
	"SplunkRestFul"
	"fmt"
	"testing"
)

type JobParam struct{
	Earliest 		string `json:"earliest"`
	Latest 			string `json:"latest"`
	SerialNumber    string `json:"serialnumber"`
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

func TestGetKey(t*testing.T){
	query := SplunkRestFul.NewSplunkQuery("dazhao.yang", "dazhao.yang")
	queryStr := NewJobParamEx("7f57b169cbb14199a46f3db0926f282d","8/2/2019:17:00:00","8/3/2019:17:00:00"	).ToString()
	sid := query.SubmitJob(	queryStr)
	logSet :=  query.QueryResults(sid)
	for _, v := range logSet {
		fmt.Println(v)
	}
}
