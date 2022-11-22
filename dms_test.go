package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDMS(t *testing.T) {
	type test struct {
		input string
		want  float64
	}

	tests := []test{
		{"25\u00b0 8\u2032 6\u2033 S", 25.135},
		{"125° 8′ 6″ Z", -125.135},
		{"90° 0′ 0″ J", -90.0},
		{"125°8′			6″    Z", -125.135},
		{"25\u00b08\u20326\u2033S", 25.135},
		{"25\u00b0    8\u2032  	6\u2033 S", 25.135},
		{"55\u00b0 30\u2032 0\u2033 V", 55.5},
		{"55\u00b0 15\u2032 0\u2033Z", -55.25},
		{"101°\u00a043′\u00a051″ Z", -101.7308},
	}

	for _, tc := range tests {
		d, err := ParseDMS(tc.input)
		assert.NoError(t, err)
		assert.Equal(t, tc.want, d)
	}
}

func TestParseDMSErrors(t *testing.T) {

	type test struct {
		input string
	}

	tests := []test{
		{""},
		{"0 0 0 Z"},
		{"55\u00b0 15 \u2032 0\u2033Z"},
		{"125° 8′ 6″"},
		{"125° 8′ 6″ J"},
		{"55\u00b0   \u2032 0\u2033Z"},
		{"5d 0m 11s S"},
		{"A\u00b0 S"},
		{"55\u00b0 J"},
		{"3° 66′ 0″ Z"},
	}

	for _, tc := range tests {
		d, err := ParseDMS(tc.input)
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
}

func BenchmarkParseDMS(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		_, _ = ParseDMS("125° 8′ 6″ Z")
	}
}
