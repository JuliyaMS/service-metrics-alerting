package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	ts := httptest.NewServer(NewRouter())
	defer ts.Close()

	tests := []struct {
		nameTest string
		method   string
		url      string
		status   int
		want     string
	}{
		{
			nameTest: "Test 1",
			method:   "POST",
			url:      "/update/gauge/TestMetric1/12.00006",
			status:   http.StatusOK,
			want:     "",
		},
		{
			nameTest: "Test 2",
			method:   "POST",
			url:      "/update/gauge/12.00006",
			status:   http.StatusNotFound,
			want:     "",
		},
		{
			nameTest: "Test 3",
			method:   "POST",
			url:      "/update/list/TestMetric1/12.00006",
			status:   http.StatusBadRequest,
			want:     "",
		},
		{
			nameTest: "Test 4",
			method:   "POST",
			url:      "/update/counter/TestMetric2/99999",
			status:   http.StatusOK,
			want:     "",
		},
		{
			nameTest: "Test 5",
			method:   "GET",
			url:      "/value/gauge/TestMetric3",
			status:   http.StatusNotFound,
			want:     "",
		},
		{
			nameTest: "Test 6",
			method:   "GET",
			url:      "/value/counter/TestMetric2",
			status:   http.StatusOK,
			want:     "99999",
		},
		{
			nameTest: "Test 7",
			method:   "GET",
			url:      "/value/gauge/TestMetric1",
			status:   http.StatusOK,
			want:     "12.00006",
		},
		{
			nameTest: "Test 8",
			method:   "GET",
			url:      "/value/",
			status:   http.StatusBadRequest,
			want:     "",
		},
	}
	for _, test := range tests {
		resp, get := testRequest(t, ts, test.method, test.url)
		assert.Equal(t, test.status, resp.StatusCode)
		assert.Equal(t, test.want, get)
		resp.Body.Close()
	}
}
