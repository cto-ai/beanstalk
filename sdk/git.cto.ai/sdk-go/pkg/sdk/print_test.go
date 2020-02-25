package sdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.cto.ai/sdk-go/pkg/testutils"
)

func Test_PrintRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/print" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PrintRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Text != "test" {
			t.Errorf("Error unexpected request body: %+v", tmp)
		}
	}))

	defer ts.Close()

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err := os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	s.Print("test")
}
