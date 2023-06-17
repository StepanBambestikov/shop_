package middleware

import (
	"net/http"
	"orderServiceGit/internal/core/entities"
	"orderServiceGit/internal/log"

	"github.com/gin-gonic/gin"
)

var ErrReplyUnknown = entities.ReplyError("Unknown error", 500)

func GinError() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if errors := c.Errors.ByType(gin.ErrorTypeAny); len(errors) > 0 {
			err := errors[0].Err
			if err, ok := err.(*entities.Error); ok {
				log.Error("[Open error]:", err)
				c.AbortWithStatusJSON(err.Code, err.ToReply())
				return
			}
			log.Error("[Hidden error]:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrReplyUnknown)
			return
		}
	}
}
