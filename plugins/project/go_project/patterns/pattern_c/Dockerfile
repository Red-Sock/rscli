FROM --platform=$BUILDPLATFORM golang as builder

WORKDIR /app

RUN --mount=target=. \
        --mount=type=cache,target=/root/.cache/go-build \
        --mount=type=cache,target=/go/pkg \
        GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 \
    go build -o /deploy/server/service ./cmd/service/main.go && \
    cp -r config /deploy/server/config &&\
     if [ -d "./migrations" ];  then \
      cp -r ./migrations /deploy/server/migrations;  \
     fi
FROM alpine

LABEL MATRESHKA_CONFIG_ENABLED=true

WORKDIR /app

COPY --from=builder /deploy/server/ .

EXPOSE 80

ENTRYPOINT ["./service"]