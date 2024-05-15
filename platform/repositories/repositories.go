package repositories

import (
	"context"
	"fmt"

	"PragatiIot/platform/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeviceRepository struct {
	pool *pgxpool.Pool
}

type HomeRepository struct {
	pool *pgxpool.Pool
}

type RoleRepository struct {
	pool *pgxpool.Pool
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{pool: pool}
}

func (r *RoleRepository) GetRoleByName(roleName string) (models.Role, error) {
	var role models.Role
	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, name FROM roles WHERE name = $1`,
		roleName,
	).Scan(&role.ID, &role.Name)
	if err != nil {
		return role, fmt.Errorf("error finding role by name %s: %w", roleName, err)
	}
	return role, nil
}

func NewHomeRepository(pool *pgxpool.Pool) *HomeRepository {
	return &HomeRepository{pool: pool}
}

func (r *HomeRepository) AddHome(home models.Home) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO homes (home_name, user_id) VALUES ($1, $2)`,
		home.HomeName, home.UserID,
	)
	return err
}

func (r *HomeRepository) GetHomesByUserID(userID int) ([]models.Home, error) {
	rows, err := r.pool.Query(
		context.Background(),
		`SELECT id, home_name, user_id, created_at FROM homes WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error finding homes by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	var homes []models.Home
	for rows.Next() {
		var home models.Home
		if err := rows.Scan(&home.ID, &home.HomeName, &home.UserID, &home.CreatedAt); err != nil {
			return nil, err
		}
		homes = append(homes, home)
	}
	return homes, nil
}

func (r *HomeRepository) AddUserToHome(homeUser models.HomeUser) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO home_users (home_id, user_id, role_id) VALUES ($1, $2, $3)`,
		homeUser.HomeID, homeUser.UserID, homeUser.RoleID,
	)
	return err
}

func (r *HomeRepository) GetHomeUserRole(homeID, userID int) (int, error) {
	var roleID int
	err := r.pool.QueryRow(
		context.Background(),
		`SELECT role_id FROM home_users WHERE home_id = $1 AND user_id = $2`,
		homeID, userID,
	).Scan(&roleID)
	if err != nil {
		return 0, fmt.Errorf("error finding role for user %d in home %d: %w", userID, homeID, err)
	}
	return roleID, nil
}

func NewDeviceRepository(pool *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{pool: pool}

}

func (r *DeviceRepository) GetDeviceByID(deviceID string) (models.Device, error) {
	var device models.Device
	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, device_id, channel_id, production_date, warranty, location, is_active, user_id, home_id, created_at
		 FROM devices
		 WHERE device_id = $1`,
		deviceID,
	).Scan(
		&device.ID, &device.DeviceID, &device.ChannelID, &device.ProductionDate, &device.Warranty,
		&device.Location, &device.IsActive, &device.UserID, &device.HomeID, &device.CreatedAt,
	)
	if err != nil {
		return device, fmt.Errorf("error finding device by ID %s: %w", deviceID, err)
	}
	return device, nil
}

func (r *DeviceRepository) AddDeviceData(deviceData models.DeviceData) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO device_data (device_id, home_id, data) VALUES ($1, $2, $3)`,
		deviceData.DeviceID, deviceData.HomeID, deviceData.Data,
	)
	return err
}

func (r *DeviceRepository) GetDevicesByUserID(userID int) ([]models.Device, error) {
	rows, err := r.pool.Query(
		context.Background(),
		`SELECT id, device_id, channel_id, production_date, warranty, location, is_active, user_id, home_id, created_at
		 FROM devices
		 WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("error finding devices by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		if err := rows.Scan(
			&device.ID, &device.DeviceID, &device.ChannelID, &device.ProductionDate, &device.Warranty,
			&device.Location, &device.IsActive, &device.UserID, &device.HomeID, &device.CreatedAt,
		); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) AddUser(user models.User) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3)`,
		user.Username, user.PasswordHash, user.Email,
	)
	if err != nil {
		return fmt.Errorf("error adding user %s: %w", user.Username, err)
	}
	return nil
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, username, password_hash, email FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email)
	if err != nil {
		return user, fmt.Errorf("error finding user by username %s: %w", username, err)
	}
	return user, nil
}
func (r *DeviceRepository) AddDevice(device models.Device) error {
	_, err := r.pool.Exec(
		context.Background(),
		`INSERT INTO devices (device_id, channel_id, production_date, warranty, location, is_active, user_id, home_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		device.DeviceID, device.ChannelID, device.ProductionDate, device.Warranty, device.Location,
		device.IsActive, device.UserID, device.HomeID, device.CreatedAt,
	)
	return err
}

func (r *DeviceRepository) UpdateDevice(device models.Device) error {
	_, err := r.pool.Exec(
		context.Background(),
		`UPDATE devices SET channel_id = $2, production_date = $3, warranty = $4, location = $5, is_active = $6, user_id = $7, home_id = $8, created_at = $9
		WHERE device_id = $1`,
		device.DeviceID, device.ChannelID, device.ProductionDate, device.Warranty, device.Location,
		device.IsActive, device.UserID, device.HomeID, device.CreatedAt,
	)
	return err
}

func (r *DeviceRepository) GetDeviceByChannel(channelID string) (models.Device, error) {
	var device models.Device
	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, device_id, channel_id, production_date, warranty, location, is_active, user_id, home_id, created_at
		FROM devices
		WHERE channel_id = $1`,
		channelID,
	).Scan(&device.ID, &device.DeviceID, &device.ChannelID, &device.ProductionDate, &device.Warranty,
		&device.Location, &device.IsActive, &device.UserID, &device.HomeID, &device.CreatedAt)
	if err != nil {
		return device, fmt.Errorf("error finding device by channel ID %s: %w", channelID, err)
	}
	return device, nil
}
