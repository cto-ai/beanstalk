package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.cto.ai/sdk-go/pkg/testutils"
)

func Test_New(t *testing.T) {
	url := "test"

	port := "1351"

	op := &Operation{
		Name:       "test",
		HttpMethod: "test",
		HttpPath:   "/test",
		Payload:    []byte("test"),
	}

	type Test struct {
		Test string
	}

	data := &Test{}

	actual := New(op, data, url, port)

	if actual.url != url {
		t.Errorf("Incorrect url: %v", url)
	}

	if actual.port != port {
		t.Errorf("Incorrect port: %v", port)
	}

	if actual.op.Name != op.Name {
		t.Errorf("Incorrect op.Name: %v", actual.op.Name)
	}

	if actual.op.HttpMethod != op.HttpMethod {
		t.Errorf("Incorrect op.HttpMethod: %v", actual.op.HttpMethod)
	}

	if actual.op.HttpPath != op.HttpPath {
		t.Errorf("Incorrect op.HttpPath: %v", actual.op.HttpPath)
	}

	if string(actual.op.Payload) != string(op.Payload) {
		t.Errorf("Incorrect op.Payload: %v", string(actual.op.Payload))
	}

	if actual.response != nil {
		t.Errorf("Incorrect response: %v", actual.response)
	}

	if actual.body.String() != bytes.NewBuffer(op.Payload).String() {
		t.Errorf("Incorrect body: %v", actual.body.String())
	}

	if actual.httpClient == nil {
		t.Errorf("Request client should not be nil")
	}

	if actual.data != data {
		t.Errorf("Incorrect request data: %v", actual.data)
	}
}

func Test_sendRequest_Positive(t *testing.T) {
	op := &Operation{
		Name:       "test",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte(`{"Text": "test"}`),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Content-Type"][0] != "application/json" {
			t.Errorf("Headers incorrect: %v", r.Header["Content-Type"])
		}

		if r.Method != op.HttpMethod {
			t.Errorf("Method incorrect: %v", r.Method)
		}

		if r.URL.Path != op.HttpPath {
			t.Errorf("Method incorrect: %v", r.URL.Path)
		}

		var tmp map[string]string
		err := json.NewDecoder(r.Body).Decode(&tmp)
		if err != nil {
			t.Errorf("Error in decoding response body: %s", err)
		}

		if tmp["Text"] != "test" {
			t.Errorf("Error unexpected request body: %+v", tmp)
		}

		fmt.Fprintf(w, `{"Text": "OK"}`)
	}))

	defer ts.Close()

	// to mock the data interface
	type Test struct {
		Test string
	}
	data := &Test{}

	url, port := testutils.UrlPortSplitter(ts.URL)

	// create new request
	actual := New(op, data, url, port)

	actual.sendRequest()

	var tmp map[string]string
	err := json.NewDecoder(actual.response.Body).Decode(&tmp)
	if err != nil {
		t.Errorf("Error in decoding response body: %s", err)
	}

	if tmp["Text"] != "OK" {
		t.Errorf("Error unexpected request body: %+v", tmp)
	}
}

// TODO: Mock negative sendRequest case
// Example no connection so we can test the httpClient.Do error catching

func Test_httpStatusSuccessful_Positive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"Text": "OK"}`)
	}))

	defer ts.Close()

	op := &Operation{
		Name:       "test",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte(`{"Text": "test"}`),
	}

	// to mock the data interface
	type Test struct {
		Test string
	}
	data := &Test{}

	url, port := testutils.UrlPortSplitter(ts.URL)

	// create new request
	actual := New(op, data, url, port)

	err := actual.Exec()
	if err != nil {
		t.Errorf("Error unexpected error in request: %v", err)
	}

	if actual.response.StatusCode != 200 {
		t.Errorf("Error unexpected status code: %v", actual.response.StatusCode)
	}

}

func Test_httpStatusSuccessful_Negative(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)

		fmt.Fprintf(w, `{"Text": "NOT OK"}`)
	}))

	defer ts.Close()

	op := &Operation{
		Name:       "test",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte(`{"Text": "test"}`),
	}

	url, port := testutils.UrlPortSplitter(ts.URL)

	// to mock the data interface
	type Test struct {
		Test string
	}
	data := &Test{}

	// create new request
	actual := New(op, data, url, port)

	err := actual.Exec()
	if err == nil {
		t.Errorf("Error expected an error in request: %v", err)
	}
}

func Test_unmarshalJson_Positive(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"Text": "OK"}`)
	}))

	defer ts.Close()

	url, port := testutils.UrlPortSplitter(ts.URL)

	op := &Operation{
		Name:       "test",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte(`{"Text": "OK?"}`),
	}

	type Test struct {
		Text string `name:"Text"`
	}
	data := &Test{}

	// create new request
	actual := New(op, data, url, port)

	err := actual.sendRequest()
	if err != nil {
		t.Errorf("Error unexpected error in request: %v", err)
	}

	err = actual.unmarshalJson()
	if err != nil {
		t.Errorf("Error unexpected error in unmarshalling response json: %v", err)
	}

	if data.Text != "OK" {
		t.Errorf("Error unexpected data text: %s", data.Text)
	}

}

