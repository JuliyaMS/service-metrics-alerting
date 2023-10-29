package requests

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/avast/retry-go"
	"net/http"
	"time"
)

func SendRequest() error {
	for k, v := range metrics.GaugeAgent.ReturnValues() {
		requestURL := fmt.Sprintf("http://%s/update/gauge/%s/%f", config.FlagRunAgAddr, k, v)
		res, err := http.Post(requestURL, "Content-Type: text/plain", nil)
		if err != nil {
			fmt.Println(err)
			return errors.New("request failed")
		}
		if er := res.Body.Close(); er != nil {
			return er
		}
	}
	requestURL := fmt.Sprintf("http://%s/update/counter/PollCount/%d", config.FlagRunAgAddr, metrics.PollCount)
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

func compressData(requestBody []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	if _, err := zb.Write(requestBody); err != nil {
		logger.Agent.Error("Error while compress data:", err.Error())
		return nil, err
	}
	if err := zb.Close(); err != nil {
		logger.Agent.Error("Error while compress data:", err.Error())
		return nil, err
	}
	return buf, nil
}

func sendRequestGaugeJSON(requestURL string) error {

	req := metrics.Metrics{MType: "gauge"}
	client := http.Client{Timeout: time.Duration(60) * time.Second}

	for k, v := range metrics.GaugeAgent.ReturnValues() {
		logger.Agent.Infow("Encode gauge metric", "addr", config.FlagRunAgAddr, "name", k, "value", v)
		req.ID = k
		req.Value = &v
		reqByte, err := json.Marshal(req)

		if err != nil {
			logger.Agent.Error(err.Error(), "event", "encode data")
			return errors.New("encoding data failed")
		}

		buf, errCompress := compressData(reqByte)

		if errCompress != nil {
			logger.Agent.Error(errCompress.Error(), "event", "compress data")
			return errors.New("compress data failed")
		}

		logger.Agent.Infow("Send gauge metric", "addr", config.FlagRunAgAddr, "data", buf)
		err = retry.Do(func() error {
			var er error
			r, _ := http.NewRequest("POST", requestURL, buf)
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Content-Encoding", "gzip")
			r.Header.Set("Accept-Encoding", "gzip")
			res, er := client.Do(r)

			if res != nil {
				if erClose := res.Body.Close(); erClose != nil {
					logger.Agent.Error(erClose.Error(), "event", "close response")
					return erClose
				}
			}
			return er
		},
			retry.Attempts(100),
			retry.OnRetry(func(n uint, err error) {
				logger.Agent.Info("Retrying request after error: %v", err)
			}))

		if err != nil {
			logger.Agent.Error(err.Error(), "event", "send request")
			return errors.New("request failed")
		}

	}
	return nil
}

func sendRequestCounterJSON(requestURL string) error {

	req := metrics.Metrics{MType: "counter", ID: "PollCount", Delta: &metrics.PollCount}
	client := http.Client{Timeout: time.Duration(60) * time.Second}

	logger.Agent.Infow("Encode counter metric", "addr", config.FlagRunAgAddr, "name", "PollCount", "value", metrics.PollCount)

	reqByte, err := json.Marshal(req)

	if err != nil {
		logger.Agent.Error(err.Error(), "event", "encode data")
		return errors.New("encoding data failed")
	}

	buf, errCompress := compressData(reqByte)

	if errCompress != nil {
		logger.Agent.Error(errCompress.Error(), "event", "compress data")
		return errors.New("compress data failed")
	}

	logger.Agent.Infow("Send counter metric", "addr", config.FlagRunAgAddr, "data", buf)
	err = retry.Do(func() error {
		r, _ := http.NewRequest("POST", requestURL, buf)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Content-Encoding", "gzip")
		r.Header.Set("Accept-Encoding", "gzip")
		res, er := client.Do(r)

		if res != nil {
			if erClose := res.Body.Close(); erClose != nil {
				logger.Agent.Error(erClose.Error(), "event", "close response")
				return erClose
			}
		}
		return er
	},
		retry.Attempts(100),
		retry.OnRetry(func(n uint, err error) { logger.Agent.Info("Retrying request after error: %v", err) }))

	if err != nil {
		logger.Agent.Error(err.Error(), "event", "send request")
		return errors.New("request failed")
	}
	return nil
}

func SendRequestJSON() error {

	requestURL := fmt.Sprintf("http://%s/update/", config.FlagRunAgAddr)

	logger.Agent.Infow("Start send metrics")
	if err := sendRequestGaugeJSON(requestURL); err != nil {
		logger.Agent.Error("Error in function sendRequestGaugeJSON")
		return err
	}
	if err := sendRequestCounterJSON(requestURL); err != nil {
		logger.Agent.Error("Error in function sendRequestCounterJSON")
		return err
	}

	return nil
}
