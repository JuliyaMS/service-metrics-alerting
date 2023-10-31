package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
		want       error
	}{
		{
			nameTest:   "Test 1",
			typeMetric: "gauge",
			nameMetric: "RandomValue",
			value:      "657.123",
			want:       nil,
		},
		{
			nameTest:   "Test 2",
			typeMetric: "list",
			nameMetric: "RandomValue",
			value:      "1.0009",
			want:       errors.New("this type of metric doesn't exists"),
		},
		{
			nameTest:   "Test 3",
			typeMetric: "counter",
			nameMetric: "PollCounter",
			value:      "21",
			want:       nil,
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			assert.Equal(t, test.want, memStor.Add(test.typeMetric, test.nameMetric, test.value))
		})
	}
}
