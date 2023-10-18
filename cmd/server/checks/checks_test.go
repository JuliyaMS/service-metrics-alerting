package checks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckDigit(t *testing.T) {
	tests := []struct {
		nameTest   string
		typeMetric string
		value      string
		want       bool
	}{
		{
			nameTest:   "Test 1",
			typeMetric: "gauge",
			value:      "657.123",
			want:       true,
		},
		{
			nameTest:   "Test 2",
			typeMetric: "counter",
			value:      "25",
			want:       true,
		},
		{
			nameTest:   "Test 3",
			typeMetric: "gauge",
			value:      "657,999",
			want:       false,
		},
		{
			nameTest:   "Test 4",
			typeMetric: "gauge",
			value:      "value",
			want:       false,
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			assert.Equal(t, test.want, CheckDigit(test.typeMetric, test.value))
		})
	}
}

func TestCheckType(t *testing.T) {
	tests := []struct {
		nameTest   string
		typeMetric string
		want       bool
	}{
		{
			nameTest:   "Test 1",
			typeMetric: "gauge",
			want:       true,
		},
		{
			nameTest:   "Test 2",
			typeMetric: "counter",
			want:       true,
		},
		{
			nameTest:   "Test 3",
			typeMetric: "type",
			want:       false,
		},
	}
	for _, test := range tests {
		t.Run(test.nameTest, func(t *testing.T) {
			assert.Equal(t, test.want, CheckType(test.typeMetric))
		})
	}
}
