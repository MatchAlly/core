FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.20 AS builder
WORKDIR /build
RUN apk add --no-cache build-base
COPY ["go.mod", "go.sum", "./"]
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/bin/service main.go

FROM --platform=$BUILDPLATFORM alpine:3.20 as image-base
WORKDIR /app
COPY --from=builder /build/bin/service /usr/bin/core
ENTRYPOINT [ "core" ]
CMD [ "serve" ]
