############################
# Build container
############################
FROM golang:1.13.4 AS build

WORKDIR $GOPATH/src/git.cto.ai/provision/
RUN go get -u github.com/jmespath/go-jmespath
RUN go get -u github.com/aws/aws-sdk-go/...
RUN go get -u golang.org/x/oauth2
RUN go get -u github.com/google/go-github/github
RUN mkdir -p $GOPATH/src/git.cto.ai
RUN mkdir -p $GOPATH/src/github.com/aws/
RUN mkdir -p $GOPATH/src/github.com/jmespath/
ADD ./sdk ./sdk
RUN mv ./sdk/git.cto.ai/sdk-go $GOPATH/src/git.cto.ai/.
ADD . .
RUN go build -o main

############################
# Final container
############################
FROM registry.cto.ai/official_images/base:latest
RUN apt-get update -y && apt-get install -y -qq zip curl
COPY --from=build /go/src/git.cto.ai/provision/main /bin/.

