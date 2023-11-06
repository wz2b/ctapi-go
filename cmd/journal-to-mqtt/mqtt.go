package main

import (
	"crypto/tls"
	"crypto/x509"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type IotConnection struct {
	options *mqtt.ClientOptions
	client  mqtt.Client
}

func NewIotConnection(config *JournalToMqttConfig) (*IotConnection, error) {
	var c IotConnection

	clientCertificate, err := tls.LoadX509KeyPair(
		config.certificate.cert, config.certificate.key)
	if err != nil {
		return nil, WrapError("unable to load private certificate", err)
	}

	certs := x509.NewCertPool()

	// Load the Amazon CA certificate
	caPem, err := os.ReadFile(config.certificate.ca)
	if err != nil {
		panic(err)
	}

	certs.AppendCertsFromPEM(caPem)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCertificate},
		RootCAs:      certs,
	}

	c.options = mqtt.NewClientOptions()
	c.options.AddBroker("tcps://" + config.mqtt.endpoint + ":8883/mqtt")
	c.options.SetMaxReconnectInterval(10 * time.Second)
	c.options.SetClientID(config.mqtt.clientId)
	c.options.SetTLSConfig(tlsConfig)
	c.options.AutoReconnect = true
	c.options.MaxReconnectInterval = 1 * time.Minute
	c.options.ConnectRetry = true

	// This sets a ping timeout that tests the connection
	// every so often. This defaults to 30 seconds
	// in PAHO but we're going to increase it to five minutes.
	c.options.KeepAlive = 300 // this is seconds

	listenPath := config.mqtt.clientId + "/cmd"

	c.options.OnConnect = func(conn mqtt.Client) {
		log.Info("connected")
		if token := conn.Subscribe(listenPath,
			byte(1),
			c.onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	c.options.OnConnectionLost = func(conn mqtt.Client, err error) {
		log.WithFields(log.Fields{"error": err, "message": err.Error()}).Error("Connection lost")
	}

	c.options.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		log.Info("websocket is reconnecting")
	}

	c.client = mqtt.NewClient(c.options)

	return &c, nil
}

func (this *IotConnection) Connect() error {

	if this.client.IsConnected() == false {
		if token := this.client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}

	}

	return nil
}

func (this *IotConnection) Close() error {
	if this.client.IsConnected() == true {
		this.client.Disconnect(500)
	}

	return nil
}

func (this *IotConnection) Publish(topic string, qos byte, retained bool, data []byte) error {
	if token := this.client.Publish(topic, qos, retained, data); token.Wait() && token.Error() != nil {
		println("ERROR")
		return token.Error()
	}
	return nil
}
