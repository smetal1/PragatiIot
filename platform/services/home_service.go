package services

import (
	"log"

	"PragatiIot/platform/models"
	"PragatiIot/platform/repositories"
)

type HomeService struct {
	homeRepo    *repositories.HomeRepository
	roleService *RoleService
	userService *UserService
}

func NewHomeService(homeRepo *repositories.HomeRepository, roleService *RoleService) *HomeService {
	return &HomeService{homeRepo: homeRepo, roleService: roleService}
}

func (s *HomeService) AddHome(home models.Home) error {
	if err := s.homeRepo.AddHome(home); err != nil {
		log.Printf("Error adding home: %v", err)
		return err
	}
	return nil
}

func (s *HomeService) AddUserToHome(homeID, userID int, roleName string) error {
	role, err := s.roleService.GetRoleByName(roleName)
	if err != nil {
		log.Printf("Error getting role by name: %v", err)
		return err
	}

	homeUser := models.HomeUser{
		HomeID: homeID,
		UserID: userID,
		RoleID: role.ID,
	}
	if err := s.homeRepo.AddUserToHome(homeUser); err != nil {
		log.Printf("Error adding user to home: %v", err)
		return err
	}
	return nil
}

func (s *HomeService) GetHomesByUserID(userID int) ([]models.Home, error) {
	homes, err := s.homeRepo.GetHomesByUserID(userID)
	if err != nil {
		log.Printf("Error getting homes by user ID: %v", err)
		return nil, err
	}
	return homes, nil
}

func (s *HomeService) GetHomeUserRole(homeID, userID int) (int, error) {
	roleID, err := s.homeRepo.GetHomeUserRole(homeID, userID)
	if err != nil {
		log.Printf("Error getting role for user %d in home %d: %v", userID, homeID, err)
		return 0, err
	}
	return roleID, nil
}

func (s *HomeService) GetUserByUsername(username string) (models.User, error) {
	user, err := s.userService.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user by username %s: %v", username, err)
		return user, err
	}
	return user, nil
}
