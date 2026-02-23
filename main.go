package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IPampurin/UrlShortener/pkg/configuration"
	"github.com/IPampurin/UrlShortener/pkg/server"
	"github.com/wb-go/wbf/logger"
)

func main() {

	// cоздаём контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// запускаем горутину обработки сигналов
	go func() {
		<-sigChan
		cancel()
	}()

	var err error

	// считываем .env файл
	cfg, err := configuration.ReadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// настраиваем логгер
	appLogger, err := logger.InitLogger(
		logger.ZapEngine,
		"UrlShortener",
		os.Getenv("APP_ENV"), // пока оставим пустым
		logger.WithLevel(logger.InfoLevel),
	)
	if err != nil {
		log.Fatalf("Ошибка создания логгера: %v", err)
	}
	defer func() { _ = appLogger.(*logger.ZapAdapter) }()
	/*
		// подключаем базу данных
		err = db.InitDB(&cfg.DB)
		if err != nil {
			appLogger.Error("ошибка подключения к БД", "error", err)
			return
		}
		defer func() { _ = db.CloseDB() }()

		// инициализируем кэш
		err = cache.InitRedis(&cfg.Redis)
		if err != nil {
			appLogger.Warn("кэш не работает", "error", err)
		}
	*/
	// запускаем сервер
	err = server.Run(ctx, &cfg.Server, appLogger)
	if err != nil {
		appLogger.Error("Ошибка сервера", "error", err)
		cancel()
		return
	}

	appLogger.Info("Приложение корректно завершено")
}
