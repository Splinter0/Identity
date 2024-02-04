package endpoints

import (
	"log"
	"os"

	"github.com/Splinter0/identity/bankid"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Providers []string             `yaml:"providers"`
	Service   *string              `yaml:"service"`
	BankID    *bankid.BankIDConfig `yaml:"bankid,omitempty"`
}

func LoadConfig() *Config {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	if config.Service == nil {
		log.Fatal("Must choose a name for 'service' in config")
	}

	if len(config.Providers) == 0 {
		log.Fatal("No provider configured!")
	}

	return &config
}
