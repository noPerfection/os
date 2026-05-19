// Package env was created for one purpose only: LoadAnyEnv
package env

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sds-framework/os-lib/arg"
	"github.com/sds-framework/os-lib/path"
)

// KeyValue is implemented by types that can be written as .env key/value pairs.
type KeyValue interface {
	MapString() map[string]string
}

// LoadAnyEnv gets the list of all .env file paths in the command line arg.
// Then load them up to the application's environment variables.
//
// The values later will be available via app/config.Config.
//
// The .env files locations are related to the exec path
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
