package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/services"
	"lembrago.com/lembrago/utils"
)

func CreateOrganization(c *gin.Context) {
	var req models.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validação falhou", "details": err.Error()})
		return
	}

	org, err := services.CreateOrganization(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, org)
}

func InviteUser(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (userID type)"})
		return
	}

	orgIDRaw, exists := c.Get("orgID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	orgID, ok := orgIDRaw.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro interno (userID type)"})
		return
	}

	var req models.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validação falhou", "details": err.Error()})
		return
	}

	token, err := services.InviteUser(userID, orgID, req.Email, req.Role)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"invitedCode": token})
}

func GetInvitedCodeToken(c *gin.Context) {
	code := c.Param("id")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	minOrgWithTokenResponse, err := services.GetInviteCodeToken(code)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, minOrgWithTokenResponse)
}
