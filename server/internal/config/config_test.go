package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileUsesDefaultsWhenMissing(t *testing.T) {
	t.Setenv("SANJIE_SERVER_ADDR", "")
	t.Setenv("SANJIE_DATABASE_DSN", "")
	t.Setenv("SANJIE_DEMO_SEED", "")

	cfg, err := LoadFile(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Addr != ":8080" || cfg.Database.Driver != "sqlite" || cfg.Database.DSN != "sanjie.db" || !cfg.Demo.Seed {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadFileParsesConfigYAML(t *testing.T) {
	t.Setenv("SANJIE_SERVER_ADDR", "")
	t.Setenv("SANJIE_DATABASE_DSN", "")
	t.Setenv("SANJIE_DEMO_SEED", "")

	path := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte("server:\n  addr: \":9090\"\n\ndatabase:\n  driver: \"sqlite\"\n  dsn: \"data/sanjie.db\"\n\ndemo:\n  seed: false\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Addr != ":9090" || cfg.Database.Driver != "sqlite" || cfg.Database.DSN != "data/sanjie.db" || cfg.Demo.Seed {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadFileAppliesEnvironmentOverrides(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte("server:\n  addr: \":9090\"\ndatabase:\n  dsn: \"file.db\"\ndemo:\n  seed: true\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("SANJIE_SERVER_ADDR", ":7070")
	t.Setenv("SANJIE_DATABASE_DSN", "env.db")
	t.Setenv("SANJIE_DEMO_SEED", "false")

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Addr != ":7070" || cfg.Database.DSN != "env.db" || cfg.Demo.Seed {
		t.Fatalf("unexpected env override config: %+v", cfg)
	}
}
