package grpcserver

import (
	"fmt"
	"log/slog"
	"ms_template/internal/grpc/notesGRPC"
	metrics "ms_template/internal/metric"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

type App struct {
    log        *slog.Logger
    gRPCServer *grpc.Server
    metrics      *metrics.Metrics // Добавляем метрики
    port         int              // gRPC порт
    metricsPort  int
}


func New(log *slog.Logger, NoteServer notesGRPC.NoteServer, port int, metricsPort int) *App {
	metrics := metrics.New("notes_service")
    
    // Настраиваем gRPC сервер с interceptors для метрик
    gRPCServer := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            recovery.UnaryServerInterceptor(),
            metrics.UnaryServerInterceptor(), // Добавляем метрики interceptor
        ),
        grpc.ChainStreamInterceptor(
            metrics.StreamServerInterceptor(), // Для stream соединений
        ),
    )

	notesGRPC.Register(gRPCServer, NoteServer)
	
	return &App{
        log:         log,
        gRPCServer:  gRPCServer,
        metrics:     metrics,
        port:        port,
        metricsPort: metricsPort,
    }
}


func (a *App) Run() error {

    l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
    if err != nil {
        return fmt.Errorf("ошибка прослушивания порта %d: %w", a.port, err)
    }

    a.log.Info("gRPC server started", 
        slog.String("addr", l.Addr().String()),
        slog.Int("metrics_port", a.metricsPort),
    )

    // Запускаем обработчик gRPC-сообщений
    if err := a.gRPCServer.Serve(l); err != nil {
        return fmt.Errorf("ошибка обслуживания grpc сервера: %w", err)
    }

    return nil
}



// MustRun запускает приложение и паникует при ошибке
func (a *App) MustRun() {
    if err := a.Run(); err != nil {
        panic(err)
    }
}

// GracefulStop останавливает сервер
func (a *App) GracefulStop() {
    a.log.Info("Shutting down gRPC server...")
    a.gRPCServer.GracefulStop()
    a.log.Info("gRPC server stopped")
}