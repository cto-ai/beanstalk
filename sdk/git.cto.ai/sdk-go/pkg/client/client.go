package client

import (
	"git.cto.ai/sdk-go/pkg/request"
	"git.cto.ai/sdk-go/pkg/utils"
)

// Slices of strings are used in the case when multiple
// env variable names are set for the same desired variable.
var (
	sdkPort = []string{
		"SDK_SPEAK_PORT",
	}
)

// A Client holds the data needed to connect to the
// external API.
type Client struct {
	Port string
	Url  string
}

// New returns a new Client pointer used for the connection
// the external API service.
func New() *Client {
	clnt := &Client{}

	utils.SetFromEnvValOrFallback(&clnt.Port, sdkPort, "")

	// the local sdk url is set to localhost
	utils.SetFromEnvValOrFallback(&clnt.Url, []string{}, "http://127.0.0.1")

	return clnt
}

// NewRequest returns a new Request pointer used for the API request.
// This passes the Sdk operation and data pointers, as well as,
// the Client url and port to the Request to be built.
func (c *Client) NewRequest(op *request.Operation, data interface{}) *request.Request {
	return request.New(op, data, c.Url, c.Port)
}
