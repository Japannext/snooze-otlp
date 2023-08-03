package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDefaults(t *testing.T) {

  initConfig()

  assert.Equal(t, "0.0.0.0", Config.GrpcListeningAddress)
  assert.Equal(t, 4317, Config.GrpcListeningPort)
  assert.Equal(t, "debug", Config.LogLevel)
  assert.Equal(t, true, Config.PrometheusEnable)
  assert.Equal(t, 9317, Config.PrometheusPort)
}

func TestConfigSnooze(t *testing.T) {

  t.Setenv("SNOOZE_OTLP_SNOOZE_URL", "https://snooze.example.com")

  initConfig()

  assert.Equal(t, "https://snooze.example.com", Config.SnoozeUrl)

}
