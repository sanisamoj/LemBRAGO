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

	router.Use(middlewares.NewRateLimiterMiddleware(time.Minute, 100))

	router.StaticFile("/favicon.ico", "./uploads/favicon.ico")
	router.Static("/uploads", "./uploads")

	adminOnly := []models.UserRole{models.RoleAdmin}
	adminOrMember := []models.UserRole{models.RoleAdmin, models.RoleMember}

	public := router.Group("/")
	public.Use(middlewares.NewRateLimiterMiddleware(time.Minute, 100))
	{
		public.POST("/organizations", controllers.CreateOrganization)
		public.POST("/auth", middlewares.DictionaryPreviewMiddleware(), controllers.SendAuthCode)
		public.POST("/login", controllers.GetLoginInfoFromUser)
		public.POST("/environment/login", controllers.UserLogin)
		public.GET("/invites/:id", controllers.GetInvitedCodeToken)
	}

	media := router.Group("/media")
	{
		media.GET("/:filename", controllers.HandleServeFile)
		media.POST("", middlewares.AuthMiddleware(appConfig.JWTSecret, []models.UserRole{}), controllers.HandleUploadFile)
	}

	organization := router.Group("/org")
	organization.Use(
		middlewares.NewRateLimiterMiddleware(time.Minute, 100),
		middlewares.AuthMiddleware(appConfig.JWTSecret, []models.UserRole{models.RoleAdmin}),
	)
	{
		organization.GET("/users", controllers.GetUsers)
		organization.DELETE("/users", controllers.DeleteUser)
	}

	invites := router.Group("/invites")
	invites.Use(middlewares.NewRateLimiterMiddleware(time.Minute, 100))
	{
		invites.POST("", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.InviteUser)
	}

	creation := router.Group("/users/creation")
	creation.Use(
		middlewares.NewRateLimiterMiddleware(time.Minute, 100),
		middlewares.AuthMiddleware(appConfig.JWTSecretUserCreation, []models.UserRole{}),
	)
	{
		creation.POST("", controllers.UserRegister)
	}

	user := router.Group("/users")
	user.Use(
		middlewares.NewRateLimiterMiddleware(time.Minute, 100),
		middlewares.AuthMiddleware(appConfig.JWTSecret, []models.UserRole{}),
	)
	{
		user.GET("/vaults", controllers.GetMyVaultsByOrgID)
	}

	vaults := router.Group("/vaults")
	vaults.Use(middlewares.NewRateLimiterMiddleware(time.Minute, 100))
	{
		vaults.POST("", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.CreateVault)
		vaults.PUT("", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.UpdateVault)
		vaults.DELETE("/:id", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.RemoveVault)

		vaults.GET("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.GetAllMembersFromTheVault)
		vaults.POST("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.AddMemberToVault)
		vaults.DELETE("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.RemoveMemberFromTheVault)
		vaults.PUT("/members", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.UpdateMemberPermission)

		vaults.GET("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.GetAllPasswordsFromVault)
		vaults.POST("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.CreatePassword)
		vaults.PUT("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.UpdatePasswordInVault)
		vaults.DELETE("/passwords", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.DeletePassword)

		vaults.GET("/medias", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOnly), controllers.GetAllMediasFromTheOrg)
		vaults.DELETE("/medias", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.DeleteMedia)
	}

	versions := router.Group("/versions")
	versions.Use(middlewares.NewRateLimiterMiddleware(time.Minute, 100))
	{
		versions.GET("", controllers.GetAllVersions)
		versions.GET("/latest", controllers.GetLatestAppVersion)
		versions.POST("", middlewares.AuthMiddleware(appConfig.JWTSecretAdmin, adminOnly), controllers.RegisterVersion)
		versions.PUT("", middlewares.AuthMiddleware(appConfig.JWTSecretAdmin, adminOnly), controllers.UpdateVersion)

		versions.GET("/:target/:arch/:version", controllers.DownloadDesktopApp)
		versions.POST("/desktop", middlewares.AuthMiddleware(appConfig.JWTSecretAdmin, adminOnly), controllers.UploadDesktopApp)
	}

	router.DELETE("/signout", middlewares.AuthMiddleware(appConfig.JWTSecret, adminOrMember), controllers.Signout)

	host := appConfig.Host
	port := appConfig.Port
	rt := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("Server running on", rt)
	router.Run(rt)
}
