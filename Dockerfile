FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /build
RUN apk add --no-cache build-base
COPY ["go.mod", "go.sum", "./"]
RUN go mod download -x
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/bin/service main.go

FROM alpine:3.19 as image-base
WORKDIR /app
COPY --from=builder /build/bin/service /usr/bin/service
ENTRYPOINT [ "service" ]
CMD [ "serve" ]
