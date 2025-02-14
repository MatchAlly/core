FROM --platform=$BUILDPLATFORM golang:1.23-alpine3.20 AS base
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .

# Development stage
FROM base AS dev
RUN apk add --no-cache curl wget tzdata
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -o /usr/local/bin/core main.go
RUN chmod +x /usr/local/bin/core
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/core"]
CMD ["api"]

# Production stage
FROM base AS builder
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -o /usr/local/bin/core main.go

FROM gcr.io/distroless/static:nonroot AS prod
COPY --from=builder /usr/local/bin/core /usr/local/bin/core
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/core"]
CMD ["api"]