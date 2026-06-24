package substrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ahmetson/mushroom"
)

func TestSubstrateImplementsMushroomSubstrate(t *testing.T) {
	var _ mushroom.Substrate = (*Substrate)(nil)
}

func TestMyceliumImplementsMushroomMycelium(t *testing.T) {
	var _ mushroom.Mycelium = (*Mycelium)(nil)
}

func root(module string) (*Mycelium, error) {
	return Root("pkg:os/" + module)
}

func TestPathCurrentDir(t *testing.T) {
	mycelium, err := root("path")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	got, err := mycelium.Spore("*pkg:os/path?func=CurrentDir()")
	if err != nil {
		t.Fatalf("Spore returned error: %v", err)
	}

	want, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("Executable returned error: %v", err)
	}
	want = filepath.Dir(exe)

	if got != want {
		t.Fatalf("Spore() = %q, want %q", got, want)
	}
}

func TestPathRejectsHashModuleSyntax(t *testing.T) {
	_, err := Root("pkg:os#path")
	if err == nil {
		t.Fatal("Root(pkg:os#path) returned nil error, want module registration error")
	}
}

func TestPathRejectsVarResource(t *testing.T) {
	mycelium, err := root("path")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	_, err = mycelium.Spore("*pkg:os/path?var")
	if err == nil {
		t.Fatal("Spore(*pkg:os/path?var) returned nil error, want error")
	}
}

func TestPathRejectsUnregisteredFunc(t *testing.T) {
	mycelium, err := root("path")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	_, err = mycelium.Spore("*pkg:os/path?func=MakeDir()")
	if err == nil {
		t.Fatal("Spore(*pkg:os/path?func=MakeDir()) returned nil error, want error")
	}

	_, err = mycelium.Spore("*pkg:os/path?func=fooBar()")
	if err == nil {
		t.Fatal("Spore(*pkg:os/path?func=fooBar()) returned nil error, want error")
	}
}

func TestPathFileNameRequiresParameter(t *testing.T) {
	mycelium, err := root("path")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	_, err = mycelium.Spore("*pkg:os/path?func=FileName()")
	if err == nil {
		t.Fatal("Spore(*pkg:os/path?func=FileName()) returned nil error, want error")
	}
}

func TestPathRejectsNestedModule(t *testing.T) {
	_, err := Root("pkg:os/path#my-module")
	if err == nil {
		t.Fatal("Root(pkg:os/path#my-module) returned nil error, want error")
	}
}

func TestPathFileName(t *testing.T) {
	mycelium, err := root("path")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	got, err := mycelium.Spore("*pkg:os/path?func=FileName(/tmp/example.txt)")
	if err != nil {
		t.Fatalf("Spore returned error: %v", err)
	}
	if got != "example.txt" {
		t.Fatalf("Spore() = %q, want %q", got, "example.txt")
	}
}

func TestProcessCurrentPid(t *testing.T) {
	mycelium, err := root("process")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	got, err := mycelium.Spore("*pkg:os/process?func=CurrentPid")
	if err != nil {
		t.Fatalf("Spore returned error: %v", err)
	}
	if got != uint64(os.Getpid()) {
		t.Fatalf("Spore() = %v, want %d", got, os.Getpid())
	}
}

func TestArgVars(t *testing.T) {
	mycelium, err := root("arg")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	prefix, err := mycelium.Spore("*pkg:os/arg?var=prefix")
	if err != nil {
		t.Fatalf("Spore(prefix) returned error: %v", err)
	}
	if prefix != "--" {
		t.Fatalf("prefix = %q, want %q", prefix, "--")
	}

	sep, err := mycelium.Spore("*pkg:os/arg?var=sep")
	if err != nil {
		t.Fatalf("Spore(sep) returned error: %v", err)
	}
	if sep != "=" {
		t.Fatalf("sep = %q, want %q", sep, "=")
	}
}

func TestEnvVarReturnsEmptyWhenUnset(t *testing.T) {
	mycelium, err := root("env")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	const name = "OS_SUBSTRATE_TEST_UNSET_VAR_12345"
	_ = os.Unsetenv(name)

	got, err := mycelium.Spore("*pkg:os/env?var=" + name)
	if err != nil {
		t.Fatalf("Spore returned error: %v", err)
	}
	if got != "" {
		t.Fatalf("Spore() = %q, want empty string", got)
	}
}

func TestNetGetFreePort(t *testing.T) {
	mycelium, err := root("net")
	if err != nil {
		t.Fatalf("Root returned error: %v", err)
	}

	got, err := mycelium.Spore("*pkg:os/net?func=GetFreePort()")
	if err != nil {
		t.Fatalf("Spore returned error: %v", err)
	}

	port, ok := got.(int)
	if !ok || port <= 0 {
		t.Fatalf("GetFreePort() = %v, want positive int", got)
	}
}

func TestSubstratePattern(t *testing.T) {
	substrate := New()
	if substrate.MushroomURL() != "pkg:os$#$" {
		t.Fatalf("MushroomURL() = %q, want %q", substrate.MushroomURL(), "pkg:os$#$")
	}
}

func TestSowIsUnsupported(t *testing.T) {
	substrate := New()
	soil := &mushroom.Soil{}
	hypha, err := soil.Hypha("pkg:os/path")
	if err != nil {
		t.Fatalf("Hypha returned error: %v", err)
	}
	if err := substrate.Sow(hypha, "data"); err == nil {
		t.Fatal("Sow returned nil error, want unsupported error")
	}
}
