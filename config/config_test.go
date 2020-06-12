package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDefaultPortIs8080(t *testing.T) {
	config := New()
	assert.Equal(t, 8080, config.Port())
}

func TestPortCanBeSetWithViper(t *testing.T) {
	viper.Set("port", 500)
	config := New()
	assert.Equal(t, 500, config.Port())
	viper.Set("port", nil)
}

func TestPortCanBeSetWithEnvironmentVariable(t *testing.T) {
	os.Setenv("LAZY_REST_PORT", "10080")
	Initialize()
	config := New()
	assert.Equal(t, 10080, config.Port())
	os.Setenv("LAZY_REST_PORT", "")
}

func TestDefaultInMemoryIsFalse(t *testing.T) {
	config := New()
	assert.Equal(t, false, config.InMemory())
}

func TestInMemoryCanBeSetWithViper(t *testing.T) {
	viper.Set("in-memory", true)
	config := New()
	assert.Equal(t, true, config.InMemory())
	viper.Set("in-memory", nil)
}

func TestInMemoryCanBeSetWithEnvironmentVariable(t *testing.T) {
	os.Setenv("LAZY_REST_IN_MEMORY", "true")
	Initialize()
	config := New()
	assert.Equal(t, true, config.InMemory())
	os.Setenv("LAZY_REST_IN_MEMORY", "")
}
