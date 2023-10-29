package cdb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type JournalRecord struct {
	SeqNo             uint64
	RecordTime        time.Time
	Source            string
	Message           *string
	CustomStringField *string
	AlarmStateDesc    *string
	Description       *string
	DisplayName       *string
}

const journalQuery = `
		SELECT cs.SeqNo,
		       cs.RecordTime,
		       cs.Source,
		       cs.Message,
		       cs.CustomStringField,
		       cs.AlarmStateDesc,
		       COALESCE( da.Description, adva.Description, aa.Description, msa.Description, tsa.Description ) Description,
		       COALESCE( da.DisplayName, adva.DisplayName, aa.DisplayName, msa.DisplayName, tsa.DisplayName ) DisplayName
		 FROM CDBEventJournal cs
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

type JournalFilter func(JournalRecord) bool

func All(record JournalRecord) bool {
	return true
}

func AlarmsOnly(record JournalRecord) bool {
	if record.AlarmStateDesc == nil {
		return false
	} else {
		switch strings.ToLower(*record.AlarmStateDesc) {
		case "cleared":
			return true
		case "acknowledged":
			return true
		case "raised":
			return true
		default:
			return false
		}
	}
}

func (this *CdbConnection) GetJournalSince(t time.Time, filter JournalFilter) (chan JournalRecord, error) {
	var compactQuery = fmt.Sprintf(strings.Join(strings.Fields(strings.TrimSpace(journalQuery)), " "),
		t.UTC().Format("2006-01-02 15:04:05Z"))

	fmt.Println(compactQuery)

	ch := make(chan JournalRecord)

	go func() {
		rows, err := this.db.Query(compactQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to execute query: %s\n", err.Error())
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			var record JournalRecord
			if err := rows.Scan(
				&record.SeqNo,
				&record.RecordTime,
				&record.Source,
				&record.Message,
				&record.CustomStringField,
				&record.AlarmStateDesc,
				&record.Description,
				&record.DisplayName); err != nil {
				log.Fatal(err)
			}

			if filter(record) {
				ch <- record
			}

		}
		close(ch)
	}()

	return ch, nil

}
