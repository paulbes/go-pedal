FROM golang:1.11.1

ARG command
WORKDIR /go/src/github.com/paulbes/go-pedal
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a $command

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
COPY --from=0 /go/src/github.com/paulbes/go-pedal/main /pedal

ENTRYPOINT ["/pedal"]