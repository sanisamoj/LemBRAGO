package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/services"
)

const (
	uploadDir   = "./uploads"     // Pasta onde os arquivos serão salvos
	maxFileSize = 2 * 1024 * 1024 // Limite de 2 MB
	formFileKey = "media"         // Nome do campo no formulário multipart
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

func HandleUploadFile(c *gin.Context) {
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

	file, header, err := c.Request.FormFile(formFileKey)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	if header.Size > maxFileSize {
		errStr := fmt.Sprintf("The file has exceeded the maximum allowed size of %d bytes", maxFileSize)
		c.JSON(400, gin.H{"error": errStr})
		return
	}

	// --- 3. Validar o tipo do arquivo (MIME type) ---
	// A forma mais segura é detectar o tipo pelo conteúdo, não só pela extensão.
	// Lemos os primeiros 512 bytes para detectar o tipo MIME.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("err reading file: %v", err)})
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detectedContentType := http.DetectContentType(buffer)
	if !allowedImageTypes[detectedContentType] {
		allowedTypesStr := make([]string, 0, len(allowedImageTypes))
		for k := range allowedImageTypes {
			allowedTypesStr = append(allowedTypesStr, k)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid Type (%s). Allowed: %s",
				detectedContentType, strings.Join(allowedTypesStr, ", ")),
		})
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := filepath.Base(header.Filename)
	rndFilename := fmt.Sprintf("%s-%s", filename, primitive.NewObjectID().Hex()) + ext
	rndFilename = strings.ReplaceAll(rndFilename, " ", "_")
	rndFilename = strings.ToLower(rndFilename)

	svMedia, err := services.SaveMedia(orgID, rndFilename, header, header.Size, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, svMedia)
}


func HandleServeFile(c *gin.Context) {
	filename := c.Param("filename")

	safeFilename := filepath.Base(filename)
	if safeFilename == "." || safeFilename == "/" || safeFilename == "" {
		c.String(http.StatusBadRequest, "Invalid filename.")
		return
	}

	filePath := filepath.Join(uploadDir, safeFilename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "Media not found.")
		return
	}
	c.File(filePath)
}
