package requests

import (
	"errors"
	"fmt"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
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
	requestURL := fmt.Sprintf("http://%s/update/counter/PollCounter/%d", config.FlagRunAgAddr, metrics.PollCounter)
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
