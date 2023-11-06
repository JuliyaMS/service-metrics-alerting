package database

import (
	"context"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/jackc/pgx/v5"
	"time"
)

type ConnectionDB struct {
	Conn *pgx.Conn
}

func NewConnectionDB() *ConnectionDB {
	logger.Logger.Info("create context with timeout")

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	logger.Logger.Info("connect to Database")
	conn, err := pgx.Connect(ctx, config.DatabaseDsn)
	if err != nil {
		logger.Logger.Error("get error while connection to database:", err)
		return nil
	}

	logger.Logger.Info("success connection")
	return &ConnectionDB{
		Conn: conn,
	}
}

func (db *ConnectionDB) Close() error {
	err := db.Conn.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}
