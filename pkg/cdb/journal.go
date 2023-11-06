package cdb

import (
	"database/sql"
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
		 WHERE cs.SeqNo > %d
	 
		 ORDER BY cs.RecordTime`

type JournalFilter func(JournalRecord) bool

func All(record JournalRecord) bool {
	return true
}

func AlarmsOnly(record JournalRecord) bool {
	return record.AlarmStateDesc != nil && len(strings.TrimSpace(*record.AlarmStateDesc)) > 0
}

func (this *CdbConnection) GetJournalSince(last uint64, filter JournalFilter) (chan JournalRecord, error) {
	var compactQuery = fmt.Sprintf(strings.Join(strings.Fields(strings.TrimSpace(journalQuery)), " "),
		last)

	fmt.Fprintf(os.Stderr, "Query: %s\n", compactQuery)

	ch := make(chan JournalRecord)

	go func() {
		rows, err := this.db.Query(compactQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to execute query: %s\n", err.Error())
			panic(err)
		}
		defer rows.Close()

		for rows.Next() {
			record := currentRecordToJournal(rows)
			if filter(*record) {
				ch <- *record
			}
		}
		close(ch)
	}()

	return ch, nil
}

func currentRecordToJournal(rows *sql.Rows) *JournalRecord {
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

	//
	// The time that comes in from the database is in some time zone, and
	// from experimentation on my system it is in UTC even though the
	// computer Plant SCADA is running on is set to local time.
	//
	record.RecordTime = time.Date(
		record.RecordTime.Year(),
		record.RecordTime.Month(),
		record.RecordTime.Day(),
		record.RecordTime.Hour(),
		record.RecordTime.Minute(),
		record.RecordTime.Second(),
		record.RecordTime.Nanosecond(),
		time.UTC)

	return &record
}
