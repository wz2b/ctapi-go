package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"github.com/wz2b/ctapi-go/cdb"
	"time"
)

func main() {

	if db, err := cdb.Open(); err != nil {
		panic(err)
	} else {
		defer db.Close()
		rows := db.GetAlarmSummarySince(time.Now().AddDate(0, -1, 0))

		for r := range rows {
			b, _ := json.Marshal(r)
			fmt.Println(string(b))
		}

	}
}
