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

// LoadAnyEnv loads environment variables from .env files.
//
// Optional arguments are arbitrary:
//   - string values are treated as .env file paths
//   - true enables envArg, which also reads repeated --env flags from os.Args
//
// envArg defaults to false, so importing this package does not implicitly parse CLI flags.
//
// When no paths are given, it falls back to ".env" in the current directory.
// Paths are resolved relative to the executable directory.
//
// Examples:
//
//	env.LoadAnyEnv()
//	env.LoadAnyEnv(".env", "./beta.env")
//	env.LoadAnyEnv(true) // ./app --env=.env --env=./beta.env
func LoadAnyEnv(params ...any) error {
	envArg := false
	paths := make([]string, 0, len(params))

	for _, param := range params {
		switch value := param.(type) {
		case string:
			if value != "" {
				paths = append(paths, value)
			}
		case bool:
			envArg = value
		default:
			return fmt.Errorf("unsupported argument type %T", param)
		}
	}

	currentDir, err := path.CurrentDir()
	if err != nil {
		return fmt.Errorf("path.CurrentDir: %w", err)
	}

	if envArg {
		paths = append(paths, arg.EnvPaths()...)
	}

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
