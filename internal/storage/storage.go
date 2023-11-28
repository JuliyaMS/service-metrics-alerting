package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
	"strconv"
	"sync"
	"time"
)

var fileMutex sync.Mutex
var DBMutex sync.Mutex

type Repositories interface {
	Init()
	Add(t, name, val string) error
	Get(tp, name string) string
	GetAll() (metrics.GaugeMetrics, metrics.CounterMetrics)
	CheckConnection() error
	AddAnyData(req []metrics.Metrics) error
	Close() error
}

func NewStorage() Repositories {
	logger.Logger.Info("Create new storage")

	if config.DatabaseDsn != "" {
		return NewConnectionDB()
	}

	if config.Restore && config.FileStoragePath != "" {
		logger.Logger.Info("restore data from file:", config.FileStoragePath)

		storage, err := ReadFromFile(config.FileStoragePath)
		if err != nil {
			logger.Logger.Errorf(err.Error(), "can't read data from file:", config.FileStoragePath)
		}
		return storage
	}

	return new(MemStorage)
}

type MemStorage struct {
	MetricsGauge   metrics.GaugeMetrics
	MetricsCounter metrics.CounterMetrics
}

func (s *MemStorage) Init() {
	s.MetricsGauge.Init()
	s.MetricsCounter.Init()
}

func (s MemStorage) Add(t, name, val string) error {
	if !metrics.CheckType(t) {
		return errors.New("this type of metric doesn't exists")
	}
	if t == "counter" {
		if s.MetricsCounter.Add(name, val) {
			return nil
		}
	}
	if s.MetricsGauge.Add(name, val) {
		return nil
	}

	return errors.New("can't add metric")
}

func (s MemStorage) Get(tp, name string) string {
	if metrics.CheckType(tp) {
		if tp == "gauge" {
			value := s.MetricsGauge.Get(name)
			return value
		}
		value := s.MetricsCounter.Get(name)
		return value

	}
	return "-1"
}

func (s *MemStorage) GetAll() (metrics.GaugeMetrics, metrics.CounterMetrics) {
	return s.MetricsGauge, s.MetricsCounter
}

func (s *MemStorage) CheckConnection() error {
	if s.MetricsGauge.Metrics != nil && s.MetricsCounter.Metrics != nil {
		return nil
	}
	return errors.New("storage for metrics isn`t initialize")
}

func (s *MemStorage) AddAnyData(req []metrics.Metrics) error {

	for _, r := range req {
		if r.MType == "gauge" {
			value := strconv.FormatFloat(*r.Value, 'g', -1, 64)
			if err := s.Add(r.MType, r.ID, value); err != nil {
				return err
			}
		}
		if r.MType == "counter" {
			value := strconv.FormatInt(*r.Delta, 10)
			if err := s.Add(r.MType, r.ID, value); err != nil {
				return err
			}
		}

	}
	return nil
}

func (s MemStorage) Close() error {
	s.MetricsGauge.Close()
	s.MetricsCounter.Close()
	return nil
}

func WriteToFile(fileName string, stor *Repositories) error {

	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(stor); err != nil {
		return err
	}
	file.Close()
	return nil
}

func ReadFromFile(fileName string) (*MemStorage, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	var stor MemStorage
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return new(MemStorage), err
	}
	encoder := json.NewDecoder(file)
	if err = encoder.Decode(&stor); err != nil {
		return new(MemStorage), err
	}
	file.Close()
	return &stor, nil
}

type ConnectionDB struct {
	Conn *pgx.Conn
}

