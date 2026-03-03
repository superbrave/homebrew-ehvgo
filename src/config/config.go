package config

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "strings"
)

type AppConfig struct {
    KubeContext   string                    `json:"kubeContext"`
    KubeNamespace string                    `json:"kubeNamespace"`
    Databases     map[string]DatabaseConfig `json:"databases"`
}

type DatabaseConfig struct {
    AwsProfile string `json:"awsProfile"`
    InstanceID string `json:"instanceId"`
    Endpoint   string `json:"endpoint"`
    Port       int    `json:"port"`
    LocalPort  int    `json:"localPort"`
}

// Read loads the config from disk or returns an empty config if missing.
func Read() (AppConfig, error) {
    path, err := Path()
    if err != nil {
        return AppConfig{}, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return AppConfig{Databases: map[string]DatabaseConfig{}}, nil
        }
        return AppConfig{}, err
    }

    var cfg AppConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return AppConfig{}, err
    }

    cfg.KubeContext = strings.TrimSpace(cfg.KubeContext)
    cfg.KubeNamespace = strings.TrimSpace(cfg.KubeNamespace)
    if cfg.Databases == nil {
        cfg.Databases = map[string]DatabaseConfig{}
    }

    return cfg, nil
}

// Write saves the config to disk.
func Write(cfg AppConfig) error {
    path, err := Path()
    if err != nil {
        return err
    }

    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0o700); err != nil {
        return err
    }

    payload, err := json.MarshalIndent(cfg, "", "    ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, payload, 0o600)
}

// Path returns the config path.
func Path() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }

    return filepath.Join(home, ".ehvgo", "config.json"), nil
}
