package db

import (
	"context"
	"fmt"

	"github.com/IPampurin/UrlShortener/pkg/configuration"
	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
	"github.com/wb-go/wbf/logger"
)

// ClientPostgres хранит подключение к БД
// делаем его публичным, чтобы другие пакеты могли использовать методы
type ClientPostgres struct {
	*pgxdriver.Postgres
}

// глобальный экземпляр клиента (синглтон)
var clientPostgres *ClientPostgres

// InitDB инициализирует подключение к PostgreSQL и применяет миграции
func InitDB(ctx context.Context, cfgDb *configuration.ConfDB, log logger.Logger) error {

	// формируем DSN из конфигурации
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfgDb.User, cfgDb.Password, cfgDb.HostName, cfgDb.Port, cfgDb.Name)

	// создаём клиент pgxdriver с параметрами по умолчанию
	pgxConn, err := pgxdriver.New(dsn, log)
	if err != nil {
		return fmt.Errorf("ошибка создания клиента pgxdriver: %w", err)
	}
	// defer pg.Close() // закрываем из main()

	// проверяем соединение
	if err = pgxConn.Ping(ctx); err != nil {
		return fmt.Errorf("ошибка соединения с клиентом pgxdriver: %w", err)
	}

	clientPostgres = &ClientPostgres{pgxConn}

	// запускаем миграции
	if err = clientPostgres.Migration(ctx); err != nil {
		return fmt.Errorf("ошибка миграций: %w", err)
	}

	log.Info("база данных PostgreSQL успешно запущена, миграции применены.")

	return nil
}

// CloseDB закрывает пул соединений с БД
func CloseDB() error {

	if clientPostgres != nil {
		clientPostgres.Close()
		clientPostgres = nil
	}

	return nil
}

// GetClientDB возвращает глобальный экземпляр клиента БД
func GetClientDB() *ClientPostgres {

	return clientPostgres
}
