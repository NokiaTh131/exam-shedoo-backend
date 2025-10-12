package handlers

import (
	"shedoo-backend/internal/app/auth"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	rawUser := c.Locals("user")
	role := c.Locals("role")
	if rawUser == nil || role == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"ok":      false,
			"message": "unauthorized",
		})
	}

	session, ok := rawUser.(auth.SessionData)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "invalid user type",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"cmuitaccount":         session.BasicInfo.Cmuitaccount,
		"cmuitaccount_name":    session.BasicInfo.CmuitaccountName,
		"firstname_EN":         session.BasicInfo.FirstnameEN,
		"firstname_TH":         session.BasicInfo.FirstnameTH,
		"itaccounttype_EN":     session.BasicInfo.ItaccounttypeEN,
		"itaccounttype_TH":     session.BasicInfo.ItaccounttypeTH,
		"itaccounttype_id":     session.BasicInfo.ItaccounttypeID,
		"lastname_EN":          session.BasicInfo.LastnameEN,
		"lastname_TH":          session.BasicInfo.LastnameTH,
		"organization_name_EN": session.BasicInfo.OrganizationNameEN,
		"organization_name_TH": session.BasicInfo.OrganizationNameTH,
		"role":                 role,
		"student_id":           session.BasicInfo.StudentID,
	})
}
