FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.22 AS builder
WORKDIR /build
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /build/bin/core main.go

FROM gcr.io/distroless/static:nonroot AS runtime
WORKDIR /app
COPY --from=builder /build/bin/core /usr/bin/core

ENTRYPOINT [ "/usr/bin/core" ]
CMD [ "api" ]