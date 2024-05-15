package services

import (
	"log"

	"PragatiIot/platform/models"
	"PragatiIot/platform/repositories"
)

type RoleService struct {
	roleRepo *repositories.RoleRepository
}

func NewRoleService(roleRepo *repositories.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) GetRoleByName(roleName string) (models.Role, error) {
	role, err := s.roleRepo.GetRoleByName(roleName)
	if err != nil {
		log.Printf("Error getting role by name: %v", err)
		return role, err
	}
	return role, nil
}
