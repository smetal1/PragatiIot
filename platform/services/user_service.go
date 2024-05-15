package services

import (
	"log"

	"PragatiIot/platform/models"
	"PragatiIot/platform/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) AddUser(user models.User) error {
	if err := s.userRepo.AddUser(user); err != nil {
		log.Printf("Error adding user: %v", err)
		return err
	}
	return nil
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user by username %s: %v", username, err)
		return user, err
	}
	return user, nil
}
