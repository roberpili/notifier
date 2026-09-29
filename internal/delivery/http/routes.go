package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/rober0xf/notifier/docs"
	"github.com/rober0xf/notifier/pkg/auth"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(userHandler *UserHandler, paymentHandler *PaymentHandler, authMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           48 * time.Hour,
	}))

	v1 := r.Group("/v1")
	admin := v1.Group("/admin")
	protected := v1.Group("/auth")

	protected.Use(authMiddleware)
	admin.Use(authMiddleware, auth.AdminOnly())

	setupAdminOnlyRoutes(admin, userHandler, paymentHandler)
	setupPublicUsersRoutes(v1, userHandler)
	setupProtectedUsersRoutes(protected, userHandler)
	setupPaymentsRoutes(protected, paymentHandler)

	// frontend
	r.Static("/static", "./frontend/dist")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	return r
}

func setupPublicUsersRoutes(v1 *gin.RouterGroup, userHandler *UserHandler) {
	users := v1.Group("/users")

	users.GET("/email_verification/:token", userHandler.VerifyEmail)
	users.POST("/register", userHandler.Create)
	users.POST("/login", userHandler.Login)
	users.POST("/login/google", userHandler.GoogleLogin)
}

func setupProtectedUsersRoutes(protected *gin.RouterGroup, userHandler *UserHandler) {
	users := protected.Group("/users")

	users.GET("/me", userHandler.Me)
	users.PUT("/:id", userHandler.Update)
}

func setupAdminOnlyRoutes(admin *gin.RouterGroup, userHandler *UserHandler, paymentHandler *PaymentHandler) {
	users := admin.Group("/users")
	payments := admin.Group("/payments")

	users.GET("", userHandler.GetAll)
	users.GET("/email/:email", userHandler.GetByEmail)
	users.GET("/:id", userHandler.GetByID)
	users.DELETE("/:id", userHandler.Delete)

	payments.GET("", paymentHandler.GetAll)
	payments.GET("/:id", paymentHandler.GetByID)
}

func setupPaymentsRoutes(protected *gin.RouterGroup, paymentHandler *PaymentHandler) {
	payments := protected.Group("/payments")

	payments.GET("/me", paymentHandler.GetMyPayments)
	payments.POST("", paymentHandler.Create)
	payments.PUT("/:id", paymentHandler.Update)
	payments.DELETE("/:id", paymentHandler.Delete)
}
