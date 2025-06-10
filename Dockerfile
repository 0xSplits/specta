FROM golang:1.24 AS build
WORKDIR /build

ARG SHA
ARG TAG

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w \
        -X 'github.com/0xSplits/specta/pkg/runtime.sha=${SHA}' \
        -X 'github.com/0xSplits/specta/pkg/runtime.tag=${TAG}'" \
    -a \
    -o specta main.go



FROM gcr.io/distroless/static:nonroot
WORKDIR /image

COPY .env .env
COPY --from=build /build/specta .
USER 65532:65532

ENV SPECTA_HTTP_HOST="0.0.0.0"
ENV SPECTA_LOG_LEVEL="info"

ENTRYPOINT ["/image/specta"]
