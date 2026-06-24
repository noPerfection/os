// Package env was created for one purpose only: LoadAnyEnv
package env

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/noPerfection/os/arg"
	"github.com/noPerfection/os/path"
)

// KeyValue is implemented by types that can be written as .env key/value pairs.
type KeyValue interface {
	MapString() map[string]string
}

// LoadAnyEnv loads environment variables from repeated --env flags.
// Example: ./app --env=.env --env=./beta.env
//
// When no --env flags are given, it falls back to ".env" in the current directory.
// Paths are resolved relative to the executable directory.
//
// The values later will be available via app/config.Config.
func LoadAnyEnv() error {
	currentDir, err := path.CurrentDir()
	if err != nil {
		return fmt.Errorf("path.CurrentDir: %w", err)
	}

	paths := arg.EnvPaths()
	for i, envPath := range paths {
		paths[i] = path.AbsDir(currentDir, envPath)
	}

	if len(paths) == 0 {
		err = godotenv.Load()
		if err != nil {
			return fmt.Errorf("godotenv.Load(\".env\"): %w", err)
		}
		return nil
	}

	err = godotenv.Load(paths...)
	if err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}
	return nil
}

// WriteEnv writes the given key value to the file.
// If the file exists, then it will be truncated.
func WriteEnv(data KeyValue, path string) error {
	err := godotenv.Write(data.MapString(), path)
	if err != nil {
		return fmt.Errorf("godotenv.Write: %w", err)
	}

	return nil
}
