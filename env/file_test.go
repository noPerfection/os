package env

import (
	"github.com/sds-framework/os-lib/path"
	"os"
	"path/filepath"
	"testing"

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
	os.Args = append(os.Args, test.envPath)

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

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCommand(t *testing.T) {
	suite.Run(t, new(TestEnvSuite))
}
