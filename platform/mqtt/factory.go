package mqtt

import (
	"PragatiIot/platform/rabbitmq"
	"PragatiIot/platform/services"
	"fmt"
)

type ProtocolHandler interface {
	ProcessMessage(deviceID string, message []byte) error
}

type ProtocolFactory struct {
	deviceService *services.DeviceService
	producer      *rabbitmq.Producer
}

func NewProtocolFactory(deviceService *services.DeviceService, producer *rabbitmq.Producer) *ProtocolFactory {
	return &ProtocolFactory{deviceService: deviceService, producer: producer}
}

func (f *ProtocolFactory) CreateHandler(protocol string) (ProtocolHandler, error) {
	switch protocol {
	case "mqtt":
		return NewMQTTHandler(f.deviceService, f.producer), nil
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}
