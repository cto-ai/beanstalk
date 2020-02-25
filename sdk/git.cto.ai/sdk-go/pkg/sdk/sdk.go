package sdk

import (
	"git.cto.ai/sdk-go/pkg/client"
	"git.cto.ai/sdk-go/pkg/request"
)

// Sdk is a wrapper for a Client. This is used
// for naming purposes as Client is a better
// name for its purpose of holding connection
// related data. Sdk is a better name for the public
// facing API.
type Sdk struct {
	*client.Client
}

// New returns a new Sdk pointer.
func New() *Sdk {
	return &Sdk{
		Client: client.New(),
	}
}

// NewRequest returns a new Request pointer.
// This passes the operation and data pointers to the client.
func (s *Sdk) NewRequest(op *request.Operation, data interface{}) *request.Request {
	return s.Client.NewRequest(op, data)
}
