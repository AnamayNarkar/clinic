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
	appointmentController := controllers.NewAppointmentController(db)
	applicationController := controllers.NewApplicationController(db, appointmentController)
	receptionistController := controllers.NewReceptionistController(db)

	apiGroup := ginEngine.Group("/api")

	setUpAuthRoutes(apiGroup, authController)
	setupPatientRoutes(apiGroup, patientController, securityManager)
	setupDoctorRoutes(apiGroup, doctorController, redisClient, securityManager)
	setupReceptionistRoutes(apiGroup, receptionistController, redisClient, securityManager)
	setupApplicationRoutes(apiGroup, applicationController, redisClient, securityManager)
	setupAppointmentRoutes(apiGroup, appointmentController, redisClient, securityManager)
}

func setupApplicationRoutes(ginGroup *gin.RouterGroup, appController *controllers.ApplicationController, redisClient *redis.Client, securityManager *security.SecurityManager) {
	g := ginGroup.Group("/application")
	g.Use(middleware.GetSession(redisClient))
	g.POST("/", securityManager.GeneralPermissionMiddleware("application", "create"), appController.CreateApplication)
	g.GET("/", securityManager.GeneralPermissionMiddleware("application", "read"), appController.ListApplications)
	g.PUT("/", securityManager.GeneralPermissionMiddleware("application", "update"), appController.UpdateApplication)
}

func setupAppointmentRoutes(ginGroup *gin.RouterGroup, appointmentController *controllers.AppointmentController, redisClient *redis.Client, securityManager *security.SecurityManager) {
	g := ginGroup.Group("/appointment")
	g.Use(middleware.GetSession(redisClient))

	// Get upcoming appointments
	g.GET("/upcoming", securityManager.GeneralPermissionMiddleware("appointment", "read"), appointmentController.GetUpcomingAppointments)

	g.PUT("/doctor/:id", securityManager.GeneralPermissionMiddleware("appointment", "update"), appointmentController.UpdateAppointmentByDoctor)
}

func setUpAuthRoutes(ginGroup *gin.RouterGroup, authController *controllers.AuthController) {
	authGroup := ginGroup.Group("/auth")

	// Public routes (no session required)
	authGroup.POST("/login/:role", authController.Login)
	authGroup.POST("/logout", authController.Logout)

}

func setupPatientRoutes(ginGroup *gin.RouterGroup, patientController *controllers.PatientController, securityManager *security.SecurityManager) {
	ginGroup = ginGroup.Group("/patient")
	ginGroup.POST("/", patientController.CreatePatient)
}

func setupDoctorRoutes(ginGroup *gin.RouterGroup, doctorController *controllers.DoctorController, redisClient *redis.Client, securityManager *security.SecurityManager) {
	g := ginGroup.Group("/doctor")
	g.Use(middleware.GetSession(redisClient))
	g.POST("/", securityManager.GeneralPermissionMiddleware("doctor", "create"), doctorController.CreateDoctor)

	doctorAPI := ginGroup.Group("/doctors")
	doctorAPI.Use(middleware.GetSession(redisClient))

	doctorAPI.GET("/appointments", doctorController.GetDoctorAppointments)
}

func setupReceptionistRoutes(ginGroup *gin.RouterGroup, receptionistController *controllers.ReceptionistController, redisClient *redis.Client, securityManager *security.SecurityManager) {
	g := ginGroup.Group("/receptionist")
	g.Use(middleware.GetSession(redisClient))
	g.POST("/", securityManager.GeneralPermissionMiddleware("receptionist", "create"), receptionistController.CreateReceptionist)
}
