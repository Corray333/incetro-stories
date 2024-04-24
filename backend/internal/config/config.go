package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	// EnvDev is the development environment
	EnvDev = "dev"
	// EnvProd is the production environment
	EnvProd = "prod"
)

type Config struct {
	Env string `yaml:"env"`
	HTTP_Server
}

type HTTP_Server struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	Timeout        int      `yaml:"timeout"`
	IdleTimeout    int      `yaml:"idle_timeout"`
}

// SetupLogger sets up the logger based on the environment
//
// If the environment is dev, it will log in text format and add the source file and line number,
// if the environment is prod, it will log in JSON format
func SetupLogger(env string) {
	var log *slog.Logger
	switch env {
	case EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(log)
}

// LoadConfig loads the configuration from the config file
func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal("error while reading config file: " + err.Error())
	}
	defer configFile.Close()
	var cfg Config
	if err := yaml.NewDecoder(configFile).Decode(&cfg); err != nil {
		log.Fatal("error while decoding config file: " + err.Error())
	}
	return &cfg
}

// Configure loads environtemt variables, the configuration from the config file,
// sets up the logger and returns the configuration.
//
// It takes the path to the environment file as an argument
func Configure(envPath string) *Config {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("error while loading environment file: " + err.Error())
	}
	cfg := LoadConfig()
	SetupLogger(cfg.Env)
	return cfg
}
