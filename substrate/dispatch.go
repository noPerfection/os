package substrate

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ahmetson/mushroom"
	"github.com/noPerfection/os/arg"
	"github.com/noPerfection/os/env"
	osnet "github.com/noPerfection/os/net"
	"github.com/noPerfection/os/path"
	"github.com/noPerfection/os/process"
)

var registeredModules = map[string]struct{}{
	"path":    {},
	"env":     {},
	"net":     {},
	"process": {},
	"arg":     {},
}

func isRegisteredModule(name string) bool {
	_, ok := registeredModules[name]
	return ok
}

func dispatchFunc(module string, hypha mushroom.Hypha) (any, error) {
	call, err := funcCall(hypha)
	if err != nil {
		return nil, err
	}

	switch module {
	case "path":
		return dispatchPathFunc(call)
	case "process":
		return dispatchProcessFunc(call)
	case "net":
		return dispatchNetFunc(call)
	case "arg":
		return dispatchArgFunc(call)
	default:
		return nil, fmt.Errorf("os substrate: module %q is not registered", module)
	}
}

func dispatchVar(module string, hypha mushroom.Hypha) (any, error) {
	switch module {
	case "env":
		return dispatchEnvVar(hypha)
	case "arg":
		return dispatchArgVar(hypha)
	default:
		return nil, fmt.Errorf("os substrate: module %q does not support var resources", module)
	}
}

func funcCall(hypha mushroom.Hypha) (mushroom.ResourceCall, error) {
	if hypha.ResourceKind == mushroom.ResourceKindFunc {
		if len(hypha.ResourcePath.Segments) == 0 || hypha.ResourcePath.Segments[0].Call == nil {
			return mushroom.ResourceCall{}, fmt.Errorf("os substrate: func resource is not registered")
		}
		return *hypha.ResourcePath.Segments[0].Call, nil
	}

	if name, ok := hypha.AdditionalProps["func"]; ok && name != "" {
		return mushroom.ResourceCall{Name: name}, nil
	}

	return mushroom.ResourceCall{}, fmt.Errorf("os substrate: func resource is not registered")
}

func dispatchPathFunc(call mushroom.ResourceCall) (any, error) {
	switch call.Name {
	case "CurrentDir":
		if len(call.Args) != 0 {
			return nil, fmt.Errorf("os substrate: CurrentDir() takes no arguments")
		}
		return path.CurrentDir()
	case "FileName":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: FileName() requires a parameter")
		}
		return path.FileName(scalarString(call.Args[0])), nil
	case "NoExtension":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: NoExtension() requires a parameter")
		}
		return path.NoExtension(scalarString(call.Args[0])), nil
	case "DirAndFileName":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: DirAndFileName() requires a parameter")
		}
		dir, name := path.DirAndFileName(scalarString(call.Args[0]))
		return []any{dir, name}, nil
	case "FileExist":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: FileExist() requires a parameter")
		}
		return path.FileExist(scalarString(call.Args[0]))
	case "DirExist":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: DirExist() requires a parameter")
		}
		return path.DirExist(scalarString(call.Args[0]))
	case "MakeDir":
		return nil, fmt.Errorf("os substrate: func %q is not registered", call.Name)
	default:
		return nil, fmt.Errorf("os substrate: func %q is not registered", call.Name)
	}
}

func dispatchProcessFunc(call mushroom.ResourceCall) (any, error) {
	switch call.Name {
	case "CurrentPid":
		if len(call.Args) != 0 {
			return nil, fmt.Errorf("os substrate: CurrentPid() takes no arguments")
		}
		return process.CurrentPid(), nil
	case "PortToPid":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: PortToPid() requires a port parameter")
		}
		port, err := scalarInt(call.Args[0])
		if err != nil {
			return nil, err
		}
		return process.PortToPid(port)
	default:
		return nil, fmt.Errorf("os substrate: func %q is not registered", call.Name)
	}
}

