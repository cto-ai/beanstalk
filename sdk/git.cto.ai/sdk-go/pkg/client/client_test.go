package client

import (
	"os"
	"testing"
)

func Test_New_Positive(t *testing.T) {

	test := []string{"SDK_SPEAK_PORT"}

	expected := "50001"

	err := os.Setenv(test[0], expected)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	actual := New()

	if actual.Port != expected {
		t.Errorf("Error unexpected port: %v", actual.Port)
	}

	if actual.Url != "http://127.0.0.1" {
		t.Errorf("Error unexpected port: %v", actual.Url)
	}
}

// the negative case of New() is tested in pkg/utls/utils_test.go
