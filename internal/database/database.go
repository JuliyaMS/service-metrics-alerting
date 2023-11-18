package database

import (
	"context"
	"errors"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
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

func (db *ConnectionDB) Init() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	_, errEx := db.Conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS gauge_metrics(Name varchar(100) PRIMARY KEY, Value double precision NOT NULL);"+
		"CREATE TABLE IF NOT EXISTS count_metrics(Name varchar(100) PRIMARY KEY, Value bigint NOT NULL);")

	if errEx != nil {
		logger.Logger.Info("Error while create tables:", errEx.Error())
		return
	}
}

func (db *ConnectionDB) Add(t, name, val string) error {

	logger.Logger.Info("Add value to DB")

	var sql string

	if !metrics.CheckType(t) {
		return errors.New("this type of metric doesn't exists")
	}
	logger.Logger.Info("Create sql string")
	if t == "counter" {
		sql = "INSERT INTO count_metrics (Name, Value) VALUES ($1, $2)" +
			"ON CONFLICT (Name) DO UPDATE SET Value = count_metrics.Value + $2;"
	} else {
		sql = "INSERT INTO gauge_metrics (Name, Value) VALUES ($1, $2)" +
			"ON CONFLICT (Name) DO UPDATE SET Value = $2;"
	}

	logger.Logger.Info("Create sql context")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	logger.Logger.Info("Execute request to add data")
	_, errEx := db.Conn.Exec(ctx, sql, name, val)

	if errEx != nil {
		logger.Logger.Info("get error while add values:", errEx.Error())
		return errEx
	}
	logger.Logger.Info("Execute successful")
	return nil
}

func (db *ConnectionDB) Get(tp, name string) string {

	logger.Logger.Info("Get value from DB")

	var sql string
	var value string

	if metrics.CheckType(tp) {
		logger.Logger.Info("Create sql string")

		if tp == "gauge" {
			sql = "SELECT Value FROM gauge_metrics WHERE Name=$1"
		} else {
			sql = "SELECT Value FROM count_metrics WHERE Name=$1"
		}

		logger.Logger.Info("Create sql context")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
		defer cancel()

		logger.Logger.Info("Execute request to get data")
		if err := db.Conn.QueryRow(ctx, sql, name).Scan(&value); err != nil {
			logger.Logger.Error("Get error while execute request: ", err.Error())
			return "-1"
		}
		logger.Logger.Info("Execute successful")
		return value

	}
	logger.Logger.Error("Metric`s type is not correct")
	return "-1"
}

func (db *ConnectionDB) getAllGaugeMetrics() metrics.GaugeMetrics {

	logger.Logger.Info("Get all gauge values from DB")
	var gauge metrics.GaugeMetrics
	gauge.Metrics = make(map[string]float64)

	logger.Logger.Info("Create sql context")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	logger.Logger.Info("Execute request to get all gauge metrics")
	rows, err := db.Conn.Query(ctx, "SELECT * FROM gauge_metrics;")
	if err != nil {
		logger.Logger.Error("Get error while execute request: ", err.Error())
		return metrics.GaugeMetrics{}
	}

	defer rows.Close()

	logger.Logger.Info("Scan data from rows")
	for rows.Next() {
		var (
			name  string
			value float64
		)

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Logger.Error("Get error while scan row: ", err.Error())
			return metrics.GaugeMetrics{}
		}
		gauge.Metrics[name] = value
	}
	logger.Logger.Info("Execute successful")
	return gauge
}

func (db *ConnectionDB) getAllCountMetrics() metrics.CounterMetrics {

	logger.Logger.Info("Get all count values from DB")
	var counter metrics.CounterMetrics
	counter.Metrics = make(map[string]int64)

	logger.Logger.Info("Create sql context")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()

	logger.Logger.Info("Execute request to get all count metrics")
	rows, err := db.Conn.Query(ctx, "SELECT * FROM count_metrics;")
	if err != nil {
		logger.Logger.Error("Get error while execute request: ", err.Error())
		return metrics.CounterMetrics{}
	}

	defer rows.Close()

	logger.Logger.Info("Scan data from rows")
	for rows.Next() {
		var (
			name  string
			value int64
		)

		err = rows.Scan(&name, &value)
		if err != nil {
			logger.Logger.Error("Get error while scan row: ", err.Error())
			return metrics.CounterMetrics{}
		}
		counter.Metrics[name] = value
	}
	logger.Logger.Info("Execute successful")
	return counter
}

func (db *ConnectionDB) GetAll() (metrics.GaugeMetrics, metrics.CounterMetrics) {
	return db.getAllGaugeMetrics(), db.getAllCountMetrics()
}

func (db *ConnectionDB) Close() error {
	err := db.Conn.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}
