package sdk

import (
	"os"
	"testing"
)

func Test_StateDir(t *testing.T) {
	sdkStateDir = []string{
		"SDK_STATE_DIR",
	}

	expected := "test"

	err := os.Setenv(sdkStateDir[0], expected)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := Sdk{}

	output, err := s.StateDir()
	if err != nil {
		t.Errorf("Error getting statedir: %s", err)
	}

	if output != "test" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_HomeDir(t *testing.T) {
	sdkStateDir = []string{
		"SDK_HOME_DIR",
	}

	expected := "test"

	err := os.Setenv(sdkStateDir[0], expected)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := Sdk{}

	output, err := s.HomeDir()
	if err != nil {
		t.Errorf("Error getting statedir: %s", err)
	}

	if output != "test" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_ConfigDir(t *testing.T) {
	sdkStateDir = []string{
		"SDK_CONFIG_DIR",
	}

	expected := "test"

	err := os.Setenv(sdkStateDir[0], expected)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := Sdk{}

	output, err := s.ConfigDir()
	if err != nil {
		t.Errorf("Error getting statedir: %s", err)
	}

	if output != "test" {
		t.Errorf("Error unexpected output: %v", output)
	}
}
