package routes

import (
	"clinic/sqlc"
	"clinic/src/controllers"
	"clinic/src/security"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetUpAllRoutes(ginEngine *gin.Engine, db *sqlc.Queries, redisClient *redis.Client, securityManager *security.SecurityManager) {
	authController := controllers.NewAuthController(db, redisClient)
	patientController := controllers.NewPatientController(db, redisClient)
	doctorController := controllers.NewDoctorController(db, redisClient)

	apiGroup := ginEngine.Group("/api")

	setUpAuthRoutes(apiGroup, authController)
	setupPatientRoutes(apiGroup, patientController, securityManager)
	setupDoctorRoutes(apiGroup, doctorController, securityManager)

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
