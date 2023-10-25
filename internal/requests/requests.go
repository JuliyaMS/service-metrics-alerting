package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
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
	requestURL := fmt.Sprintf("http://%s/update/counter/PollCounter/%d", config.FlagRunAgAddr, metrics.PollCount)
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

func SendRequestJSON() error {

	req := metrics.Metrics{MType: "gauge"}

	client := http.Client{Timeout: time.Duration(60) * time.Second}
	requestURL := fmt.Sprintf("http://%s/update/", config.FlagRunAgAddr)

	logger.Agent.Infow("Start send metrics")
	for k, v := range metrics.GaugeAgent.ReturnValues() {
		logger.Agent.Infow("Encode gauge metric", "addr", config.FlagRunAgAddr, "name", k, "value", v)
		req.ID = k
		req.Value = &v
		reqByte, err := json.Marshal(req)

		if err != nil {
			logger.Agent.Infow(err.Error(), "event", "encode data")
			return errors.New("encoding data failed")
		}

		logger.Agent.Infow("Send gauge metric", "addr", config.FlagRunAgAddr, "data", string(reqByte))

		res, err := client.Post(requestURL, "Content-Type: application/json", bytes.NewBuffer(reqByte))

		if err != nil {
			logger.Agent.Infow(err.Error(), "event", "send request")
			return errors.New("request failed")
		}
		if er := res.Body.Close(); er != nil {
			logger.Agent.Infow(er.Error(), "event", "close response")
			return er
		}
	}
	logger.Agent.Infow("Encode counter metric", "addr", config.FlagRunAgAddr, "name", "PollCount", "value", metrics.PollCount)
	req.MType = "counter"
	req.ID = "PollCount"
	req.Delta = &metrics.PollCount
	reqByte, err := json.Marshal(req)
	if err != nil {
		logger.Agent.Infow(err.Error(), "event", "encode data")
		return errors.New("encoding data failed")
	}

	logger.Agent.Infow("Send counter metric", "addr", config.FlagRunAgAddr, "data", string(reqByte))
	res, err := client.Post(requestURL, "Content-Type: application/json", bytes.NewBuffer(reqByte))

	if err != nil {
		logger.Agent.Infow(err.Error(), "event", "send request")
		return errors.New("request failed")
	}
	if er := res.Body.Close(); er != nil {
		logger.Agent.Infow(er.Error(), "event", "close response")
		return er
	}
	return nil
}
