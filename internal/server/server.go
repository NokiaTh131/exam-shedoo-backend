package server

import (
	"context"
	"time"

	"shedoo-backend/internal/app/auth"
	"shedoo-backend/internal/app/course"
	"shedoo-backend/internal/app/courseexam"
	"shedoo-backend/internal/app/enrollment"
	admin "shedoo-backend/internal/app/role"
	scrapejobs "shedoo-backend/internal/app/scrape_jobs"
	"shedoo-backend/internal/config"
	"shedoo-backend/internal/handlers"
	"shedoo-backend/internal/middlewares"
	"shedoo-backend/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
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
	AdminHandler      *handlers.AdminHandler
	AuthMiddleware    *middlewares.AuthMiddleware
}

func New() *FiberServer {
	dbService := config.NewDB()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.NewRedis().Addr,
		Password: config.NewRedis().Password,
		DB:       0,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic("cannot connect to Redis: " + err.Error())
	}

	// === Enrollment ===
	enrollmentRepo := repositories.NewEnrollmentRepository(dbService.DB)
	enrollmentService := enrollment.NewEnrollmentService(enrollmentRepo)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService)

	// === Course Exams ===
	course_examsRepo := repositories.NewCourseExamRepository(dbService.DB)
	course_examsService := courseexam.NewCourseExamService(course_examsRepo)
	course_examsHandler := handlers.NewCourseExamHandler(course_examsService)

	// === Courses ===
	courseRepo := repositories.NewCourseRepository(dbService.DB)
	courseService := course.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseService)

	// === Scrape Jobs ===
	scrapeJobRepo := repositories.NewScrapeJobRepository(dbService.DB)
	scrapeJobService := scrapejobs.NewScrapeJobService(scrapeJobRepo)
	scrapeJobHandler := handlers.NewScrapeJobHandler(scrapeJobService)

	// === Auth ===
	authRepo := repositories.NewAuthRepository(
		config.LoadAuthConfig().TokenURL,
		config.LoadAuthConfig().RedirectURL,
		config.LoadAuthConfig().ClientID,
		config.LoadAuthConfig().ClientSecret,
		config.LoadAuthConfig().Scope,
		config.LoadAuthConfig().BasicInfoURL,
	)

	accessTokenTTL := 24 * time.Hour
	refreshTokenTTL := 7 * 24 * time.Hour
	issuer := "shedoo"
	audience := "shedoo-users"

	authService := auth.NewAuthService(
		authRepo,
		config.LoadAuthConfig().JWTSecret,
		redisClient,
		accessTokenTTL,
		refreshTokenTTL,
		issuer,
		audience,
	)
	authHandler := handlers.NewAuthHandler(
		authService,
		config.LoadAuthConfig().CookieDomain,
		config.LoadAuthConfig().IsProd,
	)

	// === Admin ===
	adminRepo := repositories.NewAdminRepository(dbService.DB)
	roleService := admin.NewRoleService(adminRepo)
	adminHandler := handlers.NewAdminHandler(roleService)

	userHandler := handlers.NewUserHandler()

	authMiddleware := middlewares.NewAuthMiddleware(roleService, redisClient)

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
		AdminHandler:      adminHandler,
		AuthMiddleware:    authMiddleware,
	}

	return server
}
