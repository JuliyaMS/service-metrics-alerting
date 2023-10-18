package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestPost(t *testing.T) {
	tests := []struct {
		nameTest   string
		typeMetric string
		nameMetric string
		value      string
		want       int
	}{
		{
			nameTest:   "Test 1",
			typeMetric: "gauge",
			nameMetric: "RandomValue",
			value:      "21.0000005",
			want:       http.StatusOK,
		},
		{
			nameTest:   "Test 2",
			typeMetric: "gauge",
			nameMetric: "",
			value:      "",
			want:       http.StatusNotFound,
		},
		{
			nameTest:   "Test 3",
			typeMetric: "list",
			nameMetric: "RandomValue",
			value:      "21.0000005",
			want:       http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			msg := fmt.Sprintf(`/update/%s/%s/%s/`, test.typeMetric, test.nameMetric, test.value)
			request := httptest.NewRequest(http.MethodPost, msg, nil)
			w := httptest.NewRecorder()
			RequestPost(w, request)
			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, test.want, res.StatusCode)
		})
	}
}
