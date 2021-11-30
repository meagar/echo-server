FROM golang:1.17-alpine3.13 as builder

WORKDIR /build
COPY echo.go go.mod go.sum index.gohtml ./

RUN GOOS=linux GOARCH=amd64 go build -o echo-server echo.go

FROM alpine

COPY --from=builder /build/echo-server /echo-server

CMD ["/echo-server", "-status=/status", "3002"]
