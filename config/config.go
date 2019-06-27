package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	// defaultAuthEnv is the ENV var that will be looked up if the configuration is
	// not in the default location.
	defaultAuthEnv = "METRIC_AUTH_CONFIG"
)

var (
	basePath = "/etc/metric-auth/"
	// defaultConfigPath is the default location to look for the configuration file.
	defaultConfigPath = basePath + "metric-auth.conf"

	defaultWebServer = WebServerConfig{
		ListenAddress: "0.0.0.0",
		Port:          "80",
		UseTLS:        false,
		CertPath:      basePath + "cert.cert",
		KeyPath:       basePath + "cert.key",
	}

	defaultRedisServer = RedisConfig{
		RedisHost:          "172.17.0.1",
		RedisPort:          "6379",
		AvailableEndpoints: []string{},
	}
)

// Config is a struct that holds all the configuration options for the application.
type Config struct {
	Redis     RedisConfig     `json:"redis_server"`
	WebServer WebServerConfig `json:"web_server"`
}

// WebServerConfig holds the configuration for the http web server
type WebServerConfig struct {
	ListenAddress string `json:"listen_address"`
	Port          string `json:"listen_port"`
	UseTLS        bool   `json:"use_tls"`
	CertPath      string `json:"cert_path"`
	KeyPath       string `json:"key_path"`
}

type RedisConfig struct {
	RedisHost          string   `json:"redis_host"`
	RedisPort          string   `json:"redis_port"`
	AvailableEndpoints []string `json:"endpoints"`
}

// New creates a new configuration struct and returns to to the caller.
// The configuration file location is either defaultConfigPath const or
// METRIC_AUTH_CONFIG environment variable.
func New() (*Config, error) {
	configBytes, err := readConfigFile(determinConfigPath())
	if err != nil {
		return nil, err
	}

	c, err := hydrateConfig(configBytes)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func determinConfigPath() string {
	path := defaultConfigPath
	if value, ok := os.LookupEnv(defaultAuthEnv); ok {
		path = value
	}

	return path
}

func readConfigFile(path string) ([]byte, error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return configBytes, nil
}

func hydrateConfig(cb []byte) (*Config, error) {
	// Make new config and set defaults
	config := new(Config)
	config.WebServer = defaultWebServer
	config.Redis = defaultRedisServer

	err := json.Unmarshal(cb, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
