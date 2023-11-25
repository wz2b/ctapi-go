package cdb

import (
	"database/sql"
	"errors"
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
		SELECT %s cs.SeqNo,
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
		WHERE %s
	    ORDER BY %s`

type JournalFilter func(JournalRecord) bool

func All(record JournalRecord) bool {
	return true
}

func AlarmsOnly(record JournalRecord) bool {
	return record.AlarmStateDesc != nil && len(strings.TrimSpace(*record.AlarmStateDesc)) > 0
}

func makeJournalQuery(whereClause string, topClause string, orderClause string) string {
	cq := strings.Join(strings.Fields(strings.TrimSpace(journalQuery)), " ")
	return fmt.Sprintf(cq, topClause, whereClause, orderClause)
}

func (this *CdbConnection) GetJournalSince(last uint64, filter JournalFilter) (chan JournalRecord, error) {
	var qs = makeJournalQuery(
		fmt.Sprintf("cs.SeqNo > %d", last),
		"",
		"cs.RecordTime")

	fmt.Fprintf(os.Stderr, "----\nQuery journal: %s\n", qs)

	ch := make(chan JournalRecord)

	go func() {
		rows, err := this.db.Query(qs)
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

func (this *CdbConnection) GetLastRaisedForSource(source *JournalRecord) (*JournalRecord, error) {
	whereClause := fmt.Sprintf("cs.Source='%s' AND cs.SeqNo < %d AND cs.AlarmStateDesc='Raised'",
		source.Source,
		source.SeqNo)
	var qs = makeJournalQuery(
		whereClause,
		"TOP (1)",
		"cs.RecordTime DESC")

	//fmt.Fprintf(os.Stderr, "----\nQuery for source: %s\n", qs)

	rows, err := this.db.Query(qs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() == false {
		// not found
		return nil, nil
	}

	parent := currentRecordToJournal(rows)

	if rows.Next() {
		return parent, errors.New("fetching associated journal entry, more than one row was returned, but only one was expected")
	}
	return parent, nil
}
