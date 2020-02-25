package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

var (
	unmarshalOps = map[string]bool{"PromptServiceRequest": true}
)

// An operation is an API operation.
type Operation struct {
	Name       string
	HttpMethod string
	HttpPath   string
	Payload    []byte
}

// A request is the operation http request.
// Request contains the request parameters needed
// for the http request, as well as, memory
// allocated to hold the successive response body
// and errors.
type Request struct {
	url        string
	port       string
	op         *Operation
	response   *http.Response
	body       *bytes.Buffer
	httpClient *http.Client

	// data is a void pointer to the output struct
	// of a specific operation. This is set at
	// the service/API level, e.g. sdk/ux/prompt.
	data interface{}
}

// New returns a new Request pointer for an API operation.
func New(op *Operation, data interface{}, url, port string) *Request {
	return &Request{
		url:        url,
		port:       port,
		op:         op,
		response:   nil,
		body:       bytes.NewBuffer(op.Payload),
		httpClient: &http.Client{},
		data:       data,
	}
}

// Exec sends the API request and handles any errors returned.
// Successful responses will be validated and handled according
// the operation type.
func (r *Request) Exec() error {
	err := r.sendRequest()
	if err != nil {
		return err
	}

	err = r.httpStatusSuccessful()
	if err != nil {
		return err
	}

	// Only prompt service endpoints currently require
	// the need to unmarshal the body from a response
	// TODO: to more generically handle which requests
	// need specific post response transformations
	if unmarshalOps[r.op.Name] {
		err = r.unmarshalJson()
		if err != nil {
			return err
		}
	}

	return nil
}

// sendRequest sends the http request to the request url with
// the associated parameters.
func (r *Request) sendRequest() error {
	req, err := http.NewRequest(
		r.op.HttpMethod,
		fmt.Sprintf("%s:%s%s", r.url, r.port, r.op.HttpPath),
		bytes.NewBuffer(r.op.Payload),
	)
	if err != nil {
		return fmt.Errorf("Error in constructing request: %s", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error in request: %s", err)
	}

	r.response = resp

	return nil
}

// httpStatusSuccessful validates the http response success.
func (r *Request) httpStatusSuccessful() error {
	if r.response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Error in status code: %d", r.response.StatusCode)
	}

	return nil
}

// unmarshalJson deserializes the response body into a void pointer
// and passes the pointer to unmarshalUnknown.
//
// This is so that different request body response outputs
// can be generically deserialized into their respective types.
func (r *Request) unmarshalJson() error {
	var data interface{}

	err := json.NewDecoder(r.response.Body).Decode(&data)
	if err != nil {
		return err
	}

	err = unmarshalUnknown(r.data, data)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalUnknown deserializes a void pointer referencing
// the response body into the appropriate output struct type.
// This is achieved by using the reflect module to check the
// metadata of struct definitions and cast accordingly.
func unmarshalUnknown(output, data interface{}) error {
	if data == nil {
		return fmt.Errorf("Expected data response from sdk daemon")
	}

	mapData, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Expected a map type (json) response from sdk daemon: %v", data)
	}

	value := reflect.ValueOf(output)
	t := value.Type()

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		name := field.Tag.Get("name")

		member := value.FieldByIndex(field.Index)

		switch d := mapData[name].(type) {
		case string:
			member.Set(reflect.ValueOf(d))
		default:
			return fmt.Errorf("Unexpected and unsupported value matching: %v with (%s)", mapData[name], member.Type())
		}
	}

	return nil
}
