package interfaces

// DeviceService defines the interface for device operations
type DeviceService interface {
	SendCommand(deviceID string, command string) error
}
