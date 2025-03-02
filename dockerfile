FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.21 AS base
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .

# Development stage with Air for hot reload
FROM base AS dev
USER root
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -o /usr/local/bin/core main.go
RUN chmod +x /usr/local/bin/core
RUN go install github.com/air-verse/air@latest
ENV PATH="/go/bin:${PATH}"
EXPOSE 8080
ENTRYPOINT ["air"]
CMD ["-c", ".air.toml"]

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