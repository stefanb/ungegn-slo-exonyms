package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDMS(t *testing.T) {
	{
		d, err := ParseDMS("25\u00b0 8\u2032 6\u2033 S")
		assert.NoError(t, err)
		assert.Equal(t, 25.135, d)
	}
	{
		d, err := ParseDMS("125° 8′ 6″ Z")
		assert.NoError(t, err)
		assert.Equal(t, -125.135, d)
	}
	{
		d, err := ParseDMS("90° 0′ 0″ J")
		assert.NoError(t, err)
		assert.Equal(t, -90.0, d)
	}
	{
		d, err := ParseDMS("125°8′			6″    Z")
		assert.NoError(t, err)
		assert.Equal(t, -125.135, d)
	}
	{
		d, err := ParseDMS("25\u00b08\u20326\u2033S")
		assert.NoError(t, err)
		assert.Equal(t, 25.135, d)
	}
	{
		d, err := ParseDMS("25\u00b0    8\u2032  	6\u2033 S")
		assert.NoError(t, err)
		assert.Equal(t, 25.135, d)
	}
	{
		d, err := ParseDMS("55\u00b0 30\u2032 0\u2033 V")
		assert.NoError(t, err)
		assert.Equal(t, 55.5, d)
	}
	{
		d, err := ParseDMS("55\u00b0 15\u2032 0\u2033Z")
		assert.NoError(t, err)
		assert.Equal(t, -55.25, d)
	}
	{
		// non-breaking spaces
		d, err := ParseDMS("101°\u00a043′\u00a051″ Z")
		assert.NoError(t, err)
		assert.Equal(t, -101.7308, d)
	}
}

func TestParseDMSErrors(t *testing.T) {
	{
		d, err := ParseDMS("")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("0 0 0 Z")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("55\u00b0 15 \u2032 0\u2033Z")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("125° 8′ 6″")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("125° 8′ 6″ J")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("55\u00b0   \u2032 0\u2033Z")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("5d 0m 11s S")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("A\u00b0 S")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("55\u00b0 J")
		assert.Error(t, err)
		assert.Equal(t, 0.0, d)
	}
	{
		d, err := ParseDMS("3° 66′ 0″ Z")
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
