# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.* ./
RUN go mod download
COPY . .
# Icon collections and vendor JS are committed under internal/icons/collections
# and web/static/vendor. To refresh them, run scripts/fetch-assets.sh locally.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/minidash ./cmd/minidash

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/minidash /minidash
USER nonroot:nonroot
EXPOSE 8080
VOLUME ["/config"]
ENV MINIDASH_CONFIG=/config/config.yaml MINIDASH_ADDR=:8080
ENTRYPOINT ["/minidash"]
