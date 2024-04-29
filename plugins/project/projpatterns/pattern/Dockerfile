FROM --platform=$BUILDPLATFORM golang as builder

WORKDIR /app

ARG TARGETOS TARGETARCH

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /deploy/server/service ./cmd/service/main.go && \
    cp -r config /deploy/server/config

FROM alpine

LABEL MATRESHKA_CONFIG_ENABLED=true

WORKDIR /app
COPY --from=builder /deploy/server/ .

ENTRYPOINT ["./service"]