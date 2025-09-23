package auth

import (
	"context"
	"time"

	"shedoo-backend/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo      *repositories.AuthRepository
	jwtSecret string
}

func NewAuthService(repo *repositories.AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) SignIn(ctx context.Context, code string) (string, error) {
	accessToken, err := s.repo.ExchangeCode(ctx, code)
	if err != nil {
		return "", err
	}

	basicInfo, err := s.repo.GetBasicInfo(ctx, accessToken)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"cmuitaccount_name":    basicInfo.CmuitaccountName,
		"cmuitaccount":         basicInfo.Cmuitaccount,
		"student_id":           basicInfo.StudentID,
		"firstname_EN":         basicInfo.FirstnameEN,
		"lastname_EN":          basicInfo.LastnameEN,
		"organization_name_EN": basicInfo.OrganizationNameEN,
		"exp":                  time.Now().Add(time.Hour).Unix(),
		"iat":                  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
