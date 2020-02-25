package sdk

import (
	"encoding/json"
	"fmt"

	"git.cto.ai/sdk-go/pkg/request"
)

const spinnerServiceRequest = "SpinnerServiceRequest"

// SpinnerRequest private method generates a request to the SDK daemon from
// a SpinnerRequest type struct and fills in the appropriate request
// details for a spinner request.
// The request is subsequently executed and the error is returned.
//
// The request method and struct are exported so that an advanced user
// may choose to use the direct request as opposed to using the
// constructor/activator methods defined below.
//
// Example:
//
// See public sdk.SpinnerStart() and sdk.SpinnerStop methods.
type SpinnerRequest struct {
	Text     string `json:"text"`
}

func (s *Sdk) SpinnerRequest(payload SpinnerRequest, httpPath string) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal request body in service %s: %s", spinnerServiceRequest, err)
	}

	op := &request.Operation{
		Name:       spinnerServiceRequest,
		HttpMethod: "POST",
		HttpPath:   httpPath,
		Payload:    bytes,
	}

	req := s.NewRequest(op, nil)
	return req.Exec()
}

// SpinnerStart is a constructor method that calls the SpinnerRequest method.
// This is the public facing API to enable developers to present a spinner
// to the output interface (i.e. terminal or slack) and continues to spin
// until the SpinnerStop method is called.
//
// Example:
//
// s := sdk.New()
// err := s.SpinnerStart("Starting process...")
// if err != nil {
//     panic(err)
// }
//
// Output:
// [spinner emoji w/ spinner animation here] Starting process...
func (s *Sdk) SpinnerStart(text string) error {
	r := SpinnerRequest{
		Text:     text,
	}

	return s.SpinnerRequest(r, "/start-spinner")
}

// SpinnerStop is a constructor method that calls the SpinnerRequest method.
// This is the public facing API to enable developers to stop a spinner
// that has been previously called on the interface (i.e. terminal or slack).
//
// Example:
//
// ... //previous spinner started here
//
// err := s.SpinnerStop("Done!")
// if err != nil {
//     panic(err)
// }
//
// Output:
// [spinner completed completed here] Done!
func (s *Sdk) SpinnerStop(text string) error {
	r := SpinnerRequest{
		Text:     text,
	}

	return s.SpinnerRequest(r, "/stop-spinner")
}
