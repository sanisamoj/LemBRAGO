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

func UploadDesktopApp(c *gin.Context) {
	target := c.PostForm("target")  
	arch := c.PostForm("arch")
	version := c.PostForm("version")
	lang := c.DefaultPostForm("lang", "en-US")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
		return
	}

	dir := fmt.Sprintf("releases/%s-%s", target, arch)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar pasta"})
		return
	}

	ext := filepath.Ext(file.Filename) // ex: .msi, .AppImage, .dmg
	filename := fmt.Sprintf("lembrago_%s_%s_%s%s", version, arch, lang, ext)
	filepath := filepath.Join(dir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upload feito com sucesso", "path": filepath})
}

func DownloadDesktopApp(c *gin.Context) {
	target := c.Param("target")
	arch := c.Param("arch")
	version := c.Param("version")

	if version == "latest" {
		v, _ := services.GetLatestVersion()
		version = v.LatestDesktopVersion.Version
    }

	filename := fmt.Sprintf("lembrago_%s_%s_en-US.msi", version, arch)

	filePath := filepath.Join("releases", fmt.Sprintf("%s-%s", target, arch), filename)
	fmt.Println(filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	c.FileAttachment(filePath, filename)
}
