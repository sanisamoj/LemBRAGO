package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/services"
)

func GetLatestAppVersion(c *gin.Context) {
	version, err := services.GetLatestVersion()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, version)
}

func GetAllVersions(c *gin.Context) {
	versions, err := services.GetAllVersions()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, versions)
}

func RegisterVersion(c *gin.Context) {
	var req models.ApplicationVersion
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := services.RegisterVersion(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, res)
}

func UpdateVersion(c *gin.Context) {
	var req models.ApplicationVersion
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	version, err := services.UpdateVersion(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, version)
}

func DownloadDesktopApp(c *gin.Context) {
	target := c.Param("target")
	arch := c.Param("arch")
	version := c.Param("version")
	filename := fmt.Sprintf("lembrago_%s_%s_en-US.msi", version, arch)

	filePath := filepath.Join("releases", fmt.Sprintf("%s-%s", target, arch), filename)
	fmt.Println(filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	c.FileAttachment(filePath, filename)
}
