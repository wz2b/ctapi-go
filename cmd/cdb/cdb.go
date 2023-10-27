package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"log"
	"strings"
	"time"
)

type AlarmSummaryRecord struct {
	SeqNo         uint64
	RecordTime    string
	Source        string
	State         int
	StateDesc     *string
	ActiveTime    *string
	InactiveTime  *string
	AckTime       *string
	Message       *string
	ConditionName *string
	DisplayName   *string
	Description   *string
}

const alarmSummaryQuery = `
		SELECT cs.SeqNo,
		       cs.RecordTime,
		       cs.Source,
		       cs.Message,
		       cs.State,
		       cs.StateDesc,
		       cs.ActiveTime,
		       cs.InactiveTime,
		       cs.AckTime,
		       cs.ConditionName,
		       COALESCE( da.Description, adva.Description, aa.Description, msa.Description, tsa.Description ) Description,
		       COALESCE( da.DisplayName, adva.DisplayName, aa.DisplayName, msa.DisplayName, tsa.DisplayName ) DisplayName
		 FROM CDBAlarmSummary cs
		 LEFT OUTER JOIN CiDigitalAlarm da
		    ON cs.Id=da.Id
		 LEFT OUTER JOIN CiAdvancedAlarm adva
		    ON cs.Id=adva.Id
		 LEFT OUTER JOIN CiAnalogAlarm aa
		    ON cs.Id=aa.Id
		LEFT OUTER JOIN CiMultiStateDigAlarm msa
		    ON cs.Id=msa.Id		    	
		LEFT OUTER JOIN CiTimestampedAlarm tsa
		    ON cs.Id=tsa.Id	 
		 WHERE cs.SeqNo > %d
	 
		 ORDER BY cs.SeqNo`

func main() {

	var compactQuery = strings.Join(strings.Fields(strings.TrimSpace(alarmSummaryQuery)), " ")

	db, err := sql.Open("odbc", "dsn=CitectAlarms;uid=;pwd=")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var last uint64 = 0
	for {
		rows, err := db.Query(fmt.Sprintf(compactQuery, last))
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var result AlarmSummaryRecord
			if err := rows.Scan(
				&result.SeqNo,
				&result.RecordTime,
				&result.Source,
				&result.Message,
				&result.State,
				&result.StateDesc,
				&result.ActiveTime,
				&result.InactiveTime,
				&result.AckTime,
				&result.ConditionName,
				&result.Description,
				&result.DisplayName); err != nil {
				log.Fatal(err)
			}

			bytes, _ := json.Marshal(result)
			fmt.Println(string(bytes))

			last = result.SeqNo
			fmt.Printf("last is now %d\n", last)

		}

		time.Sleep(1 * time.Second)
	}

}
