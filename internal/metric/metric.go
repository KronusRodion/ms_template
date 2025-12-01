package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Metrics собирает основные метрики gRPC приложения
type Metrics struct {
	// gRPC метрики
	grpcRequestsTotal    *prometheus.CounterVec
	grpcRequestDuration  *prometheus.HistogramVec
	grpcRequestsInFlight *prometheus.GaugeVec

	// Бизнес-метрики
	activeConnections prometheus.Gauge
	errorsTotal       *prometheus.CounterVec

	// Регистр
	registry *prometheus.Registry
}

// New создает и регистрирует метрики
func New(appName string) *Metrics {
	registry := prometheus.NewRegistry()

	// Регистрируем стандартные коллекторы
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	m := &Metrics{
		registry: registry,
	}

	m.initializeGRPCMetrics(appName)
	return m
}

// initializeGRPCMetrics инициализирует gRPC метрики
func (m *Metrics) initializeGRPCMetrics(appName string) {
	constLabels := prometheus.Labels{"app": appName}

	// Общее количество gRPC запросов
	m.grpcRequestsTotal = promauto.With(m.registry).NewCounterVec(
		prometheus.CounterOpts{
			Name:        "grpc_requests_total",
			Help:        "Total number of gRPC requests",
			ConstLabels: constLabels,
		},
		[]string{"method", "code"},
	)

	// Длительность gRPC запросов
	m.grpcRequestDuration = promauto.With(m.registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "grpc_request_duration_seconds",
			Help:        "gRPC request duration in seconds",
			Buckets:     []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			ConstLabels: constLabels,
		},
		[]string{"method", "code"},
	)

	// Текущие выполняющиеся gRPC запросы
	m.grpcRequestsInFlight = promauto.With(m.registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "grpc_requests_in_flight",
			Help:        "Current number of gRPC requests being served",
			ConstLabels: constLabels,
		},
		[]string{"method"},
	)

	// Активные соединения
	m.activeConnections = promauto.With(m.registry).NewGauge(
		prometheus.GaugeOpts{
			Name:        "grpc_active_connections",
			Help:        "Current number of active gRPC connections",
			ConstLabels: constLabels,
		},
	)

	// Ошибки
	m.errorsTotal = promauto.With(m.registry).NewCounterVec(
		prometheus.CounterOpts{
			Name:        "grpc_errors_total",
			Help:        "Total number of gRPC errors",
			ConstLabels: constLabels,
		},
		[]string{"method", "type"},
	)
}

// UnaryServerInterceptor возвращает interceptor для gRPC метрик
func (m *Metrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		method := info.FullMethod

		// Увеличиваем счетчик текущих запросов
		m.grpcRequestsInFlight.WithLabelValues(method).Inc()
		defer m.grpcRequestsInFlight.WithLabelValues(method).Dec()

		// Обрабатываем запрос
		resp, err := handler(ctx, req)

		// Определяем код ответа
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
			
			// Регистрируем ошибку
			errorType := "business"
			if code == codes.Internal || code == codes.Unavailable {
				errorType = "internal"
			}
			m.errorsTotal.WithLabelValues(method, errorType).Inc()
		}

		// Регистрируем метрики
		duration := time.Since(start).Seconds()
		codeStr := code.String()

		m.grpcRequestsTotal.WithLabelValues(method, codeStr).Inc()
		m.grpcRequestDuration.WithLabelValues(method, codeStr).Observe(duration)

		return resp, err
	}
}

// StreamServerInterceptor возвращает stream interceptor для gRPC метрик
func (m *Metrics) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()
		method := info.FullMethod

		// Увеличиваем счетчик текущих запросов
		m.grpcRequestsInFlight.WithLabelValues(method).Inc()
		defer m.grpcRequestsInFlight.WithLabelValues(method).Dec()

		// Обрабатываем запрос
		err := handler(srv, ss)

		// Определяем код ответа
		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
			
			// Регистрируем ошибку
			errorType := "business"
			if code == codes.Internal || code == codes.Unavailable {
				errorType = "internal"
			}
			m.errorsTotal.WithLabelValues(method, errorType).Inc()
		}

		// Регистрируем метрики
		duration := time.Since(start).Seconds()
		codeStr := code.String()

		m.grpcRequestsTotal.WithLabelValues(method, codeStr).Inc()
		m.grpcRequestDuration.WithLabelValues(method, codeStr).Observe(duration)

		return err
	}
}

// ConnectionOpened увеличивает счетчик активных соединений
func (m *Metrics) ConnectionOpened() {
	m.activeConnections.Inc()
}

// ConnectionClosed уменьшает счетчик активных соединений
func (m *Metrics) ConnectionClosed() {
	m.activeConnections.Dec()
}

// ReportError регистрирует ошибку
func (m *Metrics) ReportError(method, errorType string) {
	m.errorsTotal.WithLabelValues(method, errorType).Inc()
}

// Handler возвращает http.Handler для метрик Prometheus
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// GetRegistry возвращает регистр метрик (для кастомных метрик)
func (m *Metrics) GetRegistry() *prometheus.Registry {
	return m.registry
}
