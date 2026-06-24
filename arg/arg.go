// Package arg is used to read command line flags of the application.
//
// The SDS supports a Flag which is composed of name and optionally a value.
//
// Flags() returns all flags
// EnvPaths() returns paths from repeated --env flags.
// FlagExist(name) flag exists?
// ExtractFlagValue(flag) finds the flag and returns the value of it.
// ExtractFlagName(flag) returns the flag name.
// IsFlag(str) returns true if the given string has a prefix
// NewFlag(name string, values ...string) returns a new flag by its name
// FlagValue(name) returns the flag value
package arg

import (
	"os"
	"strings"
)

const (
	Prefix = "--"
	Sep    = "="
)

// NewFlag creates a new flag with the given name and optionally with a value
func NewFlag(name string, values ...string) string {
	flag := Prefix + name
	if len(values) > 0 {
		flag += Sep + values[0]
	}

	return flag
}

// Flags returns the flags from application flags.
func Flags() []string {
	args := os.Args[1:]
	if len(args) == 0 {
		return []string{}
	}

	count := 0
	for _, str := range args {
		if IsFlag(str) {
			count++
		}
	}

	flags := make([]string, count)

	i := 0
	for _, str := range args {
		if IsFlag(str) {
			flags[i] = strings.TrimPrefix(str, Prefix)
			i++
		}
	}

	return flags
}

// IsFlag returns true, if the given string contains a flag prefix
func IsFlag(str string) bool {
	return strings.HasPrefix(str, Prefix)
}

// FlagExist is given flag exists or not.
func FlagExist(name string) bool {
	flags := Flags()
	for _, flag := range flags {
		if ExtractFlagName(flag) == name {
			return true
		}
	}

	return false
}

// ExtractFlagName returns the flag name.
// If the flag is prefixed, then it will be trimmed.
func ExtractFlagName(flag string) string {
	return strings.Split(strings.TrimPrefix(flag, Prefix), Sep)[0]
}

// ExtractFlagValue Extracts the value of the arg if it exists.
func ExtractFlagValue(flag string) string {
	parts := strings.Split(flag, Sep)
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

func FlagValue(name string) string {
	names := Flags()
	for _, flag := range names {
		if ExtractFlagName(flag) == name {
			return ExtractFlagValue(flag)
		}
	}

	return ""
}

// EnvPaths returns paths from repeated --env flags.
func EnvPaths() []string {
	paths := make([]string, 0)

	for _, flag := range Flags() {
		if ExtractFlagName(flag) != "env" {
			continue
		}

		value := ExtractFlagValue(flag)
		if value != "" {
			paths = append(paths, value)
		}
	}

	return paths
}
