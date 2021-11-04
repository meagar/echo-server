FROM golang:1.15-alpine3.13
RUN apk add git
RUN go get -u github.com/cosmtrek/air
WORKDIR /run
COPY echo.go go.mod go.sum .air.toml ./
ENTRYPOINT [ "air", "-c", ".air.toml" ]
