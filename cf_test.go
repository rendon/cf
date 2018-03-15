package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
	testcli.Run("make", "clean")
	t.Assert(testcli.Success())

	testcli.Run("make")
	t.Assert(testcli.Success())

	t.AssertNil(os.RemoveAll("arena/"))
	t.AssertNil(os.Mkdir("arena/", 0775))
}

func (t *TestSuite) TestGenAndTest() {
	supportedLangs := map[string]string{
		"c":      "c",
		"cpp":    "cpp",
		"golang": "go",
		"ruby":   "rb",
	}

	root := os.Getenv("PWD")
	cf := fmt.Sprintf("%s%cbin%ccf", root, os.PathSeparator, os.PathSeparator)
	for lang, ext := range supportedLangs {
		dir := fmt.Sprintf("arena/%s", lang)
		srcFile := fmt.Sprintf("%s/%s.%s", dir, lang, ext)
		fmt.Printf("srcFile: %s\n", srcFile)
		testcli.Run(cf, "gen", srcFile)
		if !testcli.Success() {
			fmt.Printf(testcli.Stderr())
			t.Fail("Failed to run `cf gen` for lang " + lang)
		}

		body, err := ioutil.ReadFile(srcFile)
		t.AssertNil(err)
		t.Assert(len(body) > 0)

		t.AssertNil(os.Chdir(dir))
		body, err = ioutil.ReadFile(".settings.yml")
		t.AssertNil(err)
		t.Assert(len(body) > 0)

		testcli.Run(cf, "config", "tests", "1")
		t.Assert(testcli.Success())
		testcli.Run("touch", ".in1.txt", ".out1.txt")
		t.Assert(testcli.Success())

		testcli.Run(cf, "-v", "test")
		if !testcli.Success() {
			fmt.Printf("Error: " + testcli.Stderr())
			t.Fail("Failed to run `cf test` for lang " + lang)
		}
		t.AssertNil(os.Chdir(root))
	}
}

func (t *TestSuite) TestSetup() {
	root := os.Getenv("PWD")
	dir := fmt.Sprintf("%s%carena", root, os.PathSeparator)
	t.AssertNil(os.Chdir(dir))

	cf := fmt.Sprintf("%s%cbin%ccf", root, os.PathSeparator, os.PathSeparator)
	testcli.Run(cf, "setup", "399")
	if !testcli.Success() {
		fmt.Fprintf(os.Stderr, testcli.Stderr())
		t.Fail("Failed to run `cf setup`")
	}
	text := "\nContest directory: CodeforcesRound233Div.2/\n"
	t.AssertContainsStr(testcli.Stdout(), text)
	t.AssertEqualInt(1, strings.Count(testcli.Stdout(), text))
}
