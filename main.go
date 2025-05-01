package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/controllers"
	"lembrago.com/lembrago/handlers"
	"lembrago.com/lembrago/internal/config"
	"lembrago.com/lembrago/middlewares"
	"lembrago.com/lembrago/models"
)

func main() {
	appConfig := config.GetServerConfig()

	router := gin.Default()
	router.Use(handlers.ErrorHandler())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Use(middlewares.NewRateLimiterMiddleware(time.Second, 3))

	adminOnly := []models.UserRole{models.RoleAdmin}
	adminOrMember := []models.UserRole{models.RoleAdmin, models.RoleMember}

	router.POST("/organizations", controllers.CreateOrganization)
	router.POST("/auth", middlewares.DictionaryPreviewMiddleware(), controllers.SendAuthCode)
	router.POST("/login", controllers.GetLoginInfoFromUser)
	router.POST("/environment/login", controllers.UserLogin)

	organizationRoute := router.Group("/org")
	organizationRoute.Use(middlewares.AuthMiddleware(appConfig.JWTSecret, []models.UserRole{models.RoleAdmin}))
	{
		organizationRoute.GET("/users", controllers.GetUsers)
	}

	inviteRoute := router.Group("/invites")
	{
		inviteRoute.POST("", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.InviteUser)
		inviteRoute.GET("/:id", controllers.GetInvitedCodeToken)
	}

	creationUserRoute := router.Group("/users/creation")
	creationUserRoute.Use(middlewares.AuthMiddleware(appConfig.JWTSecretUserCreation, []models.UserRole{}))
	{
		creationUserRoute.POST("", controllers.UserRegister)
	}

	userRoute := router.Group("/users")
	userRoute.Use(middlewares.AuthMiddleware(appConfig.JWTSecret, []models.UserRole{}))
	{
		userRoute.GET("/vaults", controllers.GetMyVaultsByOrgID)
	}

	vaultRoute := router.Group("/vaults")
	{
		vaultRoute.POST("", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.CreateVault)
		vaultRoute.DELETE("/:id", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.RemoveVault)
		vaultRoute.GET("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.GetAllMembersFromTheVault)
		vaultRoute.POST("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.AddMemberToVault)
		vaultRoute.DELETE("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.RemoveMemberFromTheVault)

		vaultRoute.PUT("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.UpdateMemberPermission)

		vaultRoute.GET("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.GetAllPasswordsFromVault)
		vaultRoute.POST("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.CreatePassword)
		vaultRoute.PUT("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.UpdatePasswordInVault)
		vaultRoute.DELETE("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.DeletePassword)
	}

	router.DELETE("/signout", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.Signout)

	host := appConfig.Host
	port := appConfig.Port
	rt := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("Server running on", rt)
	router.Run(rt)
}
