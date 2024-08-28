FROM golang:1.23-alpine AS builder

LABEL author="Mystic"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-w -s" -o /setup-custom-action-by-docker

FROM gcr.io/distroless/static

COPY --from=builder /setup-custom-action-by-docker /setup-custom-action-by-docker

ENTRYPOINT ["/setup-custom-action-by-docker"]
