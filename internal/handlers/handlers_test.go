package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/JuliyaMS/service-metrics-alerting/internal/logger"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body []byte) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	logger.NewLogger()
	ts := httptest.NewServer(NewRouter())
	defer ts.Close()
	randomValues := []float64{127.765, 154789.33200}
	randomDeltas := []int64{657, 325}
	tests := []struct {
		nameTest string
		method   string
		url      string
		body     metrics.Metrics
		status   int
		want     string
	}{
		{
			nameTest: "Test 1",
			method:   "POST",
			url:      "/update/gauge/TestMetric1/12.00006",
			body:     metrics.Metrics{},
			status:   http.StatusOK,
			want:     "",
		},
		{
			nameTest: "Test 2",
			method:   "POST",
			url:      "/update/gauge/12.00006",
			body:     metrics.Metrics{},
			status:   http.StatusNotFound,
			want:     "",
		},
		{
			nameTest: "Test 3",
			method:   "POST",
			url:      "/update/list/TestMetric1/12.00006",
			body:     metrics.Metrics{},
			status:   http.StatusBadRequest,
			want:     "",
		},
		{
			nameTest: "Test 4",
			method:   "POST",
			url:      "/update/counter/TestMetric2/99999",
			body:     metrics.Metrics{},
			status:   http.StatusOK,
			want:     "",
		},
		{
			nameTest: "Test 5",
			method:   "GET",
			url:      "/value/gauge/TestMetric3",
			body:     metrics.Metrics{},
			status:   http.StatusNotFound,
			want:     "",
		},
		{
			nameTest: "Test 6",
			method:   "GET",
			url:      "/value/counter/TestMetric2",
			body:     metrics.Metrics{},
			status:   http.StatusOK,
			want:     "99999",
		},
		{
			nameTest: "Test 7",
			method:   "GET",
			url:      "/value/gauge/TestMetric1",
			body:     metrics.Metrics{},
			status:   http.StatusOK,
			want:     "12.00006",
		},
		{
			nameTest: "Test 8",
			method:   "GET",
			url:      "/value/",
			body:     metrics.Metrics{},
			status:   http.StatusBadRequest,
			want:     "",
		},
		{
			nameTest: "Test 9",
			method:   "POST",
			url:      "/update/",
			body:     metrics.Metrics{MType: "gauge", ID: "RandomValue", Value: &randomValues[0]},
			status:   http.StatusOK,
			want:     "{\"id\":\"RandomValue\",\"type\":\"gauge\",\"value\":127.765}\n",
		},
		{
			nameTest: "Test 9",
			method:   "POST",
			url:      "/value/",
			body:     metrics.Metrics{MType: "gauge", ID: "RandomValue"},
			status:   http.StatusOK,
			want:     "{\"id\":\"RandomValue\",\"type\":\"gauge\",\"value\":127.765}\n",
		},
		{
			nameTest: "Test 10",
			method:   "POST",
			url:      "/update/",
			body:     metrics.Metrics{MType: "counter", ID: "PollCount", Delta: &randomDeltas[0]},
			status:   http.StatusOK,
			want:     "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":657}\n",
		},
		{
			nameTest: "Test 11",
			method:   "POST",
			url:      "/value/",
			body:     metrics.Metrics{MType: "counter", ID: "PollCount"},
			status:   http.StatusOK,
			want:     "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":657}\n",
		},
	}
	for _, test := range tests {
		reqByte, err := json.Marshal(test.body)
		if err != nil {
			panic(err.Error())
		}
		resp, get := testRequest(t, ts, test.method, test.url, reqByte)
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
		resp.Body.Close()
	}
}
