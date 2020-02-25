// -build ignore

package testutils

import (
	"strings"
)

// utility function for httptest server mocks
// url and port need to be separate attributes in the request struct
//
// Takes input serverUrl in the form:
// http://url.com:9000
func UrlPortSplitter(serverUrl string) (string, string) {
	split := strings.Split(serverUrl, ":")

	base := strings.Join(split[:len(split)-1], ":")
	port := split[len(split)-1]

	return base, port
}
