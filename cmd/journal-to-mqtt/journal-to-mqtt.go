package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golink/pkg/cdb"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"time"
)

var config *JournalToMqttConfig
var mqttConnection *IotConnection

func main() {
	var err error
	f := &lumberjack.Logger{
		Filename:  "agent.log",
		MaxSize:   10, /* max size per log file, megabytes */
		MaxAge:    14, /* max days to save old log files before removing them */
		Compress:  true,
		LocalTime: true,
	}

	// don't forget to close it
	defer f.Close()

	log.SetFormatter(&log.TextFormatter{ForceQuote: true})

	// Output to stderr instead of stdout, could also be a file.
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	config = loadConfig()
	mqttConnection, err = NewIotConnection(config)
	if err != nil {
		panic(err)
	}

	db, err := cdb.Open()
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	var last uint64 = 0

	for {

		entries, err := db.GetJournalSince(last, cdb.AlarmsOnly)
		if err != nil {
			log.WithError(err).Error("Unable to get journal")
		}

		for entry := range entries {
			if entry.AlarmStateDesc != nil {
				switch strings.ToLower(*entry.AlarmStateDesc) {
				case "raised":
					//sendJournalEntry(entry)

				case "cleared":
					//sendJournalEntry(entry)

				case "acknowledged":
					//sendJournalEntry(entry)
				}
			}
			if entry.SeqNo > last {
				last = entry.SeqNo
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func sendJournalEntry(entry cdb.JournalRecord) (err error) {
	const topic = "$aws/rules/fclab_alarm_input"

	if err = mqttConnection.Connect(); err == nil {
		msgBytes, _ := json.Marshal(entry)
		source := strings.TrimSpace(entry.Source)
		if source != "" {
			log.WithFields(log.Fields{"seq": entry.SeqNo}).Info("Publishing")
			mqttConnection.Publish(topic, 1, false, msgBytes)
		} else {
			log.Error("Not publishing journal entry with no source")

		}
	}
	return
}
