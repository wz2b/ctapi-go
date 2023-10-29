package main

import (
	"encoding/json"
	"fmt"
	"golink/pkg/cdb"
	"time"
)

func main() {

	db, err := cdb.Open()
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	filter := cdb.AlarmsOnly

	entries, err := db.GetJournalSince(time.Now().AddDate(0, -1, 0), filter)
	if err != nil {
		panic(err)
	}

	for entry := range entries {
		b, _ := json.Marshal(entry)
		fmt.Println(string(b))
	}

}
