package main

import (
	"fmt"
	_ "github.com/alexbrainman/odbc"
	log "github.com/sirupsen/logrus"
	"github.com/wz2b/ctapi-go/cdb"
)

func main() {

	if db, err := cdb.Open(); err != nil {
		panic(err)
	} else {
		defer db.Close()
		rows, err := db.GetJournalSince(0, cdb.AlarmsOnly)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Fatal("Unable to get journal")
		}
		for r := range rows {

			if r.AlarmStateDesc == nil {
				log.WithFields(log.Fields{"seq": r.SeqNo}).Error("no state")
			} else if *r.AlarmStateDesc == "Raised" {
				fmt.Printf("%s  %d (%s)\n",
					r.Source,
					r.SeqNo,
					*r.AlarmStateDesc)
			} else if *r.AlarmStateDesc == "Cleared" || *r.AlarmStateDesc == "Acknowledged" {
				orig, err := db.GetLastRaisedForSource(&r)
				if err != nil {
					log.WithFields(log.Fields{"error": err}).Fatal("error finding source")
				} else if orig != nil {
					fmt.Printf("%s  %d (%s) -> %d (%s)\n",
						r.Source,
						orig.SeqNo,
						orig.RecordTime,
						r.SeqNo,
						*r.AlarmStateDesc)
				} else {
					fmt.Printf("%s  ??? (???) -> %d (%s)\n",
						r.Source,
						r.SeqNo,
						*r.AlarmStateDesc)
				}

			}
		}

	}
}
