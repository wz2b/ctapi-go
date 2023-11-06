package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type JournalToMqttConfig struct {
	certificate struct {
		cert string
		ca   string
		key  string
	}

	mqtt struct {
		endpoint string
		clientId string
	}
}

func loadConfig() *JournalToMqttConfig {

	var config JournalToMqttConfig

	viper.AddConfigPath(".")
	viper.SetConfigName("config") // name of config file (without extension)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	pflag.String("cert", "", "path to private (device) certificate")
	pflag.String("key", "", "path to key file for private certificate")
	pflag.String("ca", "", "path to CA certificate")

	pflag.String("broker", "", "URL to broker endpoint")
	pflag.String("client-id", "", "MQTT client id")

	pflag.Parse()

	viper.BindPFlag("certificate.cert", pflag.Lookup("cert"))
	viper.BindPFlag("certificate.ca", pflag.Lookup("key"))
	viper.BindPFlag("certificate.key", pflag.Lookup("ca"))

	viper.BindPFlag("mqtt.endpoint", pflag.Lookup("broker"))
	viper.BindPFlag("mqtt.clientId", pflag.Lookup("client-id"))

	config.certificate.cert = viper.GetString("certificate.cert")
	config.certificate.key = viper.GetString("certificate.key")
	config.certificate.ca = viper.GetString("certificate.ca")

	config.mqtt.endpoint = viper.GetString("mqtt.endpoint")
	config.mqtt.clientId = viper.GetString("mqtt.clientId")

	return &config
}
