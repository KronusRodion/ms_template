package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ms_template/internal/app"
	"ms_template/internal/config"
	"ms_template/internal/logger"
	"golang.org/x/sync/errgroup"
)

func main() {

	path, exist := os.LookupEnv("CONF_PATH")
	if !exist {
		log.Fatalf("Не задана обязательная переменная CONF_PATH")
	}

	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatalf("Не удалось загрузить конфиг: %v", err)
	}

	logger := logger.Setup(cfg.Env)

	app := app.New(logger, cfg)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для получения сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Группа для управления горутинами
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return app.Run()
	})

	g.Go(func() error {
		return app.RunMetricsServer()
	})

	g.Go(func() error {
		select {
		case sig := <-sigChan:
			logger.Info("Получен сигнал завершения", "signal", sig.String())
			cancel()
			
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer shutdownCancel()
			
			if err := app.Shutdown(shutdownCtx); err != nil {
				logger.Error("Ошибка при graceful shutdown", "error", err)
				return err
			}
			
			logger.Info("Приложение корректно завершено")
			return nil
		case <-gCtx.Done():
			return gCtx.Err()
		}
	})

	// Ожидание завершения всех горутин
	if err := g.Wait(); err != nil && err != context.Canceled {
		logger.Error("Ошибка в работе приложения", "error", err)
		os.Exit(1)
	}
}