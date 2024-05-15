package mqtt

import (
	"encoding/json"
	"log"

	"PragatiIot/platform/models"
	"PragatiIot/platform/rabbitmq"
	"PragatiIot/platform/services"
)

type MQTTHandler struct {
	deviceService *services.DeviceService
	producer      *rabbitmq.Producer
}

func NewMQTTHandler(deviceService *services.DeviceService, producer *rabbitmq.Producer) *MQTTHandler {
	return &MQTTHandler{
		deviceService: deviceService,
		producer:      producer,
	}
}

func (h *MQTTHandler) ProcessMessage(deviceID string, message []byte) error {
	device, err := h.deviceService.GetDeviceByID(deviceID)
	if err != nil {
		log.Printf("Error finding device %s: %v", deviceID, err)
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(message, &data); err != nil {
		log.Printf("Error parsing MQTT message for device %s: %v", deviceID, err)
		return err
	}

	deviceData := models.DeviceData{
		DeviceID: device.DeviceID,
		HomeID:   device.HomeID,
		Data:     data,
	}
	if err := h.deviceService.AddDeviceData(deviceData); err != nil {
		log.Printf("Error storing data for device %s: %v", device.DeviceID, err)
		return err
	}

	// Publish to RabbitMQ
	messageBytes, err := json.Marshal(deviceData)
	if err != nil {
		log.Printf("Error marshalling device data for RabbitMQ: %v", err)
		return err
	}

	if err := h.producer.Publish(messageBytes); err != nil {
		log.Printf("Error publishing device data to RabbitMQ: %v", err)
		return err
	}

	return nil
}
