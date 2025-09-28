package server

import (
	"shedoo-backend/internal/app/auth"
	"shedoo-backend/internal/app/course"
	"shedoo-backend/internal/app/courseexam"
	"shedoo-backend/internal/app/enrollment"
	admin "shedoo-backend/internal/app/role"
	scrapejobs "shedoo-backend/internal/app/scrape_jobs"
	"shedoo-backend/internal/config"
	"shedoo-backend/internal/handlers"
	"shedoo-backend/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
	EnrollmentHandler *handlers.EnrollmentHandler
	CourseHandler     *handlers.CourseHandler
	CourseExamHandler *handlers.CourseExamHandler
	ScrapeJobHandler  *handlers.ScrapeJobHandler
	AuthHandler       *handlers.AuthHandler
	UserHandler       *handlers.UserHandler
	RoleService       *admin.RoleService
}

func New() *FiberServer {
	dbService := config.New()
	enrollmentRepo := repositories.NewEnrollmentRepository(dbService.DB)
	enrollmentService := enrollment.NewEnrollmentService(enrollmentRepo)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)

	course_examsRepo := repositories.NewCourseExamRepository(dbService.DB)
	course_examsService := courseexam.NewCourseExamService(course_examsRepo)
	course_examsHandler := handlers.NewCourseExamHandler(course_examsService)

	courseRepo := repositories.NewCourseRepository(dbService.DB)
	courseService := course.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseService)

	scrapeJobRepo := repositories.NewScrapeJobRepository(dbService.DB)
	scrapeJobService := scrapejobs.NewScrapeJobService(scrapeJobRepo)
	scrapeJobHandler := handlers.NewScrapeJobHandler(scrapeJobService)

	authRepo := repositories.NewAuthRepository(
		config.LoadAuthConfig().TokenURL,
		config.LoadAuthConfig().RedirectURL,
		config.LoadAuthConfig().ClientID,
		config.LoadAuthConfig().ClientSecret,
		config.LoadAuthConfig().Scope,
		config.LoadAuthConfig().BasicInfoURL,
	)
	authService := auth.NewAuthService(authRepo, config.LoadAuthConfig().JWTSecret)
	authHandler := handlers.NewAuthHandler(authService, config.LoadAuthConfig().CookieDomain, config.LoadAuthConfig().IsProd)

	userHandler := handlers.NewUserHandler()

	adminRepo := repositories.NewAdminRepository(dbService.DB)
	roleService := admin.NewRoleService(adminRepo)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "shedoo-backend",
			AppName:      "shedoo-backend",
		}),
		EnrollmentHandler: enrollmentHandler,
		CourseHandler:     courseHandler,
		CourseExamHandler: course_examsHandler,
		ScrapeJobHandler:  scrapeJobHandler,
		AuthHandler:       authHandler,
		UserHandler:       userHandler,
		RoleService:       roleService,
	}

	return server
}
