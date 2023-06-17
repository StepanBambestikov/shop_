package metrics

import (
	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetHandlerFunc() gin.HandlerFunc

	AddCounter(name string, labelNames []string) error
	IncCounter(name string, labels map[string]string) error

	AddGauge(name string) error
	SetGauge(name string, value float64, labels map[string]string) error
}
