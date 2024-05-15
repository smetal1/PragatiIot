package services

import (
	"log"

	"PragatiIot/platform/models"
	"PragatiIot/platform/repositories"
)

type DeviceService struct {
	deviceRepo  *repositories.DeviceRepository
	homeService *HomeService
}

func NewDeviceService(deviceRepo *repositories.DeviceRepository, homeService *HomeService) *DeviceService {
	return &DeviceService{deviceRepo: deviceRepo, homeService: homeService}
}

func (s *DeviceService) AddDevice(device models.Device) error {
	if err := s.deviceRepo.AddDevice(device); err != nil {
		log.Printf("Error adding device: %v", err)
		return err
	}
	return nil
}

func (s *DeviceService) AssignDeviceToHome(deviceID string, homeID *int) error {
	device, err := s.deviceRepo.GetDeviceByID(deviceID)
	if err != nil {
		log.Printf("Error finding device %s: %v", deviceID, err)
		return err
	}

	device.HomeID = homeID
	if err := s.deviceRepo.UpdateDevice(device); err != nil {
		log.Printf("Error updating device %s: %v", deviceID, err)
		return err
	}
	return nil
}

func (s *DeviceService) GetDevicesByUserID(userID int) ([]models.Device, error) {
	devices, err := s.deviceRepo.GetDevicesByUserID(userID)
	if err != nil {
		log.Printf("Error getting devices by user ID: %v", err)
		return nil, err
	}
	return devices, nil
}

func (s *DeviceService) GetDeviceByID(deviceID string) (models.Device, error) {
	device, err := s.deviceRepo.GetDeviceByID(deviceID)
	if err != nil {
		log.Printf("Error finding device by ID %s: %v", deviceID, err)
		return device, err
	}
	return device, nil
}

func (s *DeviceService) AddDeviceData(deviceData models.DeviceData) error {
	if err := s.deviceRepo.AddDeviceData(deviceData); err != nil {
		log.Printf("Error adding device data: %v", err)
		return err
	}
	return nil
}

func (s *DeviceService) GetDeviceAnalytics(deviceID string, homeID, roleID int) ([]models.DeviceData, error) {
	// Implement analytics retrieval logic here
	var analytics []models.DeviceData
	// Example logic (to be replaced with actual implementation):
	log.Printf("Fetching analytics for device %s in home %d with role %d", deviceID, homeID, roleID)
	return analytics, nil
}
func (s *DeviceService) GetDeviceByChannel(channelID string) (models.Device, error) {
	device, err := s.deviceRepo.GetDeviceByChannel(channelID)
	if err != nil {
		log.Printf("Error finding device by channel ID %s: %v", channelID, err)
		return device, err
	}
	return device, nil
}
