package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Demo     DemoConfig
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type DemoConfig struct {
	Seed bool
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Addr: ":8080",
		},
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "sanjie.db",
		},
		Demo: DemoConfig{
			Seed: true,
		},
	}
}

func Load() (Config, error) {
	path := strings.TrimSpace(os.Getenv("SANJIE_CONFIG"))
	if path == "" {
		path = "config.yaml"
	}
	return LoadFile(path)
}

func LoadFile(path string) (Config, error) {
	cfg := Default()
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return Config{}, fmt.Errorf("read config %s: %w", path, err)
			}
		} else if err := applyYAML(&cfg, data); err != nil {
			return Config{}, fmt.Errorf("parse config %s: %w", path, err)
		}
	}
	if err := applyEnv(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func applyEnv(cfg *Config) error {
	if value := strings.TrimSpace(os.Getenv("SANJIE_SERVER_ADDR")); value != "" {
		cfg.Server.Addr = value
	}
	if value := strings.TrimSpace(os.Getenv("SANJIE_DATABASE_DSN")); value != "" {
		cfg.Database.DSN = value
	}
	if value := strings.TrimSpace(os.Getenv("SANJIE_DEMO_SEED")); value != "" {
		seed, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("parse SANJIE_DEMO_SEED: %w", err)
		}
		cfg.Demo.Seed = seed
	}
	return nil
}

func applyYAML(cfg *Config, data []byte) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	section := ""
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		raw := strings.TrimRight(scanner.Text(), " \t")
		line := stripComment(raw)
		if strings.TrimSpace(line) == "" {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			if !strings.HasSuffix(trimmed, ":") {
				return fmt.Errorf("line %d: expected section", lineNo)
			}
			section = strings.TrimSuffix(trimmed, ":")
			continue
		}
		if section == "" {
			return fmt.Errorf("line %d: key without section", lineNo)
		}
		key, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			return fmt.Errorf("line %d: expected key value", lineNo)
		}
		if err := setValue(cfg, section, strings.TrimSpace(key), cleanValue(value)); err != nil {
			return fmt.Errorf("line %d: %w", lineNo, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func setValue(cfg *Config, section string, key string, value string) error {
	switch section + "." + key {
	case "server.addr":
		cfg.Server.Addr = value
	case "database.driver":
		cfg.Database.Driver = value
	case "database.dsn":
		cfg.Database.DSN = value
	case "demo.seed":
		seed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		cfg.Demo.Seed = seed
	}
	return nil
}

func cleanValue(value string) string {
	value = strings.TrimSpace(value)
	if unquoted, err := strconv.Unquote(value); err == nil {
		return unquoted
	}
	return strings.Trim(value, `"'`)
}

func stripComment(line string) string {
	inSingle := false
	inDouble := false
	for i, r := range line {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return line[:i]
			}
		}
	}
	return line
}
