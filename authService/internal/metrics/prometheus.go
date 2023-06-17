package metrics

import (
	"fmt"
	"gitea.teneshag.ru/gigabit/goauth/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusMetricsRepository struct {
	metricsConfig *core.MetricsConfig
	registry      *prometheus.Registry

	counters map[string]*prometheus.CounterVec
}

func (p *prometheusMetricsRepository) GetHandlerFunc() gin.HandlerFunc {
	return gin.WrapH(promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))
}

func (p *prometheusMetricsRepository) AddCounter(name string, labelNames []string) error {
	if _, ok := p.counters[name]; ok {
		return fmt.Errorf("AddCounter: metrics %s already exists", name)
	}

	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{Name: name}, labelNames)
	p.counters[name] = counterVec
	return p.registry.Register(counterVec)
}

func (p *prometheusMetricsRepository) IncCounter(name string, labels map[string]string) error {
	metric, ok := p.counters[name]
	if !ok {
		return fmt.Errorf("IncCounter: no such metric: %s", name)
	}

	metric.With(labels).Inc()

	return nil
}

func (p *prometheusMetricsRepository) AddGauge(name string) error {
	//TODO implement me
	panic("implement me")
}

func (p *prometheusMetricsRepository) SetGauge(name string, value float64, labels map[string]string) error {
	//TODO implement me
	panic("implement me")
}

func NewPrometheusRepository(metricsConfig *core.MetricsConfig) (Repository, error) {
	registry := prometheus.NewRegistry()
	// Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
	)

	return &prometheusMetricsRepository{
		metricsConfig: metricsConfig,
		registry:      registry,
		counters:      make(map[string]*prometheus.CounterVec),
	}, nil
}
