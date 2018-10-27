FROM golang:1.11.1

WORKDIR /go/src/github.com/paulbes/go-pedal
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -a cmd/pedal/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
COPY --from=0 /go/src/github.com/paulbes/go-pedal/main /pedal

ENTRYPOINT ["/pedal"]