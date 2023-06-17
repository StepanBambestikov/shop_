package middleware

import (
	"catalogServiceGit/internal/log"
	"catalogServiceGit/internal/metrics"
	"github.com/gin-gonic/gin"
	"strconv"
)

const counterName = "gaimp_http_response_codes"

func Metrics(repository metrics.Repository) gin.HandlerFunc {
	err := repository.AddCounter(counterName, []string{"code", "path"})
	if err != nil {
		log.Fatal("Can't add metric: ", err)
	}
	return func(context *gin.Context) {
		context.Next()
		code := context.Writer.Status()
		path := context.FullPath()

		err := repository.IncCounter(counterName, map[string]string{
			"code": strconv.Itoa(code),
			"path": path,
		})
		if err != nil {
			log.Error("Metrics middleware can't inc metric "+counterName, err.Error())
		}
	}
}
