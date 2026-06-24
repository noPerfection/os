package arg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing orchestra
type TestArgSuite struct {
	suite.Suite
	flags []string
}

// Make sure that Account is set to five
// before each test
func (suite *TestArgSuite) SetupTest() {
	os.Args = []string{"app"}
	os.Args = append(os.Args, "--plain")
	os.Args = append(os.Args, "--account")
	os.Args = append(os.Args, "--number-key=5")
	os.Args = append(os.Args, "--env=./.test.env")
	os.Args = append(os.Args, "--env=./.other.env")

	suite.flags = []string{
		"plain",
		"account",
		"number-key=5",
		"env=./.test.env",
		"env=./.other.env",
	}
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *TestArgSuite) Test_0_Run() {
	flags := Flags()
	suite.Require().EqualValues(suite.flags, flags)

	suite.True(FlagExist("number-key"))
	suite.Equal("5", FlagValue("number-key"))
	suite.True(FlagExist("plain"))
	suite.True(FlagExist("account"))
	suite.False(FlagExist("./.test.env"))

	paths := EnvPaths()
	suite.Require().Len(paths, 2)
	suite.Require().Equal([]string{"./.test.env", "./.other.env"}, paths)
}

func (suite *TestArgSuite) Test_1_Flag() {
	name := "name"
	value := "value"
	expected := "--name"
	flag := NewFlag(name)
	suite.Require().Equal(expected, flag)

	expected = "--name=value"
	flag = NewFlag(name, value)
	suite.Require().Equal(expected, flag)

	suite.Require().Equal(name, ExtractFlagName(flag))
	suite.Require().Equal(value, ExtractFlagValue(flag))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCommand(t *testing.T) {
	suite.Run(t, new(TestArgSuite))
}
