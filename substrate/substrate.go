package substrate

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ahmetson/mushroom"
)

type Substrate struct {
	url mushroom.Hypha
	mu  sync.RWMutex
}

var _ mushroom.Substrate = (*Substrate)(nil)

type Mycelium struct {
	url       mushroom.Hypha
	module    string
	soil      *mushroom.Soil
	substrate mushroom.Substrate
}

var _ mushroom.Mycelium = (*Mycelium)(nil)

// New returns an OS substrate registered with the pattern pkg:os/$ .
func New() mushroom.Substrate {
	substrate := &Substrate{}
	soil := &mushroom.Soil{}
	substrate.url, _ = soil.Hypha("pkg:os/$")
	return substrate
}

// Root creates the initial mycelium colony for the given OS module URL.
//
// Example:
//
//	mycelium, err := substrate.Root("pkg:os/path")
//	mycelium, err := substrate.Root("pkg:os/path", otherSubstrate)
func Root(mushroomURL string, substrates ...mushroom.Substrate) (*Mycelium, error) {
	substrate := New()
	soil := &mushroom.Soil{}
	for _, s := range substrates {
		if err := soil.AddSubstrate(s); err != nil {
			return nil, err
		}
	}
	hypha, err := soil.Hypha(mushroomURL)
	if err != nil {
		return nil, err
	}
	got, err := soil.Germinate(hypha, substrate)
	if err != nil {
		return nil, err
	}
	return got.(*Mycelium), nil
}

func (substrate *Substrate) Digest(url mushroom.Hypha, data any, soil *mushroom.Soil) (mushroom.Mycelium, error) {
	if !url.URL {
		return nil, fmt.Errorf("os substrate: digest URL must be a Mushroom URL")
	}
	if url.Dereference {
		return nil, fmt.Errorf("os substrate: digest URL must be a link")
	}
	if !substrate.url.Satisfies(url) {
		return nil, fmt.Errorf("os substrate: digest URL %q does not satisfy %q", url.String(), substrate.url.String())
	}
	if url.ModuleID != "" {
		return nil, fmt.Errorf("os substrate: no module %q is registered", url.ModuleID)
	}
	if !isRegisteredModule(url.PackageID) {
		return nil, fmt.Errorf("os substrate: module %q is not registered", url.PackageID)
	}

	module, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("os substrate: unsupported digest data %T", data)
	}
	if module != url.PackageID {
		return nil, fmt.Errorf("os substrate: foraged module %q does not match URL package %q", module, url.PackageID)
	}

	return &Mycelium{
		url:       url.ModuleURL(),
		module:    url.PackageID,
		soil:      soil,
		substrate: substrate,
	}, nil
}

func (substrate *Substrate) MushroomURL() string {
	return substrate.url.String()
}

func (substrate *Substrate) Forage(url mushroom.Hypha) (any, error) {
	if !substrate.url.Satisfies(url) {
		return nil, fmt.Errorf("os substrate: forage URL %q does not satisfy pattern %q", url.String(), substrate.url.String())
	}
	if url.ModuleID != "" {
		return nil, fmt.Errorf("os substrate: no module %q is registered", url.ModuleID)
	}
	if !isRegisteredModule(url.PackageID) {
		return nil, fmt.Errorf("os substrate: module %q is not registered", url.PackageID)
	}

	substrate.mu.RLock()
	defer substrate.mu.RUnlock()
	return url.PackageID, nil
}

func (substrate *Substrate) Sow(url mushroom.Hypha, data any) error {
	return fmt.Errorf("os substrate: sow is not supported")
}

func (mycelium *Mycelium) Link(path string) (string, error) {
	hypha, err := mycelium.soil.Hypha(path, mycelium.url)
	if err != nil {
		return "", err
	}
	if hypha.Dereference {
		return "", errors.New("os substrate: link cannot contain a dereference")
	}
	if !hypha.URL {
		return path, nil
	}
	if hypha.ModuleID != "" {
		return "", fmt.Errorf("os substrate: no module %q is registered", hypha.ModuleID)
	}
	if _, _, err := mycelium.soil.Recognize(hypha); err != nil {
		return "", err
	}
	return hypha.String(), nil
}

