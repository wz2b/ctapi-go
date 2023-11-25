package cdb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type AlarmSummaryRecord struct {
	SeqNo         uint64
	RecordTime    time.Time
	Source        string
	State         int
	StateDesc     *string
	ActiveTime    *time.Time
	InactiveTime  *time.Time
	AckTime       *time.Time
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
		 WHERE cs.RecordTime > {ts '%s'}
	 
		 ORDER BY cs.RecordTime`

func (this *CdbConnection) GetAlarmSummarySince(t time.Time) chan AlarmSummaryRecord {
	var compactQuery = fmt.Sprintf(strings.Join(strings.Fields(strings.TrimSpace(alarmSummaryQuery)), " "),
		t.UTC().Format("2006-01-02 15:04:05Z"))

	ch := make(chan AlarmSummaryRecord)

	go func() {
		rows, err := this.db.Query(compactQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to execute query: %s\n", err.Error())
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var record AlarmSummaryRecord
			if err := rows.Scan(
				&record.SeqNo,
				&record.RecordTime,
				&record.Source,
				&record.Message,
				&record.State,
				&record.StateDesc,
				&record.ActiveTime,
				&record.InactiveTime,
				&record.AckTime,
				&record.ConditionName,
				&record.Description,
				&record.DisplayName); err != nil {
				log.Fatal(err)
			}

			ch <- record

		}
		close(ch)
	}()

	return ch

}
