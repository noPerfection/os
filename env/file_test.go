package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/noPerfection/os/path"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing orchestra
type TestEnvSuite struct {
	suite.Suite
	envPath string
}

// Make sure that Account is set to five
// before each test
func (test *TestEnvSuite) SetupTest() {
	currentDir, err := path.CurrentDir()
	test.Require().NoError(err)

	test.envPath = filepath.Join(currentDir, ".test.env")
	os.Args = []string{"app", "--env=" + test.envPath}

	file, err := os.Create(test.envPath)
	test.Require().NoError(err)
	_, err = file.WriteString("")
	test.Require().NoError(err, "failed to write the data into: "+test.envPath)
	err = file.Close()
	test.Require().NoError(err, "delete the dump file: "+test.envPath)
}

func (test *TestEnvSuite) TearDownTest() {
	err := os.Remove(test.envPath)
	test.Require().NoError(err, "delete the dump file: "+test.envPath)
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (test *TestEnvSuite) TestRun() {
	err := LoadAnyEnv()
	test.Require().NoError(err)
}

func TestLoadAnyEnvMultipleFiles(t *testing.T) {
	originalArgs := os.Args
	t.Cleanup(func() {
		os.Args = originalArgs
	})

	currentDir, err := path.CurrentDir()
	require.NoError(t, err)

	alphaPath := filepath.Join(currentDir, ".alpha.env")
	betaPath := filepath.Join(currentDir, ".beta.env")
	t.Cleanup(func() {
		_ = os.Remove(alphaPath)
		_ = os.Remove(betaPath)
		_ = os.Unsetenv("ALPHA_KEY")
		_ = os.Unsetenv("BETA_KEY")
	})

	require.NoError(t, os.WriteFile(alphaPath, []byte("ALPHA_KEY=alpha\n"), 0o644))
	require.NoError(t, os.WriteFile(betaPath, []byte("BETA_KEY=beta\n"), 0o644))

	os.Args = []string{"app", "--env=" + alphaPath, "--env=" + betaPath}
	require.NoError(t, LoadAnyEnv())

	require.Equal(t, "alpha", os.Getenv("ALPHA_KEY"))
	require.Equal(t, "beta", os.Getenv("BETA_KEY"))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCommand(t *testing.T) {
	suite.Run(t, new(TestEnvSuite))
}
