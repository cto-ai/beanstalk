package sdk

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.cto.ai/sdk-go/pkg/testutils"
)

func Test_ProgressbarRequest_ProgressbarStart(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/progress-bar/start" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp ProgressbarRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Text != "start" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Length != 2 {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Initial != 1 {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Increment != 0 {
			t.Errorf("Error unexpected request body: %v", tmp)
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
	s.ProgressbarStart("start", 2, 1)
}

func Test_ProgressbarRequest_ProgressbarAdvance(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/progress-bar/advance" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp ProgressbarRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Text != "" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Length != 0 {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Initial != 0 {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Increment != 1 {
			t.Errorf("Error unexpected request body: %v", tmp)
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
	s.ProgressbarAdvance(1)
}

func Test_ProgressbarRequest_ProgressbarStop(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/progress-bar/stop" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp ProgressbarRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Text != "Done" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Length != 0 {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Initial != 0 {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Increment != 0 {
			t.Errorf("Error unexpected request body: %v", tmp)
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
	s.ProgressbarStop("Done")
}
