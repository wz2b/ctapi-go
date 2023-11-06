package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type CommandMessage struct {
	CommandText string   `json:"command"`
	Args        []string `json:"args"`
}

func (this *IotConnection) onMessageReceived(client mqtt.Client, message mqtt.Message) {

	var cmd CommandMessage
	err := json.Unmarshal(message.Payload(), &cmd)

	if err == nil {
		switch cmd.CommandText {

		default:
			log.Warnf("Invalid remote command '%s' ignored", cmd.CommandText)
		}
	}
}