func Test_unmarshalJson_Negative(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"Text": "Incomplete json"`)
	}))

	defer ts.Close()

	url, port := testutils.UrlPortSplitter(ts.URL)

	op := &Operation{
		Name:       "test",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte{},
	}

	type Test struct {
		Text string `name:"Text"`
	}
	data := &Test{}

	// create new request
	actual := New(op, data, url, port)

	err := actual.sendRequest()
	if err != nil {
		t.Errorf("Error unexpected error in request: %v", err)
	}

	err = actual.unmarshalJson()
	if err.Error() != "unexpected EOF" {
		t.Errorf("Error unexpected error: %v", err.Error())
	}

	if data.Text != "" {
		t.Errorf("Error unexpected deserialized data in data: %v", data)
	}
}

func Test_unmarshalUnknown_nil(t *testing.T) {
	type Test struct {
		Text string
	}

	output := &Test{}

	err := unmarshalUnknown(output, nil)
	if err.Error() != "Expected data response from sdk daemon" {
		t.Errorf("Error unexpected error: %s", err)
	}

	if output.Text != "" {
		t.Errorf("Error unexpected deserialized data in output: %v", output)
	}
}

func Test_unmarshalUnknown_nonmap(t *testing.T) {
	type Test struct {
		Text string
	}

	output := &Test{}

	test := []string{"test", "test"}

	err := unmarshalUnknown(output, test)
	if err.Error() != "Expected a map type (json) response from sdk daemon: [test test]" {
		t.Errorf("Error unexpected error: %s", err)
	}

	if output.Text != "" {
		t.Errorf("Error unexpected deserialized data in output: %v", output)
	}
}

func Test_Exec_Prompt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"Text":"test"}`)
	}))

	defer ts.Close()

	url, port := testutils.UrlPortSplitter(ts.URL)

	op := &Operation{
		Name:       "PromptServiceRequest",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte{},
	}

	type Test struct {
		Text string `name:"Text"`
	}
	data := &Test{}

	// create new request
	actual := New(op, data, url, port)
	err := actual.Exec()
	if err != nil {
		t.Errorf("Error unexpected error in request: %v", err)
	}

	if data.Text != "test" {
		t.Errorf("Error unexpected data deserialized into output: %v", data.Text)
	}
}

func Test_Exec_NonPrompt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"Text":"test"}`)
	}))

	defer ts.Close()

	url, port := testutils.UrlPortSplitter(ts.URL)

	op := &Operation{
		Name:       "notaprompt",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte{},
	}

	type Test struct {
		Text string `name:"Text"`
	}
	data := &Test{}

	// create new request
	actual := New(op, data, url, port)
	err := actual.Exec()
	if err != nil {
		t.Errorf("Error unexpected error in request: %v", err)
	}

	if data.Text != "" {
		t.Errorf("Error unexpected data deserialized into output: %v", data.Text)
	}
}

func Test_unmarshalUnknown_Negative(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"Text":true}`)
	}))

	defer ts.Close()

	url, port := testutils.UrlPortSplitter(ts.URL)

	op := &Operation{
		Name:       "PromptServiceRequest",
		HttpMethod: "POST",
		HttpPath:   "/test",
		Payload:    []byte{},
	}

	// mock a non string non float64 type
	// should cause an error
	type Test struct {
		Text bool `name:"Text"`
	}
	data := &Test{}

	// create new request
	actual := New(op, data, url, port)

	err := actual.Exec()
	if err == nil {
		t.Errorf("Error expected an error in request: %v", err)
	}

	if err.Error() != "Unexpected and unsupported value matching: true with (bool)" {
		t.Errorf("Error expected an error to catch bool case: %v", err.Error())
	}
}
