package server

import (

  log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ConfigModel struct {
  // IP address the gRPC server should bind to
  GrpcListeningAddress string `mapstructure:"GRPC_LISTENING_ADDRESS"`
  // Port number the gRPC server should bind to
  GrpcListeningPort int `mapstructure:"GRPC_LISTENING_PORT"`
  // Snooze URL to send alerts to. Will use the
  // path <url>/api/alerts to send alerts
  SnoozeUrl string `mapstructure:"SNOOZE_URL"`
  // Path to a pem formatted certificate authority when
  // communicating with the snooze server in HTTPS
  SnoozeCaPath string `mapstructure:"SNOOZE_CA_PATH"`
  // A logrus log level (trace/debug/info/warning/error/fatal/panic).
  LogLevel string `mapstructure:"LOG_LEVEL"`
  // Whether to enable prometheus metrics
  PrometheusEnable bool `mapstructure:"PROMETHEUS_ENABLE"`
  // Port the prometheus exporter should listen to
  PrometheusPort int `mapstructure:"PROMETHEUS_PORT"`
}

var Config ConfigModel

// Configuration defaults
/*
var Config = &ConfigModel{
  GrpcListeningAddress: "0.0.0.0",
  GrpcListeningPort: 4317,
  LogLevel: "debug",
  PrometheusEnable: true,
  PrometheusPort: 9317,
}
*/

// Declare the default variables.
// Note: Variables with no default should be
// declared with a BindEnv.
func setDefaults() {
  viper.SetDefault("GRPC_LISTENING_ADDRESS", "0.0.0.0")
  viper.SetDefault("GRPC_LISTENING_PORT", 4317)
  viper.BindEnv("SNOOZE_URL")
  viper.BindEnv("SNOOZE_CA_PATH")
  viper.SetDefault("LOG_LEVEL", "debug")
  viper.SetDefault("PROMETHEUS_ENABLE", true)
  viper.SetDefault("PROMETHEUS_PORT", 9317)
}

func initConfig() {
  viper.SetEnvPrefix("SNOOZE_OTLP")
  setDefaults()
  viper.AutomaticEnv()

  err = viper.Unmarshal(&Config)
  if err != nil {
    log.Fatalf("Error unmarshaling config: %s", err)
  }

  log.Debugf("Loaded config: %+v", Config)
}
