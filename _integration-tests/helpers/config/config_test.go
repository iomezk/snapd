// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2014-2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	check "gopkg.in/check.v1"
)

// Hook up check.v1 into the "go test" runner
func Test(t *testing.T) { check.TestingT(t) }

type ConfigSuite struct{}

var _ = check.Suite(&ConfigSuite{})

func patchStdoutToDevNull(c *check.C) {
	devnull, err := os.Open(os.DevNull)
	c.Assert(err, check.IsNil)
	oldStdout := os.Stdout
	os.Stdout = devnull
	defer func() {
		os.Stdout = oldStdout
	}()
}

func testConfigFileName(c *check.C) string {
	tmpDir, err := ioutil.TempDir("", "")
	c.Assert(err, check.IsNil, check.Commentf(
		"Error creating a temporary directory: %v", err))
	return filepath.Join(tmpDir, "test.config")
}

func testConfigStruct(fileName string) *Config {
	return NewConfig(
		fileName,
		"testrelease", "testchannel", "testtargetrelease", "testtargetchannel",
		true, true)
}
func testConfigContents(fileName string) string {
	return `{` +
		fmt.Sprintf(`"FileName":"%s",`, fileName) +
		`"Release":"testrelease",` +
		`"Channel":"testchannel",` +
		`"TargetRelease":"testtargetrelease",` +
		`"TargetChannel":"testtargetchannel",` +
		`"Update":true,` +
		`"Rollback":true` +
		`}`
}

func (s *ConfigSuite) TestWriteConfig(c *check.C) {
	// Do not print to stdout.
	patchStdoutToDevNull(c)
	configFileName := testConfigFileName(c)

	cfg := testConfigStruct(configFileName)
	cfg.Write()

	writtenConfig, err := ioutil.ReadFile(configFileName)
	c.Assert(err, check.IsNil, check.Commentf("Error reading config: %v", err))
	c.Assert(string(writtenConfig), check.Equals, testConfigContents(configFileName))
}

func (s *ConfigSuite) TestReadConfig(c *check.C) {
	configFileName := testConfigFileName(c)

	configContents := testConfigContents(configFileName)
	ioutil.WriteFile(configFileName, []byte(configContents), 0644)

	cfg, err := ReadConfig(configFileName)

	c.Assert(err, check.IsNil, check.Commentf("Error reading config: %v", err))
	c.Assert(cfg, check.DeepEquals, testConfigStruct(configFileName))
}
