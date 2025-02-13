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
    -o /app/bin/service main.go
EXPOSE 8080
ENTRYPOINT ["core"]
CMD ["serve"]

# Production stage
FROM base AS builder
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -o /app/bin/service main.go
FROM gcr.io/distroless/static:nonroot AS prod
COPY --from=builder /app/bin/service /usr/bin/core
COPY --from=gcr.io/distroless/static:healthcheck /healthcheck /healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
    CMD ["/healthcheck", "http://localhost:8080/health"]
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["core"]
CMD ["serve"]