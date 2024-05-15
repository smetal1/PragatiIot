package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"PragatiIot/platform/rabbitmq"
	"PragatiIot/platform/services"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	mu            sync.Mutex
	deviceService *services.DeviceService
	producer      *rabbitmq.Producer
	factory       *ProtocolFactory
	mqttClient    mqtt.Client
	topics        map[string]struct{}
}

func NewMQTTClient(deviceService *services.DeviceService, producer *rabbitmq.Producer, factory *ProtocolFactory) *MQTTClient {
	return &MQTTClient{
		deviceService: deviceService,
		producer:      producer,
		factory:       factory,
		topics:        make(map[string]struct{}),
	}
}

func (c *MQTTClient) handleMessage(client mqtt.Client, msg mqtt.Message) {
	device, err := c.deviceService.GetDeviceByChannel(msg.Topic())
	if err != nil {
		log.Printf("Error getting device for topic %s: %v", msg.Topic(), err)
		return
	}

	handler, err := c.factory.CreateHandler("mqtt")
	if err != nil {
		log.Printf("Error creating handler: %v", err)
		return
	}

	if err := handler.ProcessMessage(device.DeviceID, msg.Payload()); err != nil {
		log.Printf("Error processing message: %v", err)
	}
}

func (c *MQTTClient) StartMQTT(broker, clientID, caCert, clientCert, clientKey string, insecure bool) {
	var tlsConfig *tls.Config
	if insecure {
		// Configure to skip certificate validation
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else {
		// Use normal TLS configuration with certificates
		tlsConfig = c.newTLSConfig(caCert, clientCert, clientKey)
	}

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetTLSConfig(tlsConfig).
		SetDefaultPublishHandler(c.handleMessage)

	c.mqttClient = mqtt.NewClient(opts)
	if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection error: %v", token.Error())
	}

	go c.subscribeTopics()

	log.Println("MQTT client connected and ready")
	select {}
}

func (c *MQTTClient) newTLSConfig(caCert, clientCert, clientKey string) *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCert)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	certpool.AppendCertsFromPEM(ca)

	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("Failed to load client certificate/key pair: %v", err)
	}

	return &tls.Config{
		RootCAs:            certpool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
	}
}

func (c *MQTTClient) subscribeTopics() {
	for {
		devices, err := c.deviceService.GetDevicesByUserID(1) // or other userID if required
		if err != nil {
			log.Printf("Error getting devices: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		c.mu.Lock()
		for _, device := range devices {
			if _, ok := c.topics[device.ChannelID]; !ok {
				if token := c.mqttClient.Subscribe(device.ChannelID, 0, nil); token.Wait() && token.Error() != nil {
					log.Printf("Error subscribing to topic %s: %v", device.ChannelID, token.Error())
					continue
				}
				c.topics[device.ChannelID] = struct{}{}
				log.Printf("Subscribed to new topic: %s", device.ChannelID)
			}
		}
		c.mu.Unlock()

		time.Sleep(10 * time.Second)
	}
}
