package main

import (
	"log"
	"ms_template/internal/app"
	"ms_template/internal/config"
	"ms_template/internal/logger"
	"os"
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

	app.GRPCServer.Run()
}
