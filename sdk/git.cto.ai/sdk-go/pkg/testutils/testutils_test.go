package testutils

import (
	"testing"
)

func Test_UrlPortSplitter(t *testing.T) {
	test := "http://localhost:3000"

	url, port := UrlPortSplitter(test)

	if url != "http://localhost" {
		t.Errorf("Error unexpected url: %v", url)
	}

	if port != "3000" {
		t.Errorf("Error unexpected port: %v", port)
	}
}
