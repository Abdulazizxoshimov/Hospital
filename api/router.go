package api

import (
	"time"

	_ "github.com/Abdulazizxoshimov/Hospital/api/docs"
	"github.com/Abdulazizxoshimov/Hospital/api/handlers"
	"github.com/Abdulazizxoshimov/Hospital/api/middleware"
	"github.com/Abdulazizxoshimov/Hospital/config"
	"github.com/Abdulazizxoshimov/Hospital/internal/repo"
	redis "github.com/Abdulazizxoshimov/Hospital/internal/repo/redisdb"
	"github.com/Abdulazizxoshimov/Hospital/pkg/logger"
	"github.com/Abdulazizxoshimov/Hospital/pkg/token"
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteOption struct {
	Config         config.Config
	Logger         logger.Logger
	ContextTimeout time.Duration
	Cache          redis.Cache
	Enforcer       *casbin.Enforcer
	RefreshToken   token.JWTHandler
	Service        repo.StorageI
}

// NewRoute
// @Title Hospital 
// @Description Contacs: https://t.me/Abuzada0401
// @securityDefinitions.apikey BearerAuth
// @in 			header
// @name 		Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := handlers.New(&handlers.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Redis:          option.Cache,
		RefreshToken:   option.RefreshToken,
		Enforcer:       option.Enforcer,
		Service:        option.Service,
	})

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:7777"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, 
		AllowHeaders:     []string{"Content-Type", "Authorization"}, 
		AllowCredentials: true, 
	}
	
	router.Use(cors.New(corsConfig))
	
	router.Use(middleware.CheckCasbinPermission(option.Enforcer, option.Config))

	// login
	router.POST("/register", HandlerV1.Register)
	router.POST("/login", HandlerV1.Login)
	router.POST("/forgot/:email", HandlerV1.Forgot)
	router.POST("/verify", HandlerV1.VerifyOTP)
	router.PUT("/reset-password", HandlerV1.ResetPassword)
	router.GET("/token/:refresh", HandlerV1.Token)
	router.POST("/users/verify", HandlerV1.Verify)

	//user
	router.POST("/user", HandlerV1.CreateUser)
	router.PUT("/user", HandlerV1.UpdateUser)
	router.DELETE("/user/:id", HandlerV1.DeleteUser)
	router.GET("/user/:id", HandlerV1.GetUser)
	router.GET("/users", HandlerV1.ListUsers)
	router.PUT("/user/password", HandlerV1.UpdatePassword)

	//doctor
	router.POST("/doctor", HandlerV1.CreateDoctor)
	router.GET("/doctor/:id", HandlerV1.GetDoctor)
	router.PUT("/doctor", HandlerV1.UpdateDoctor)
	router.GET("/doctors", HandlerV1.ListDoctors)
	router.DELETE("/doctor/:id", HandlerV1.DeleteDoctor)

	//appointment
	router.POST("/appointment", HandlerV1.CreateAppointment)
	router.GET("/appointments", HandlerV1.GetAppointments)
	router.GET("/appointment/:id", HandlerV1.GetAppointmentByID)
	router.PUT("/appointment", HandlerV1.UpdateAppointment)
	router.DELETE("/appointment/:id", HandlerV1.DeleteAppointment)
	router.GET("/availabilities", HandlerV1.GetDoctorAvailabilities)
	router.GET("/availability/:id", HandlerV1.GetAvailabilityByID)


	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
