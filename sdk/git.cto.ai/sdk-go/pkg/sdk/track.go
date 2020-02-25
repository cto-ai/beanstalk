package sdk

import (
	"encoding/json"
	"fmt"

	"git.cto.ai/sdk-go/pkg/request"
)

const trackServiceRequest = "TrackServiceRequest"

// TrackRequest private method generates a request to the SDK daemon from
// a TrackRequest type struct and fills in the appropriate request
// details for a track request.
// The request is subsequently executed and the error is returned.
//
// The request method and struct are exported so that an advanced user
// may choose to use the direct request as opposed to using the
// constructor/activator methods defined below.
//
// Example:
//
// See public sdk.Track() method.
type TrackRequest struct {
	Tags     []string               `json:"tags"`
	Event    string                 `json:"event"`
	MetaData map[string]interface{} `json:"details"`
	Error    string                 `json:"error"`
}

func (s *Sdk) TrackRequest(payload TrackRequest) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal request body in service %s: %s", trackServiceRequest, err)
	}

	op := &request.Operation{
		Name:       trackServiceRequest,
		HttpMethod: "POST",
		HttpPath:   "/track",
		Payload:    bytes,
	}

	req := s.NewRequest(op, nil)
	return req.Exec()
}

// Track is a constructor method that calls the TrackRequest method.
// This is the public facing API to enable developers to send data
// that they want to be tracked by cto.ai
//
// Example:
//
// s := sdk.New()
// err := s.Track("testing", []string{"sdk", "go", "tracked"}, map[string]interface{}{
//     "user": "name",
// })
// if err != nil {
//     panic(err)
// }
//
// The event, tags, and payload will be logged.
func (s *Sdk) Track(event string, tags []string, payload map[string]interface{}) error {
	r := TrackRequest{
		Event:    event,
		Tags:     tags,
		MetaData: payload,
	}

	return s.TrackRequest(r)
}
