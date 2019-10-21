package main

import (
	"github.com/splunk/splunk-cloud-sdk-go/sdk"
	"github.com/splunk/splunk-cloud-sdk-go/services"
	"os"
)

func main(){
	//query := SplunkRestFul.NewSplunkQuery("dazhao.yang", "dazhao.yang")
	//logSet :=  query.QueryResults(query.SubmitJob(`index="ysb-syslog" `+ os.Args[1] +` |reverse`))
	//for _, v := range logSet {
	//	fmt.Println(v)
	//}

	client, err := sdk.NewClient(&services.Config{
		Token:  os.Getenv("BEARER_TOKEN"),
		Tenant: os.Getenv("TENANT"),
	})
}