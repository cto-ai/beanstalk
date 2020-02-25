package utils

import (
	"os"
	"os/exec"
	"testing"
)

func Test_SetFromEnvValOrFallback_Positive(t *testing.T) {
	test := []string{
		"TEST",
	}

	expected := "test"

	err := os.Setenv(test[0], expected)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	var actual string

	SetFromEnvValOrFallback(&actual, test, "fallback")

	if actual != expected {
		t.Errorf("Env variable not found: %s", actual)
	}

	err = os.Unsetenv(test[0])
	if err != nil {
		t.Errorf("Error unsetting env: %s", err)
	}

}

func Test_SetFromEnvValOrFallback_Fallback(t *testing.T) {
	test := []string{
		"TEST",
	}

	expected := "test"

	var actual string

	SetFromEnvValOrFallback(&actual, test, expected)

	if actual != expected {
		t.Errorf("Fallback not set: %s", actual)
	}
}

func Test_SetFromEnvValOrFallback_Negative(t *testing.T) {
	if os.Getenv("EXEC") == "1" {
		test := []string{
			"TEST",
		}

		var dst string

		SetFromEnvValOrFallback(&dst, test, "")
	}

	cmd := exec.Command(os.Args[0], "-test.run=Test_SetFromEnvValOrFallback_Negative")

	cmd.Env = append(os.Environ(), "EXEC=1")

	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}

	t.Errorf("Process ran with err %v, but we want exit status 1", err)
}
