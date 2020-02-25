package sdk

import (
	"encoding/json"
	"fmt"

	"git.cto.ai/sdk-go/pkg/request"
)

const progressbarServiceRequest = "ProgressbarServiceRequest"

// ProgressbarRequest private method generates a request to the SDK daemon from
// a ProgressbarRequest type struct and fills in the appropriate request
// details for a progressbar request.
// The request is subsequently executed and the error is returned.
//
// The request method and struct are exported so that an advanced user
// may choose to use the direct request as opposed to using the
// constructor/activator methods defined below.
//
// Example:
//
// See public sdk.ProgressbarStart(), sdk.ProgressbarAdvance,
// and sdk.ProgressbarStop methods.
type ProgressbarRequest struct {
	Length    int    `json:"length"`
	Initial   int    `json:"initial"`
	Increment int    `json:"increment"`
	Text      string `json:"text"`
}

func (s *Sdk) ProgressbarRequest(payload ProgressbarRequest, httpPath string) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal request body in service %s: %s", progressbarServiceRequest, err)
	}

	op := &request.Operation{
		Name:       progressbarServiceRequest,
		HttpMethod: "POST",
		HttpPath:   httpPath,
		Payload:    bytes,
	}

	req := s.NewRequest(op, nil)
	return req.Exec()
}

// ProgressbarStart is a constructor method that calls the ProgressbarRequest method.
// This is the public facing API to enable developers to present a progressbar
// to the output interface (i.e. terminal or slack) and will stay present until
// the progressbar advance or stop method is called.
//
// The input length is the total length of the progress bar, e.g.
// if you have 5 steps in your logic, then a unit length of 5 might be
// and appropriate length.
//
// The initial length indicates the unit length (out of total length) that is initially
// filled at the start.
//
// Example:
//
// s := sdk.New()
// err := s.ProgressbarStart("Downloading...", 5, 1)
// if err != nil {
//     panic(err)
// }
//
// Output:
// [progressbar animation with 1/5 of the bar filled here] Downloading...
func (s *Sdk) ProgressbarStart(text string, length, initial int) error {
	r := ProgressbarRequest{
		Length:   length,
		Initial:  initial,
		Text:     text,
	}

	return s.ProgressbarRequest(r, "/progress-bar/start")
}

// ProgressbarAdvance is a constructor method that calls the ProgressbarRequest method.
// This is the public facing API to enable developers to add onto a progressbar
// that is already present on the interface (i.e. terminal or slack).
//
// The increment indicates the unit length (out of total length) that will be filled.
//
// Example:
//
// ...
// [progressbar animation with 1/5 of the bar filled here] Downloading...
// err := s.ProgressbarAdvance(1)
// if err != nil {
//     panic(err)
// }
//
// Output:
// [progressbar animation with 2/5 of the bar filled here] Downloading...
func (s *Sdk) ProgressbarAdvance(increment int) error {
	r := ProgressbarRequest{
		Increment: increment,
	}

	return s.ProgressbarRequest(r, "/progress-bar/advance")
}

// ProgressbarStop is a constructor method that calls the ProgressbarRequest method.
// This is the public facing API to enable developers to complete the progressbar
// that is already present on the interface (i.e. terminal or slack).
//
// Ideally the progressbar should be completed (using sdk.ProgressbarAdvance) prior
// to stopping the progressbar.
//
// The text will change the initial text (set from the sdk.ProgressbarStart method).
//
// Example:
//
// ...
// [progressbar animation with 5/5 of the bar filled here] Downloading...
// err := s.ProgressbarStop("Done!")
// if err != nil {
//     panic(err)
// }
//
// Output:
// [progressbar animation with 5/5 of the bar filled here] Done!
func (s *Sdk) ProgressbarStop(text string) error {
	r := ProgressbarRequest{
		Text:     text,
	}

	return s.ProgressbarRequest(r, "/progress-bar/stop")
}
