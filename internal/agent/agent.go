package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/hash"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/avast/retry-go"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Agent struct {
	URL      string
	logger   *zap.SugaredLogger
	Compress bool
}

func NewAgent(url string, compress bool) *Agent {
	return &Agent{
		URL:      url,
		logger:   logger.NewLogger(),
		Compress: compress,
	}
}

func (a *Agent) SendRequest() error {
	for k, v := range metrics.GaugeAgent.ReturnValues() {
		requestURL := fmt.Sprintf("%sgauge/%s/%f", a.URL, k, v)
		res, err := http.Post(requestURL, "Content-Type: text/plain", nil)
		if err != nil {
			fmt.Println(err)
			return errors.New("request failed")
		}
		if er := res.Body.Close(); er != nil {
			return er
		}
	}
	requestURL := fmt.Sprintf("%scounter/PollCount/%d", a.URL, metrics.PollCount)
	res, err := http.Post(requestURL, "Content-Type: text/plain", nil)
	if err != nil {
		return errors.New("request failed")
	}
	if er := res.Body.Close(); er != nil {
		return er
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (a *Agent) compressData(requestBody []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	compressor := gzip.NewWriter(buf)

	if _, err := compressor.Write(requestBody); err != nil {
		a.logger.Error("Error while compress data:", err.Error())
		return nil, err
	}
	if err := compressor.Close(); err != nil {
		logger.Logger.Error("Error while compress data:", err.Error())
		return nil, err
	}
	return buf, nil
}

func (a *Agent) sendRequestGaugeJSON() error {

	req := metrics.Metrics{MType: "gauge"}
	client := http.Client{Timeout: time.Duration(60) * time.Second}

	for k, v := range metrics.GaugeAgent.ReturnValues() {
		a.logger.Infow("Encode gauge metric", "addr", config.FlagRunAgAddr, "name", k, "value", v)
		req.ID = k
		req.Value = &v
		reqByte, err := json.Marshal(req)
		data := bytes.NewBuffer(reqByte)

		if err != nil {
			a.logger.Error(err.Error(), "event", "encode data")
			return errors.New("encoding data failed")
		}

		if a.Compress {
			var errCompress error
			data, errCompress = a.compressData(reqByte)

			if errCompress != nil {
				a.logger.Error(errCompress.Error(), "event", "compress data")
				return errors.New("compress data failed")
			}
		}
		a.logger.Infow("Send gauge metric", "addr", config.FlagRunAgAddr, "data", data)
		attempt := 0
		err = retry.Do(func() error {
			var er error
			r, _ := http.NewRequest("POST", a.URL, data)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Accept-Encoding", "gzip")
			res, er := client.Do(r)

			if res != nil {
				if erClose := res.Body.Close(); erClose != nil {
					a.logger.Error(erClose.Error(), "event", "close response")
					return erClose
				}
			}
			return er
		},
			retry.Attempts(3),
			retry.OnRetry(func(n uint, err error) {
				time.Sleep(time.Duration(1 + 2*attempt))
				attempt += 1
				a.logger.Info("Retrying request after error: %v", err)
			}))

		if err != nil {
			a.logger.Error(err.Error(), "event", "send request")
			return errors.New("request failed")
		}

	}
	return nil
}

func (a *Agent) sendRequestCounterJSON() error {

	req := metrics.Metrics{MType: "counter", ID: "PollCount", Delta: &metrics.PollCount}
	client := http.Client{Timeout: time.Duration(60) * time.Second}

	a.logger.Infow("Encode counter metric", "addr", config.FlagRunAgAddr, "name", "PollCount", "value", metrics.PollCount)

	reqByte, err := json.Marshal(req)
	data := bytes.NewBuffer(reqByte)

	if err != nil {
		a.logger.Error(err.Error(), "event", "encode data")
		return errors.New("encoding data failed")
	}

	if a.Compress {
		var errCompress error
		data, errCompress = a.compressData(reqByte)

		if errCompress != nil {
			a.logger.Error(errCompress.Error(), "event", "compress data")
			return errors.New("compress data failed")
		}
	}

	a.logger.Infow("Send counter metric", "addr", config.FlagRunAgAddr, "data", data)
	attempt := 0
	err = retry.Do(func() error {
		r, _ := http.NewRequest("POST", a.URL, data)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "gzip")
		res, er := client.Do(r)

		if res != nil {
			if erClose := res.Body.Close(); erClose != nil {
				a.logger.Error(erClose.Error(), "event", "close response")
				return erClose
			}
		}
		return er
	},
		retry.Attempts(3),
		retry.OnRetry(func(n uint, err error) {
			time.Sleep(time.Duration(1 + 2*attempt))
			attempt += 1
			a.logger.Info("Retrying request after error: %v", err)

		}))

	if err != nil {
		a.logger.Error(err.Error(), "event", "send request")
		return errors.New("request failed")
	}
	return nil
}

func (a *Agent) SendRequestJSON() error {

	a.logger.Infow("Start send metrics")
	if err := a.sendRequestGaugeJSON(); err != nil {
		a.logger.Error("Error in function sendRequestGaugeJSON")
		return err
	}
	if err := a.sendRequestCounterJSON(); err != nil {
		a.logger.Error("Error in function sendRequestCounterJSON")
		return err
	}

	return nil
}

func (a *Agent) SendBatchDataJSON() error {

	a.logger.Infow("Start send metrics")

	var req []metrics.Metrics
	client := http.Client{Timeout: time.Duration(60) * time.Second}

	for k, v := range metrics.GaugeAgent.ReturnValues() {
		a.logger.Infow("Add gauge metric to list", "name", k, "value", v)
		value := new(float64)
		*value = v
		req = append(req, metrics.Metrics{MType: "gauge", ID: k, Value: value})
	}
	a.logger.Infow("Add counter metric to list", "name", "PollCount", "value", metrics.PollCount)
	req = append(req, metrics.Metrics{MType: "counter", ID: "PollCount", Delta: &metrics.PollCount})

	a.logger.Infow("Encode all metrics")

	reqByte, err := json.Marshal(req)

	if err != nil {
		a.logger.Error(err.Error(), "event", "encode data")
		return errors.New("encoding data failed")
	}

	data := bytes.NewBuffer(reqByte)

	if a.Compress {
		var errCompress error
		data, errCompress = a.compressData(reqByte)

		if errCompress != nil {
			a.logger.Error(errCompress.Error(), "event", "compress data")
			return errors.New("compress data failed")
		}
	}
	attempt := 0
	a.logger.Infow("Send all metrics", "addr", config.FlagRunAgAddr, "data", data)
	err = retry.Do(func() error {
		r, _ := http.NewRequest("POST", a.URL, data)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "gzip")
		if config.HashKeyAgent != "" {
			sign := hash.GetSignature(reqByte, config.HashKeyAgent)
			r.Header.Set("HashSHA256", base64.StdEncoding.EncodeToString(sign))
		}
		res, er := client.Do(r)

		if res != nil {
			if erClose := res.Body.Close(); erClose != nil {
				a.logger.Error(erClose.Error(), "event", "close response")
				return erClose
			}
		}
		return er
	},
		retry.Attempts(3),
		retry.OnRetry(func(n uint, err error) {
			time.Sleep(time.Duration(1 + 2*attempt))
			attempt += 1
			a.logger.Info("Retrying request after error: %v", err)
		}))

	if err != nil {
		a.logger.Error(err.Error(), "event", "send request")
		return errors.New("request failed")
	}

	return nil
}
