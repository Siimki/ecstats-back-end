package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
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
