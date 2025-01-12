package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	authCtr "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/auth/controller"
	authRepo "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/auth/repository"
	authSvc "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/auth/service"
	deptCtr "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/department/controller"
	deptRepo "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/department/repository"
	deptSvc "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/department/service"
	userCtr "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/user/controller"
	userRepo "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/user/repository"
	userSvc "github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/app/user/service"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/env"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/middlewares"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/bcrypt"
	errorhandler "github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/helpers/http/error_handler"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/helpers/http/response"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/jwt"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
	timePkg "github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/time"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/uuid"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/validator"
)

type HttpServer interface {
	Start(part string)
	MountMiddlewares()
	MountRoutes(db *sqlx.DB)
	GetApp() *fiber.App
}

type httpServer struct {
	app *fiber.App
}

func NewHttpServer() HttpServer {
	config := fiber.Config{
		CaseSensitive: true,
		AppName:       "GoGo Manager",
		ServerHeader:  "GoGo Manager",
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		ErrorHandler:  errorhandler.ErrorHandler,
	}

	app := fiber.New(config)

	return &httpServer{
		app: app,
	}
}

func (s *httpServer) GetApp() *fiber.App {
	return s.app
}

func (s *httpServer) Start(port string) {
	if port[0] != ':' {
		port = ":" + port
	}

	err := s.app.Listen(port)

	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[SERVER][Start] failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {
	s.app.Use(middlewares.LoggerConfig())
	s.app.Use(middlewares.Helmet())
	s.app.Use(middlewares.Compress())
	s.app.Use(middlewares.Cors())
	if env.AppEnv.AppEnv != "development" {
		s.app.Use(middlewares.ApiKey())
	}
	s.app.Use(middlewares.RecoverConfig())
}

func (s *httpServer) MountRoutes(db *sqlx.DB) {
	bcrypt := bcrypt.Bcrypt
	_ = timePkg.Time
	uuid := uuid.UUID
	validator := validator.Validator
	jwt := jwt.Jwt

	_ = middlewares.NewMiddleware(jwt)

	s.app.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "GoGoManager API")
	})

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, "GoGoManager API")
	})

	userRepository := userRepo.NewUserRepository(db)
	authRepository := authRepo.NewAuthRepository(db)
	departmentRepository := deptRepo.NewDepartmentRepository(db)

	userService := userSvc.NewUserService(userRepository, validator, uuid, bcrypt)
	authService := authSvc.NewAuthService(authRepository, validator, uuid, jwt, bcrypt)
	departmentService := deptSvc.NewDepartmentService(departmentRepository)

	userCtr.InitNewController(v1, userService)
	authCtr.InitAuthController(v1, authService)
	deptCtr.InitNewController(v1, departmentService)
	s.app.Use(func(c *fiber.Ctx) error {
		return c.SendFile("./web/not-found.html")
	})
}
