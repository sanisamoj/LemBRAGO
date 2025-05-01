package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/cache"
	"lembrago.com/lembrago/models"
)

func DictionaryPreviewMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AuthCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}

		actAtt, _ := cache.GetInt(fmt.Sprintf("att-%s", req.Email))
		fmt.Println(actAtt)
		if actAtt >= 5 {
			c.AbortWithStatusJSON(429, gin.H{"error": "Too many attempts"})
			return
		}

		go attackDicPrev(req.Email)
		c.Set("validatedAuthCodeRequest", req)

		c.Next()
	}
}

func attackDicPrev(email string) {
	attKey := fmt.Sprintf("att-%s", email)
	cache.Increment(attKey)
	cache.SetTTL(attKey, 30 * time.Minute)
}


