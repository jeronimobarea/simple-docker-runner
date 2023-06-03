package env

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEnvWithFallback(t *testing.T) {
	t.Run("get correct value from env variable", func(t *testing.T) {
		var (
			key   = "test"
			value = "test"
		)
		t.Setenv(key, value)
		got := GetEnvWithFallback(key, "")
		require.Equal(t, value, got)
	})

	t.Run("get fallback value", func(t *testing.T) {
		var (
			key      = "fallback"
			fallback = "fallback"
		)
		got := GetEnvWithFallback(key, fallback)
		require.Equal(t, fallback, got)
	})
}
