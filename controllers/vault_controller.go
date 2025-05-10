package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/services"
	"lembrago.com/lembrago/utils"
)

func CreateVault(c *gin.Context) {
	var req models.CreateVaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

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
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid Role"})
		return
	}

	vault, err := services.CreateVault(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, vault)
}

func RemoveVault(c *gin.Context) {
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

	vaultID := c.Param("id")
	if vaultID == "" {
		c.JSON(400, gin.H{"error": "vaultId is required"})
		return
	}

	err := services.RemoveVault(userID, vaultID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Vault, members and passwords removed successfully"})
}

func GetMyVaultsByOrgID(c *gin.Context) {
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

	orgIDRaw, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}

	orgID, ok := orgIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (orgID type)"})
		return
	}

	vaults, err := services.FindMyAllVaultsByOrgID(userID, orgID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, vaults)
}

func AddMemberToVault(c *gin.Context) {
	var req models.CreateVaultMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

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

	vaultWithMemberResponse, err := services.AddMemberToVault(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, vaultWithMemberResponse)
}

func RemoveMemberFromTheVault(c *gin.Context) {
	memberId := c.Query("id")
	if memberId == "" {
		c.JSON(400, gin.H{"error": "memberId is required"})
		return
	}

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

	err := services.RemoveMemberFromVault(userID, memberId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"message": "Member removed from the vault"})
}

func UpdateMemberPermission(c *gin.Context) {
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

	orgIDRaw, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	orgID, ok := orgIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (orgID type)"})
		return
	}

	var req models.UpdateVaultMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = services.UpdateMemberPermission(userID, orgID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func CreatePassword(c *gin.Context) {
	var req models.CreatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

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

	password, err := services.AddPasswordToVault(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, password)
}

func DeletePassword(c *gin.Context) {
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

	passwordId := c.Query("id")
	if passwordId == "" {
		c.JSON(400, gin.H{"error": "passwordId is required"})
		return
	}

	err := services.DeletePasswordFromVault(userID, passwordId)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

func UpdatePasswordInVault(c *gin.Context) {
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

	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	err := utils.GetValidator().Struct(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pRes, err := services.UpdatePasswordInVault(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, pRes)
}

func GetAllPasswordsFromVault(c *gin.Context) {
	vaultId := c.Query("vaultId")
	if vaultId == "" {
		c.JSON(400, gin.H{"error": "vaultId is required"})
		return
	}

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

	passwords, err := services.GetAllPasswordsFromVault(userID, vaultId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, passwords)
}
