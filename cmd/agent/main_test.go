package main

import (
	"errors"
	"github.com/JuliyaMS/service-metrics-alerting/internal/headers"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"regexp"
	"runtime"
	"strconv"
	"testing"
)

func TestSendRequest(t *testing.T) {
	srv := httptest.NewServer(headers.Router())
	defer srv.Close()

	var rtm runtime.MemStats
	changeMetrics(rtm)

	r, _ := regexp.Compile(`:\d+`)
	port, _ := strconv.Atoi(r.FindAllString(srv.URL, -1)[0][1:])

	tests := []struct {
		nameTest string
		port     int
		want     error
	}{
		{
			nameTest: "Test 1",
			port:     port,
			want:     nil,
		},
		{
			nameTest: "Test 2",
			port:     8070,
			want:     errors.New("request failed"),
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			assert.Equal(t, test.want, SendRequest(test.port))
		})
	}
}
