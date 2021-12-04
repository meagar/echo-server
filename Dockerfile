FROM golang:1.17-alpine3.13 as builder

ARG GITHUB_SHA
ARG GITHUB_REF

WORKDIR /build
COPY echo.go go.mod go.sum index.gohtml ./

RUN GOOS=linux GOARCH=amd64 go build -o echo-server echo.go

FROM alpine

COPY --from=builder /build/echo-server /echo-server

ARG GITHUB_SHA
ARG GITHUB_REF
ENV SHA=$GITHUB_SHA
ENV REF=$GITHUB_REF

CMD ["/echo-server", "-status=/status", "3002"]
