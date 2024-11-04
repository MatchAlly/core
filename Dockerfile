FROM golang:1.23-alpine3.20 AS builder
WORKDIR /build
RUN apk add --no-cache build-base
COPY go.* ./
RUN go mod download -x
COPY . .
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/bin/service main.go

FROM alpine:3.20 as image-base
WORKDIR /app
COPY --from=builder /build/bin/service /usr/bin/service
ENTRYPOINT [ "service" ]
CMD [ "serve" ]
