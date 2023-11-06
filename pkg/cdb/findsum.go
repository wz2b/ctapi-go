package cdb

import (
	"fmt"
	"os"
	"time"
)

const findAcked = `
	SELECT asum.RecordTime
	 FROM CDBEventJournal cs
	 INNER JOIN CDBAlarmSummary asum
	   	ON asum.AckTime=cs.RecordTime
		AND asum.Source=cs.Source
	 WHERE asum.AckTime={ts ?}
	    AND cs.Source=?
		AND cs.AlarmStateDesc='Acknowledged'`

const findCleared = `
	SELECT asum.RecordTime
	 FROM CDBEventJournal cs
	 INNER JOIN CDBAlarmSummary asum
		ON asum.InactiveTime=cs.RecordTime
		AND asum.Source=cs.Source
	 WHERE asum.InactiveTime={ts ?}
	    AND cs.Source=?
	 	AND cs.AlarmStateDesc='Cleared'`

func (this *CdbConnection) FindAssociatedJournalEntry(ack JournalRecord) *time.Time {

	var query string
	if ack.AlarmStateDesc == nil {
		return nil
	} else if *ack.AlarmStateDesc == "Cleared" {
		query = findCleared
	} else if *ack.AlarmStateDesc == "Acknowledge" {
		query = findAcked
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Unknown record type %s\n", *ack.AlarmStateDesc)
		return nil
	}

	row := this.db.QueryRow(query,
		ack.RecordTime.UTC().Format("2006-01-02 15:04:05.999Z"),
		ack.Source)

	if row == nil {
		return nil
	}

	var t time.Time
	err := row.Scan(&t)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error parsing time: %s\n", err.Error())
		return nil
	}

	t = forceUtc(t)
	return &t
}
