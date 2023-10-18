package requests

import (
	"errors"
	"github.com/JuliyaMS/service-metrics-alerting/internal/config"
	"github.com/JuliyaMS/service-metrics-alerting/internal/handlers"
	"github.com/JuliyaMS/service-metrics-alerting/internal/metrics"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"runtime"
	"testing"
)

func TestSendRequest(t *testing.T) {
	srv := httptest.NewServer(handlers.NewRouter())
	defer srv.Close()

	var rtm runtime.MemStats
	metrics.ChangeMetrics(&rtm)

	tests := []struct {
		nameTest string
		addr     string
		want     error
	}{
		{
			nameTest: "Test 1",
			addr:     srv.URL[7:],
			want:     nil,
		},
		{
			nameTest: "Test 2",
			addr:     ":8080",
			want:     errors.New("request failed"),
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			config.FlagRunAgAddr = test.addr
			assert.Equal(t, test.want, SendRequest())
		})
	}
}
