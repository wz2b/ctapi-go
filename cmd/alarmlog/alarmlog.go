package main

import (
	"fmt"
	"github.com/wz2b/ctapi-go/ctapi"
)

var dllPath = "C:\\Program Files (x86)\\AVEVA Plant SCADA\\Bin\\Bin (x64)"

func main() {
	api, err := ctapi.Init("CtApi.dll", dllPath)
	if err != nil {
		panic(err)
	}

	hConnection, err := api.CtOpen()
	if err != nil {
		panic(err)
	}
	defer api.CtClose(hConnection)

	const table = "AlarmSummary"

	for hObj := range api.FindAll(hConnection, table) {
		v, _ := api.GetPropertyAsString(hObj, "TAG")
		s, _ := api.GetPropertyAsString(hObj, "SUMSTATE")
		sd, _ := api.GetPropertyAsString(hObj, "STATE_DESC")
		onDate, _ := api.GetPropertyAsString(hObj, "ONDATE")
		onTime, _ := api.GetPropertyAsString(hObj, "ONTIME")
		offDate, _ := api.GetPropertyAsString(hObj, "OFFDATE")
		offTime, _ := api.GetPropertyAsString(hObj, "OFFTIME")
		ackDate, _ := api.GetPropertyAsString(hObj, "ACKDATE")
		ackTime, _ := api.GetPropertyAsString(hObj, "ACKTIME")

		fmt.Printf("TAG=\"%s\" SUMSTATE=\"%s\" STATE_DESC=\"%s\" ON=\"%s %s\" OFF=\"%s %s\" ACK=\"%s %s\"\n",
			v, s, sd, onDate, onTime, offDate, offTime, ackDate, ackTime)
	}
}
