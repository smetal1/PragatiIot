package models

import "time"

// User model
// User defines the structure for an API user.
// swagger:model User
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
}

// Role model
// Role represents user roles within the system.
// swagger:model Role
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Home model
// Home represents a household or location managed by a user.
// swagger:model Home
type Home struct {
	ID        int       `json:"id"`
	HomeName  string    `json:"home_name"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// HomeUser model
// HomeUser links a user with a home and assigns a role.
// swagger:model HomeUser
type HomeUser struct {
	HomeID int `json:"home_id"`
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

// Device model
// Device represents a physical or virtual device within the system.
// swagger:model Device
type Device struct {
	ID             int       `json:"id"`
	DeviceID       string    `json:"device_id"`
	ChannelID      string    `json:"channel_id"`
	ProductionDate time.Time `json:"production_date"`
	Warranty       int       `json:"warranty"`
	Location       string    `json:"location"`
	IsActive       bool      `json:"is_active"`
	UserID         int       `json:"user_id"`
	HomeID         *int      `json:"home_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// DeviceData model
// DeviceData represents data generated or consumed by a device.
// swagger:model DeviceData
type DeviceData struct {
	DeviceID string                 `json:"device_id"`
	HomeID   *int                   `json:"home_id"`
	Data     map[string]interface{} `json:"data"`
}

// ApiResponse model
// ApiResponse represents a standard response for API endpoints.
// swagger:model ApiResponse
type ApiResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Error   string `json:"error,omitempty"`
	token   string `json:"error,omitempty"`
}

// AddUserToHomeRequest model
// AddUserToHomeRequest defines the required parameters to add a user to a home.
// swagger:model AddUserToHomeRequest
type AddUserToHomeRequest struct {
	HomeID int    `json:"home_id"`
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}

// AssignDeviceRequest model
// AssignDeviceRequest defines the JSON structure for assigning a device to a home.
// swagger: model AssignDeviceRequest
type AssignDeviceRequest struct {
	DeviceID string `json:"device_id"`
	HomeID   *int   `json:"home_id"`
}
