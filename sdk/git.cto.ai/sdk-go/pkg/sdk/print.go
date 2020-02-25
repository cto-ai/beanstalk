package sdk

import (
	"encoding/json"
	"fmt"

	"git.cto.ai/sdk-go/pkg/request"
)

const printServiceRequest = "PrintServiceRequest"

// PrintRequest private method generates a request to the SDK daemon from
// a PrintRequest type struct and fills in the appropriate request
// details for a print request.
// The request is subsequently executed and the error is returned.
//
// The request method and struct are exported so that an advanced user
// may choose to use the direct request as opposed to using the
// constructor/activator methods defined below.
//
// Example:
//
// See public sdk.Print() method.
type PrintRequest struct {
	Text string `json:"text"`
}

func (s *Sdk) PrintRequest(payload PrintRequest) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal request body in service %s: %s", printServiceRequest, err)
	}

	op := &request.Operation{
		Name:       printServiceRequest,
		HttpMethod: "POST",
		HttpPath:   "/print",
		Payload:    bytes,
	}

	req := s.NewRequest(op, nil)
	return req.Exec()
}

// Print is a constructor method that calls the PrintRequest method.
// This is the public facing API to enable developers to print text
// to the output interface (i.e. terminal/slack).
//
// Example:
//
// s := sdk.New()
// err := s.Print("testing")
// if err != nil {
//     panic(err)
// }
//
// Output:
//
// testing
func (s *Sdk) Print(a ...interface{}) error {
	r := PrintRequest{
		Text: fmt.Sprint(a...),
	}

	return s.PrintRequest(r)
}
