package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/cache"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/services"
	"lembrago.com/lembrago/utils"
)

func SendAuthCode(c *gin.Context) {
	var reqData, exists = c.Get("validatedAuthCodeRequest")
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "Request data not found in context"})
		return
	}

	req, ok := reqData.(models.AuthCodeRequest)
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "Request data type assertion failed"})
		return
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = services.SendAuthCode(req.Email)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Code sent successfully"})
}

func GetLoginInfoFromUser(c *gin.Context) {
	var req models.AuthCodeSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Code == "" {
		c.JSON(400, gin.H{"error": "Code is empty"})
		return
	}

	loginInfo, err := services.GetLoginInfoFromUser(req.Email, req.Code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, loginInfo)
}

func UserLogin(c *gin.Context) {
	var req models.UserLoginComparison
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UserLogin(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, user)
}

func UserRegister(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	orgIDRaw, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	orgID, ok := orgIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (orgID type)"})
		return
	}

	roleRaw, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, ok := roleRaw.(models.UserRole)
	fmt.Println(role)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (role type)"})
		return
	}

	emailRaw, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, ok := emailRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (email type)"})
		return
	}

	exists, err = cache.Exists(req.Code)
	if err != nil {
		c.Error(err)
		return
	}
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	userRole := models.UserRole(role)
	err = services.UserRegister(&req, orgID, email, userRole)
	if err != nil {
		c.Error(err)
		return
	}
	cache.Delete(req.Code)

	c.Status(http.StatusCreated)
}

func DeleteUser(c *gin.Context) {
	ownerID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	strOwnerID, ok := ownerID.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (userID type)"})
		return
	}

	userID := c.Query("userId")
	if userID == strOwnerID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	err := services.DeleteUser(strOwnerID, userID)
	if err != nil {
		c.Error(err)
	}

	c.Status(200)
}

func GetUsers(c *gin.Context) {
	id := c.Query("userId")

	orgIDRaw, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	orgID, ok := orgIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (orgID type)"})
		return
	}

	roleRaw, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role, ok := roleRaw.(models.UserRole)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (role type)"})
		return
	}
	if role != models.RoleAdmin {
		c.AbortWithStatusJSON(403, gin.H{"error": "Invalid Permission"})
		return
	}

	if id != "" {
		user, err := services.GetUserByID(id)
		if err != nil {
			c.Error(err)
			return
		}
		if user.OrgID != orgID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.JSON(200, user)
		return
	} else {
		users, err := services.GetUsersByOrgID(orgID)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, users)
	}
}

func GetAllMembersFromTheVault(c *gin.Context) {
	vaultId := c.Query("vaultId")

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (userID type)"})
		return
	}

	members, err := services.GetAllMembersFromTheVault(vaultId, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, members)
}

func Signout(g *gin.Context) {
	authHeader := g.GetHeader("Authorization")
	token, _ := utils.GetTokenInH(authHeader)
	services.SignOut(token)
	g.Status(200)
}
