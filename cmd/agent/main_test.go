package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestSendRequest(t *testing.T) {
	var rtm runtime.MemStats
	changeMetrics(rtm)

	tests := []struct {
		nameTest string
		port     int
		want     error
	}{
		{
			nameTest: "Test 1",
			port:     8080,
			want:     errors.New("request failed"),
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
