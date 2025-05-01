package middlewares

import (
	"time"

	"github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"lembrago.com/lembrago/cache"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

var store ratelimit.Store

func NewRateLimiterMiddleware(rate time.Duration, limit uint) gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: cache.RedisClient,
		Rate:        rate,
		Limit:       limit,
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}
