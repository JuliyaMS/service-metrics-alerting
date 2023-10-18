package storage

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMemStorage_Add(t *testing.T) {
	var memStor = MemStorage{}
	memStor.Init()

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
			value:      "657.123",
			want:       http.StatusOK,
		},
		{
			nameTest:   "Test 2",
			typeMetric: "list",
			nameMetric: "RandomValue",
			value:      "1.0009",
			want:       http.StatusBadRequest,
		},
		{
			nameTest:   "Test 3",
			typeMetric: "counter",
			nameMetric: "PollCounter",
			value:      "21",
			want:       http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			assert.Equal(t, test.want, memStor.Add(test.typeMetric, test.nameMetric, test.value))
		})
	}
}
