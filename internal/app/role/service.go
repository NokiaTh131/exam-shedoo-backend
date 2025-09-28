package admin

import "shedoo-backend/internal/repositories"

type RoleService struct {
	AdminRepo *repositories.AdminRepository
}

func NewRoleService(repo *repositories.AdminRepository) *RoleService {
	return &RoleService{AdminRepo: repo}
}

func (s *RoleService) ClassifyRole(accountTypeID, accountName string) (string, error) {
	isAdmin, err := s.AdminRepo.IsAdmin(accountName)
	if err != nil {
		return "", err
	}

	if isAdmin {
		return "admin", nil
	}

	if accountTypeID == "StdAcc" {
		return "student", nil
	}

	return "professor", nil
}
