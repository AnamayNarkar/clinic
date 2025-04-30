package routes

import (
	"clinic/sqlc"
	"clinic/src/controllers"
	"clinic/src/middleware"
	"clinic/src/security"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetUpAllRoutes(ginEngine *gin.Engine, db *sqlc.Queries, redisClient *redis.Client, securityManager *security.SecurityManager) {
	authController := controllers.NewAuthController(db, redisClient)
	patientController := controllers.NewPatientController(db, redisClient)
	doctorController := controllers.NewDoctorController(db, redisClient)
	applicationController := controllers.NewApplicationController(db)
	appointmentController := controllers.NewAppointmentController(db)

	apiGroup := ginEngine.Group("/api")

	setUpAuthRoutes(apiGroup, authController)
	setupPatientRoutes(apiGroup, patientController, securityManager)
	setupDoctorRoutes(apiGroup, doctorController, securityManager)
	setupApplicationRoutes(apiGroup, applicationController, redisClient, securityManager)
	setupAppointmentRoutes(apiGroup, appointmentController, redisClient, securityManager)
}

// Application route setup
func setupApplicationRoutes(ginGroup *gin.RouterGroup, appController *controllers.ApplicationController, redisClient *redis.Client, securityManager *security.SecurityManager) {
	g := ginGroup.Group("/application")
	g.Use(middleware.GetSession(redisClient))
	g.POST("/", securityManager.GeneralPermissionMiddleware("application", "create"), appController.CreateApplication)
	g.GET("/:id", securityManager.GeneralPermissionMiddleware("application", "read"), appController.GetApplication)
	g.GET("/", securityManager.GeneralPermissionMiddleware("application", "read"), appController.ListApplications)
	g.PUT("/:id", securityManager.GeneralPermissionMiddleware("application", "update"), appController.UpdateApplication)
}

// Appointment route setup
func setupAppointmentRoutes(ginGroup *gin.RouterGroup, appointmentController *controllers.AppointmentController, redisClient *redis.Client, securityManager *security.SecurityManager) {
	g := ginGroup.Group("/appointment")
	g.Use(middleware.GetSession(redisClient))
	g.POST("/", securityManager.GeneralPermissionMiddleware("appointment", "create"), appointmentController.CreateAppointment)
	g.GET("/:id", securityManager.GeneralPermissionMiddleware("appointment", "read"), appointmentController.GetAppointment)
	g.GET("/", securityManager.GeneralPermissionMiddleware("appointment", "read"), appointmentController.ListAppointments)
	g.DELETE("/:id", securityManager.GeneralPermissionMiddleware("appointment", "delete"), appointmentController.CancelAppointment)
}

func setUpAuthRoutes(ginGroup *gin.RouterGroup, authController *controllers.AuthController) {
	authGroup := ginGroup.Group("/auth")
	{
		authGroup.POST("/login/:role", authController.Login)
		authGroup.POST("/logout", authController.Logout)
	}
}

func setupPatientRoutes(ginGroup *gin.RouterGroup, patientController *controllers.PatientController, securityManager *security.SecurityManager) {
	ginGroup = ginGroup.Group("/patient")
	ginGroup.POST("/", patientController.CreatePatient)
}

func setupDoctorRoutes(ginGroup *gin.RouterGroup, doctorController *controllers.DoctorController, securityManager *security.SecurityManager) {
	ginGroup = ginGroup.Group("/doctor")
	ginGroup.POST("/", securityManager.GeneralPermissionMiddleware("doctor", "create"), doctorController.CreateDoctor)
}
