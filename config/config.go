package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
	"log"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`

	App struct {
		LogLevel    string `yaml:"log_level"`
		Port        int    `yaml:"port"`
		Environment string `yaml:"environment"`
	} `yaml:"app"`

	RegexGroups struct {
		BibNumber int `yaml:"bib_number"`
		Position  int `yaml:"position"`
		FullName  int `yaml:"full_name"`
		Time      int `yaml:"time"`
		Points    int `yaml:"points"`
		Regex    string `yaml:"regex"`
		BirthYear    int `yaml:"birth_year"`
		Gender int `yaml:"gender"`
		Nationality int `yaml:"nationality"`
		Team int `yaml:"team"`
	} `yaml:"regexgroups"`

	Race struct {
		RaceID       int    `yaml:"raceid"`
		StageNumber  int    `yaml:"stagenr"`
		FileToRead   string `yaml:"file_to_read"`
		DuplicateFile string `yaml:"duplicate_file"`
		Year         int    `yaml:"year"`
	} `yaml:"race"`
}
// LoadConfig reads the config.yaml file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

func LoadConfigFromEnv() *Config {
	get := func(key string, required bool) string {
		val := os.Getenv(key)
		if val == "" && required {
			log.Fatalf("Missing required env variable: %s", key)
		}
		return val
	}

	getInt := func(key string, fallback int) int {
		val := os.Getenv(key)
		if val == "" {
			return fallback
		}
		n, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("Invalid integer for %s: %v", key, err)
		}
		return n
	}

	cfg := &Config{}

	cfg.Database.Host = get("DB_HOST", true)
	cfg.Database.Port = getInt("DB_PORT", 5432)
	cfg.Database.User = get("DB_USER", true)
	cfg.Database.Password = get("DB_PASSWORD", true)
	cfg.Database.Name = get("DB_NAME", true)
	cfg.Database.SSLMode = get("DB_SSLMODE", false)

	cfg.App.LogLevel = get("APP_LOG_LEVEL", false)
	cfg.App.Port = getInt("APP_PORT", 8080)
	cfg.App.Environment = get("APP_ENV", false)

	// Optional fields (or hardcode in your app)
	cfg.RegexGroups.Regex = get("REGEX", false)

	return cfg
}
