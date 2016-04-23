package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rendon/asserting"
	"github.com/rendon/testcli"
)

type TestSuite struct {
	*asserting.TestCase
}

func TestAll(t *testing.T) {
	ts := &TestSuite{asserting.NewTestCase(t)}
	asserting.Run(ts)
}

func (t *TestSuite) BeforeAll() {
	testcli.Run("make")
	t.Assert(testcli.Success())

	t.AssertNil(os.RemoveAll("arena/"))
	t.AssertNil(os.Mkdir("arena/", 0775))
}

func (t *TestSuite) TestGen() {
	testcli.Run("bin/cf", "gen", "arena/A/A.cpp")
	if !testcli.Success() {
		fmt.Printf(testcli.Stderr())
	}
	t.Assert(testcli.Success())
	body, err := ioutil.ReadFile("arena/A/A.cpp")
	t.AssertNil(err)
	t.Assert(len(body) > 0)

	body, err = ioutil.ReadFile("arena/A/.settings.yml")
	t.AssertNil(err)
	t.Assert(len(body) > 0)
}
