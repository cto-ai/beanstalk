# sdk-go

The Ops SDK for Go

# Usage

## Example sdk usage

SDK definitions can be found in the `pkg/sdk` directory. Documentation and example usage can be found in their respective source code.

```go
package main

import (
	"fmt"
	"git.cto.ai/sdk-go/pkg/sdk"
	"log"
	"time"
)

func main() {
	// instantiate a new sdk client
	s := sdk.New()

	// printing text to interface
	err := s.Print("starting e2e test of Go SDK")
	if err != nil {
		log.Fatal(err)
	}

	// prompt user to confirm
	// captures user stdin and converts to a bool
	output, err := s.PromptConfirm("cofirmation", "confirm?", "C", true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("You confirmed with a %v\n", output)

	// displays a progress bar to the interface
	err = s.ProgressbarStart("Finishing...", 5, 1)
	if err != nil {
		log.Fatal(err)
	}

	// do some work
	time.Sleep(2 * time.Second)

	err = s.ProgressbarAdvance(4)
	if err != nil {
		log.Fatal(err)
	}

	err = s.ProgressbarStop("Done!")
	if err != nil {
		log.Fatal(err)
	}
}
```

## Running inside an op

This uses the example code above.

NOTE: currently not publically available so using `go get` or `dep init` will not work properly unless you have access to the repo.

The current workaround is to clone the current repo and place inside the `$GOPATH` of the container.

Below is an example Dockerfile:

This assumes you have the `sdk-go` in the root of your op.

```docker
############################
# Build container
############################
FROM golang:1.13.6-alpine AS build

WORKDIR /ops

ADD . .

RUN mkdir -p $GOPATH/src/git.cto.ai

RUN mv ./sdk-go $GOPATH/src/git.cto.ai/.

RUN go build -ldflags="-s -w" -o main

############################
# Final container
############################
FROM registry.cto.ai/official_images/base:latest

COPY --from=build /ops/main /ops/main
```

Corresponding ops.yml run command:
```yaml
run: /ops/main
```
