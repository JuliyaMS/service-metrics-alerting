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

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.DisableKeepAlives = true

	req := metrics.Metrics{MType: "gauge"}

	client := http.Client{Timeout: time.Duration(60) * time.Second, Transport: t}

	for k, v := range metrics.GaugeAgent.ReturnValues() {
		fmt.Println(k, v)
		req.ID = k
		req.Value = &v
		reqByte, err := json.Marshal(req)
		if err != nil {
			fmt.Println(err)
			return errors.New("encoding data failed")
		}
		fmt.Println(string(reqByte))
		requestURL := fmt.Sprintf("http://%s/update/gauge/%s/%f", config.FlagRunAgAddr, k, v)
		logger.Logger.Infow("Send request", "addr", config.FlagRunAgAddr)
		var req2, _ = http.NewRequest("POST", requestURL, nil)

		res, err := client.Do(req2)
		fmt.Println(res, err)

		if err != nil {
			fmt.Println(err)
			return errors.New("request failed")
		}
		if er := res.Body.Close(); er != nil {
			return er
		}
	}

	requestURL := fmt.Sprintf("http://%s/update/", config.FlagRunAgAddr)
	req.MType = "counter"
	req.ID = "PollCount"
	req.Delta = &metrics.PollCount
	reqByte, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return errors.New("encoding data failed")
	}
	res, err := client.Post(requestURL, "Content-Type: application/json", bytes.NewBuffer(reqByte))

	if err != nil {
		return errors.New("request failed")
	}
	if er := res.Body.Close(); er != nil {
		return er
	}
	return nil
}