func (mycelium *Mycelium) Spore(path string) (any, error) {
	var hypha mushroom.Hypha
	for {
		var err error
		hypha, err = mycelium.soil.Hypha(path, mycelium.url)
		if err == nil {
			break
		}
		var unrecognized *mushroom.ErrUnrecognizedMycelium
		if !errors.As(err, &unrecognized) {
			return nil, err
		}
		if _, germinateErr := mycelium.soil.Germinate(unrecognized.Hypha, unrecognized.Substrate); germinateErr != nil {
			return nil, fmt.Errorf("os substrate: spore %q: %w", path, germinateErr)
		}
	}

	if !hypha.URL {
		return path, nil
	}
	if !hypha.Dereference {
		return nil, fmt.Errorf("os substrate: spore requires a dereference URL, got link %q", path)
	}
	if hypha.DereferenceType != mushroom.DereferenceTypeResource {
		return nil, fmt.Errorf("os substrate: spore requires a resource dereference, got %q", hypha.DereferenceType)
	}
	if hypha.ModuleID != "" {
		return nil, fmt.Errorf("os substrate: no module %q is registered", hypha.ModuleID)
	}
	if !isRegisteredModule(hypha.PackageID) {
		return nil, fmt.Errorf("os substrate: module %q is not registered", hypha.PackageID)
	}

	colony, substrate, recognizeErr := mycelium.soil.Recognize(hypha)
	if recognizeErr != nil {
		return nil, fmt.Errorf("os substrate: spore %q: %w", path, recognizeErr)
	}

	if colony != nil && colony != mycelium {
		return colony.Spore(path)
	}
	if colony == nil && substrate != nil {
		m, germinateErr := mycelium.soil.Germinate(hypha, substrate)
		if germinateErr != nil {
			return nil, fmt.Errorf("os substrate: spore %q: %w", path, germinateErr)
		}
		return m.Spore(path)
	}

	switch hypha.ResourceKind {
	case mushroom.ResourceKindFunc:
		return dispatchFunc(hypha.PackageID, hypha)
	case mushroom.ResourceKindVar:
		return dispatchVar(hypha.PackageID, hypha)
	case "":
		if hypha.PackageID == "process" {
			if _, ok := hypha.AdditionalProps["func"]; ok {
				return dispatchFunc(hypha.PackageID, hypha)
			}
		}
		return nil, fmt.Errorf("os substrate: resource is not registered")
	default:
		return nil, fmt.Errorf("os substrate: unsupported resource kind %q", hypha.ResourceKind)
	}
}

func (mycelium *Mycelium) Fruit(value any) (any, error) {
	switch typed := value.(type) {
	case string:
		hypha, err := mycelium.soil.Hypha(typed)
		if err != nil {
			return typed, nil
		}
		if hypha.URL && hypha.Dereference {
			return mycelium.Spore(typed)
		}
		return typed, nil
	case map[string]any:
		clone := make(map[string]any, len(typed))
		for key, item := range typed {
			fruited, err := mycelium.Fruit(item)
			if err != nil {
				return nil, err
			}
			clone[key] = fruited
		}
		return clone, nil
	case []any:
		clone := make([]any, len(typed))
		for index, item := range typed {
			fruited, err := mycelium.Fruit(item)
			if err != nil {
				return nil, err
			}
			clone[index] = fruited
		}
		return clone, nil
	default:
		return value, nil
	}
}

func (mycelium *Mycelium) Mineralize() (any, error) {
	return mycelium.module, nil
}

func (mycelium *Mycelium) MushroomURL() string {
	return mycelium.url.String()
}

func (mycelium *Mycelium) MyceliumURL() mushroom.Hypha {
	return mycelium.url
}

func (mycelium *Mycelium) Soil() *mushroom.Soil {
	return mycelium.soil
}

func (mycelium *Mycelium) Substrate() *mushroom.Substrate {
	return &mycelium.substrate
}
