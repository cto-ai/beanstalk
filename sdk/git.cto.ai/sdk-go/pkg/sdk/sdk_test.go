package sdk

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

	if actual.Client.Url != "http://127.0.0.1" {
		t.Errorf("Error unexpected client url: %v", actual.Client.Url)
	}

	if actual.Client.Port != expected {
		t.Errorf("Error unexpected client url: %v", actual.Client.Port)
	}
}

// The negative case is tested in pkg/utils/utils_test.go
