package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"ms_template/internal/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type App struct {
	log         *slog.Logger
	cfg         *config.Config
	grpcServer  *grpc.Server
	httpServer  *http.Server
	port        int
	metricsPort int
}

func New(log *slog.Logger, cfg *config.Config) *App {
	return &App{
		log:         log,
		cfg:         cfg,
		grpcServer:  grpc.NewServer(),
		port:        *cfg.GRPC.Port,
		metricsPort: *cfg.Prometheus.Port,
	}
}

func (a *App) Run() error {
	// Запуск gRPC сервера
	go func() {
		if err := a.runGRPCServer(); err != nil {
			a.log.Error("Ошибка gRPC сервера", "error", err)
		}
	}()

	// Запуск HTTP сервера для метрик
	go func() {
		if err := a.runMetricsServer(); err != nil && err != http.ErrServerClosed {
			a.log.Error("Ошибка метрик сервера", "error", err)
		}
	}()

	// Блокируем основную горутину
	select {}
}

func (a *App) runGRPCServer() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("ошибка прослушивания порта %d: %w", a.port, err)
	}

	a.log.Info("gRPC server started",
		slog.String("addr", l.Addr().String()),
		slog.Int("metrics_port", a.metricsPort),
	)

	if err := a.grpcServer.Serve(l); err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("ошибка обслуживания grpc сервера: %w", err)
	}

	return nil
}

func (a *App) runMetricsServer() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.metricsPort),
		Handler: mux,
	}

	a.log.Info("Metrics server started", "port", a.metricsPort)
	return a.httpServer.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("Начало graceful shutdown...")
	
	// Останавливаем gRPC сервер
	if a.grpcServer != nil {
		a.log.Info("Остановка gRPC сервера...")
		a.grpcServer.GracefulStop()
		a.log.Info("gRPC сервер остановлен")
	}
	
	// Останавливаем HTTP сервер
	if a.httpServer != nil {
		a.log.Info("Остановка HTTP сервера метрик...")
		if err := a.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("ошибка остановки HTTP сервера: %w", err)
		}
		a.log.Info("HTTP сервер метрик остановлен")
	}
	
	// Здесь можно добавить закрытие других ресурсов
	// (базы данных, пулов соединений, кэшей и т.д.)
	
	return nil
}

// GRPCServer возвращает gRPC сервер (для совместимости со старым кодом)
func (a *App) GRPCServer() *grpc.Server {
	return a.grpcServer
}

// RunMetricsServer запускает сервер метрик (публичный метод)
func (a *App) RunMetricsServer() error {
	return a.runMetricsServer()
}