FROM golang:1.22-alpine AS builder

WORKDIR /builder

ADD go.mod go.mod
ADD go.sum go.sum

RUN go mod download
ADD ./cmd ./cmd
ADD ./pkg ./pkg

RUN go build -ldflags '-extldflags "-static" -w -s' -o smart-app-builder ./cmd/builder

FROM mgoltzsche/podman:rootless

WORKDIR /podman
COPY --from=builder /builder/smart-app-builder /usr/local/bin/smart-app-builder