func NewConnectionDB() *ConnectionDB {
	if config.DatabaseDsn == "" {
		return nil
	}

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

func (db *ConnectionDB) CheckConnection() error {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	logger.Logger.Info("check connection to Database")
	err := Retry(4, time.Duration(1), db.Conn.Ping, ctx, "", "", 0)
	if err != nil {
		return err
	}
	return nil
}

func (db *ConnectionDB) Init() {

	logger.Logger.Info("Start creation tables for metrics")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
	defer cancel()
	sql := "CREATE TABLE IF NOT EXISTS gauge_metrics(Name varchar(100) PRIMARY KEY, Value double precision NOT NULL);" +
		"CREATE TABLE IF NOT EXISTS count_metrics(Name varchar(100) PRIMARY KEY, Value bigint NOT NULL);"

	_, errEx := db.Conn.Exec(ctx, sql)

	if errEx != nil {
		logger.Logger.Info("Error while create table gauge_metrics: ", errEx.Error())
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
	errEx := Retry(4, time.Duration(1), db.Conn.Exec, ctx, sql, name, val)

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

		logger.Logger.Info("Execute request to get data")
		if err := RetryQueryRow(4, time.Duration(1), db.Conn, sql, name, &value); err != nil {
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

	logger.Logger.Info("Execute request to get all gauge metrics")
	rows, err := RetryQuery(4, time.Duration(1), db.Conn, "SELECT * FROM gauge_metrics;")
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

	logger.Logger.Info("Execute request to get all count metrics")
	rows, err := RetryQuery(4, time.Duration(1), db.Conn, "SELECT * FROM count_metrics;")
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

func (db *ConnectionDB) AddAnyData(req []metrics.Metrics) error {
	DBMutex.Lock()
	defer DBMutex.Unlock()

	logger.Logger.Infow("Start transaction")
	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		return err
	}
	for _, el := range req {
		if el.MType == "counter" {

			sql := "INSERT INTO count_metrics (Name, Value) VALUES ($1, $2)" +
				"ON CONFLICT (Name) DO UPDATE SET Value = count_metrics.Value + $2;"

			logger.Logger.Info("Execute request to add counter metric")
			err = Retry(4, time.Duration(1), tx.Exec, context.Background(), sql, el.ID, el.Delta)
		} else {
			sql := "INSERT INTO gauge_metrics (Name, Value) VALUES ($1, $2)" +
				"ON CONFLICT (Name) DO UPDATE SET Value = $2;"

			logger.Logger.Info("Execute request to add gauge metric")
			err = Retry(4, time.Duration(1), tx.Exec, context.Background(), sql, el.ID, el.Value)
		}

		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
		logger.Logger.Info("Execute successful")

	}

	logger.Logger.Info("Close transaction")
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	logger.Logger.Info("All data added successful")
	time.Sleep(1 * time.Second)
	return nil
}

func (db *ConnectionDB) Close() error {
	logger.Logger.Info("Close connection")
	err := Retry(4, time.Duration(1), db.Conn.Close, context.Background(), "", "", "")
	if err != nil {
		return err
	}
	return nil
}

func Retry(attempts int, sleep time.Duration, f interface{}, ctx context.Context, sql string, val1 any, val2 any) (err error) {
	logger.Logger.Info("Start retry function")
	for i := 0; ; i++ {
		logger.Logger.Info("Execute function, attempt:", i+1)
		switch t := f.(type) {
		case func(context.Context) error:
			logger.Logger.Info("Function type:", t)
			err = f.(func(context.Context) error)(ctx)
		case func(context.Context, string, ...any) (pgconn.CommandTag, error):
			logger.Logger.Info("Function type:", t)
			_, err = f.(func(context.Context, string, ...any) (pgconn.CommandTag, error))(ctx, sql, val1, val2)
		}
		if err != nil {
			logger.Logger.Info("Check type of error")
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				logger.Logger.Info("Get retryable-error")
				if i >= (attempts - 1) {
					logger.Logger.Info("Number of attempts exhausted")
					break
				}
				logger.Logger.Info("Sleep...")
				time.Sleep((sleep + time.Duration(2*i)) * time.Second)
			} else {
				logger.Logger.Info("Get not retryable-error")
				return
			}
		} else {
			logger.Logger.Info("Function execute successfully")
			return
		}

	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func RetryQueryRow(attempts int, sleep time.Duration, Conn *pgx.Conn, sql string, name string, value *string) (err error) {
	logger.Logger.Info("Start retry function")
	for i := 0; ; i++ {
		logger.Logger.Info("Create sql context")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
		defer cancel()

		logger.Logger.Info("Execute function, attempt:", i+1)
		err = Conn.QueryRow(ctx, sql, name).Scan(value)
		if err != nil {
			logger.Logger.Info("Check type of error")
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				logger.Logger.Info("Get retryable-error")
				if i >= (attempts - 1) {
					logger.Logger.Info("Number of attempts exhausted")
					break
				}
				logger.Logger.Info("Sleep...")
				time.Sleep((sleep + time.Duration(2*i)) * time.Second)
			} else {
				logger.Logger.Info("Get not retryable-error")
				return
			}
		} else {
			logger.Logger.Info("Function execute successfully")
			return
		}

	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func RetryQuery(attempts int, sleep time.Duration, Conn *pgx.Conn, sql string) (rows pgx.Rows, err error) {
	logger.Logger.Info("Start retry function")
	for i := 0; ; i++ {
		logger.Logger.Info("Create sql context")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*200)
		defer cancel()

		logger.Logger.Info("Execute function, attempt:", i+1)
		rows, err = Conn.Query(ctx, sql)
		if err != nil {
			logger.Logger.Info("Check type of error")
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
				logger.Logger.Info("Get retryable-error")
				if i >= (attempts - 1) {
					logger.Logger.Info("Number of attempts exhausted")
					break
				}
				logger.Logger.Info("Sleep...")
				time.Sleep((sleep + time.Duration(2*i)) * time.Second)
			} else {
				logger.Logger.Info("Get not retryable-error")
				return
			}
		} else {
			logger.Logger.Info("Function execute successfully")
			return
		}

	}
	return rows, fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
