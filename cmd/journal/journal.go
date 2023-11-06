package main

import (
	"encoding/json"
	"fmt"
	"golink/pkg/cdb"
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

	entries, err := db.GetJournalSince(0, filter)
	if err != nil {
		panic(err)
	}

	for entry := range entries {
		b, _ := json.Marshal(entry)
		fmt.Println(string(b))
	}

}
