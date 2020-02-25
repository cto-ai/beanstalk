package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"git.cto.ai/sdk-go/pkg/testutils"
)

func Test_PromptRequest_PromptInput(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "input" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "type test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Default != "error" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "I" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.AllowEmpty != true {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": "test"}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptInput("test", "type test", "error", "I", true)
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != "test" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptNumber(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "number" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "type 2" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "N" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Default != float64(0) {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Min != float64(1) {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Max != float64(3) {
			t.Errorf("Error unexpected request body: %v", tmp)
		}
		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": 2}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptNumber("test", "type 2", "N", 0, 1, 3)
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != 2 {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptSecret(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "secret" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "what is secret" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "S" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": "secret"}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptSecret("test", "what is secret", "S")
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != "secret" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptPassword(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "password" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "what is password" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "P" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Confirm != true {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": "password"}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptPassword("test", "what is password", "P", true)
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != "password" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptConfirm(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "confirm" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "confirm?" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "C" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Default != true {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": false}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptConfirm("test", "confirm?", "C", true)
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != false {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptList(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	choices := []string{"aws", "gcd"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "list" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "choose" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Choices[0] != choices[0] {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Choices[1] != choices[1] {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Default != "aws" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "L" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": "gcd"}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptList(choices, "test", "choose", "aws", "L")
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != "gcd" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptAutocomplete(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	choices := []string{"aws", "gcd"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "autocomplete" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "choose" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Choices[0] != choices[0] {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Choices[1] != choices[1] {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Default != "aws" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "A" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": "gcd"}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptAutocomplete(choices, "test", "choose", "aws", "A")
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != "gcd" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptCheckbox(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	choices := []string{"aws", "gcd", "azure"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "checkbox" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "choose" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "C" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Choices[0] != choices[0] {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Choices[1] != choices[1] {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Choices[2] != choices[2] {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": ["gcd", "azure"]}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptCheckbox(choices, "test", "choose", "C")
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output[0] != "gcd" {
		t.Errorf("Error unexpected output: %v", output[0])
	}

	if output[1] != "azure" {
		t.Errorf("Error unexpected output: %v", output[1])
	}
}

func Test_PromptRequest_PromptEditor(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "editor" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "edit" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "E" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Default != "default" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": "test"}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	s := New()
	output, err := s.PromptEditor("test", "edit", "default", "E")
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	if output != "test" {
		t.Errorf("Error unexpected output: %v", output)
	}
}

func Test_PromptRequest_PromptDatetime(t *testing.T) {
	expectedResponse := `{"replyFilename": "/tmp/response-mocktest"}`

	expectedRequest := "2006-01-02T15:04:05Z"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != "POST" {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != "/prompt" {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp PromptRequest
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp.Name != "test" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptType != "datetime" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.PromptMsg != "what date" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Variant != "datetime" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Flag != "D" {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Min != expectedRequest {
			t.Errorf("Error unexpected request body: %v", tmp.Min)
		}

		if tmp.Max != expectedRequest {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		if tmp.Default != expectedRequest {
			t.Errorf("Error unexpected request body: %v", tmp)
		}

		fmt.Fprintf(w, expectedResponse)
	}))

	defer ts.Close()

	// write a fake file
	err := ioutil.WriteFile("/tmp/response-mocktest", []byte(`{"test": "2006-01-02 15:04:05"}`), 0777)

	_, port := testutils.UrlPortSplitter(ts.URL)

	test := []string{
		"SDK_SPEAK_PORT",
	}

	err = os.Setenv(test[0], port)
	if err != nil {
		t.Errorf("Error setting test env variable: %s", err)
	}

	input, err := time.Parse(time.RFC3339, expectedRequest)
	if err != nil {
		t.Errorf("Error parsing expected time: %v", err)
	}

	s := New()
	output, err := s.PromptDatetime("test", "what date", "datetime", "D", input, input, input)
	if err != nil {
		t.Errorf("Error in prompt request: %v", err)
	}

	const format = "2006-01-02 15:04:05"

	expected, err := time.Parse(format, "2006-01-02 15:04:05")
	if err != nil {
		t.Errorf("Error parsing expected time: %v", err)
	}

	if output != expected {
		t.Errorf("Error unexpected output: %v", output)
	}
}