func dispatchNetFunc(call mushroom.ResourceCall) (any, error) {
	switch call.Name {
	case "GetFreePort":
		if len(call.Args) != 0 {
			return nil, fmt.Errorf("os substrate: GetFreePort() takes no arguments")
		}
		return osnet.GetFreePort(), nil
	case "IsPortUsed":
		if len(call.Args) != 2 {
			return nil, fmt.Errorf("os substrate: IsPortUsed() requires host and port parameters")
		}
		host := scalarString(call.Args[0])
		port, err := scalarInt(call.Args[1])
		if err != nil {
			return nil, err
		}
		return osnet.IsPortUsed(host, port), nil
	default:
		return nil, fmt.Errorf("os substrate: func %q is not registered", call.Name)
	}
}

func dispatchArgFunc(call mushroom.ResourceCall) (any, error) {
	switch call.Name {
	case "NewFlag":
		if len(call.Args) == 0 {
			return nil, fmt.Errorf("os substrate: NewFlag() requires a name parameter")
		}
		name := scalarString(call.Args[0])
		if len(call.Args) == 1 {
			return arg.NewFlag(name), nil
		}
		return arg.NewFlag(name, scalarString(call.Args[1])), nil
	case "Flags":
		if len(call.Args) != 0 {
			return nil, fmt.Errorf("os substrate: Flags() takes no arguments")
		}
		return arg.Flags(), nil
	case "IsFlag":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: IsFlag() requires a parameter")
		}
		return arg.IsFlag(scalarString(call.Args[0])), nil
	case "FlagExist":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: FlagExist() requires a name parameter")
		}
		return arg.FlagExist(scalarString(call.Args[0])), nil
	case "ExtractFlagName":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: ExtractFlagName() requires a parameter")
		}
		return arg.ExtractFlagName(scalarString(call.Args[0])), nil
	case "ExtractFlagValue":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: ExtractFlagValue() requires a parameter")
		}
		return arg.ExtractFlagValue(scalarString(call.Args[0])), nil
	case "FlagValue":
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("os substrate: FlagValue() requires a name parameter")
		}
		return arg.FlagValue(scalarString(call.Args[0])), nil
	case "EnvPaths":
		if len(call.Args) != 0 {
			return nil, fmt.Errorf("os substrate: EnvPaths() takes no arguments")
		}
		return arg.EnvPaths(), nil
	default:
		return nil, fmt.Errorf("os substrate: func %q is not registered", call.Name)
	}
}

func dispatchEnvVar(hypha mushroom.Hypha) (any, error) {
	if len(hypha.ResourcePath.Segments) == 0 {
		return nil, fmt.Errorf("os substrate: env var name is required")
	}

	if err := runLoadAnyEnv(hypha.AdditionalProps); err != nil {
		return nil, err
	}

	name := hypha.ResourcePath.Segments[0].Name
	return os.Getenv(name), nil
}

func dispatchArgVar(hypha mushroom.Hypha) (any, error) {
	if len(hypha.ResourcePath.Segments) == 0 {
		return nil, fmt.Errorf("os substrate: var name is required")
	}

	switch hypha.ResourcePath.Segments[0].Name {
	case "sep":
		return arg.Sep, nil
	case "prefix":
		return arg.Prefix, nil
	default:
		return nil, fmt.Errorf("os substrate: var %q is not registered", hypha.ResourcePath.Segments[0].Name)
	}
}

func runLoadAnyEnv(props map[string]string) error {
	if props == nil {
		return nil
	}

	if value, ok := props["arg"]; ok && value == "true" {
		return env.LoadAnyEnv(true)
	}

	loadAny, hasLoadAny := props["LoadAnyEnv"]
	if !hasLoadAny {
		return nil
	}

	switch loadAny {
	case "true":
		return env.LoadAnyEnv(true)
	case "false":
		return env.LoadAnyEnv(false)
	default:
		return env.LoadAnyEnv()
	}
}

func scalarString(scalar mushroom.ResourceScalar) string {
	switch scalar.Kind {
	case mushroom.ResourceScalarKeyValue:
		return scalar.Value
	case mushroom.ResourceScalarNumber:
		return scalar.Value
	case mushroom.ResourceScalarKey:
		return scalar.Key
	default:
		return scalar.Value
	}
}

func scalarInt(scalar mushroom.ResourceScalar) (int, error) {
	raw := scalarString(scalar)
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("os substrate: invalid integer argument %q", raw)
	}
	return value, nil
}
