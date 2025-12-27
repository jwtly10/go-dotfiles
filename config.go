package main

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const MigrationFile = "migrate.yaml"
const ConfigFile = "dotfiles.yaml"

type DotfilesConfig struct {
	Ignore []string `yaml:"ignore, omitempty"`
}

type MigrateConfig struct {
	Paths  []string `yaml:"paths"`
	Ignore []string `yaml:"ignore, omitempty"`
}

func DefaultDotfilesConfig() *DotfilesConfig {
	return &DotfilesConfig{
		Ignore: []string{
			MigrationFile,
			".DS_Store",
			"*.log",
		},
	}
}

func DefaultMigrateConfig() *MigrateConfig {
	return &MigrateConfig{
		Ignore: []string{
			".DS_Store",
			"*.log",
		},
	}
}

func LoadDotfilesConfig(dotfilesDir string) (*DotfilesConfig, error) {
	path := filepath.Join(dotfilesDir, ConfigFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultDotfilesConfig(), nil
	}
	var config DotfilesConfig
	err := marshal(path, &config)
	return &config, err
}

func (c *DotfilesConfig) Save(dotfilesDir string) error {
	path := filepath.Join(dotfilesDir, ConfigFile)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, FilePermReadWriteUser)
}

func LoadMigrationConfig(dotfilesDir string) (*MigrateConfig, error) {
	path := filepath.Join(dotfilesDir, MigrationFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultMigrateConfig(), nil
	}
	var config MigrateConfig
	err := marshal(path, &config)
	return &config, err
}

func (c *MigrateConfig) Save(dotfilesDir string) error {
	path := filepath.Join(dotfilesDir, MigrationFile)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, FilePermReadWriteUser)
}

func marshal(path string, out interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}
