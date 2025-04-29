package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/errors"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		fmt.Println("ErroHandler: ", c.Errors)
		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

			if appErr, ok := err.(*errors.AppError); ok {
				log.Printf("[AppError] Ocorreu um erro: Code=%d, Message=%s\n", appErr.Code, appErr.Message)
				c.AbortWithStatusJSON(appErr.Code, gin.H{
					"message": appErr.Message,
				})
				return
			}

			log.Printf("[InternalError] Ocorreu um erro inesperado: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Ocorreu um erro interno inesperado no servidor.",
			})
			return
		}
	}
}
