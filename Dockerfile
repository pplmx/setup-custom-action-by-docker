FROM golang:1.23-alpine AS builder

LABEL author="Mystic"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -trimpath -ldflags="-w -s" -o /setup-my-action

FROM gcr.io/distroless/static

COPY --from=builder /setup-my-action /setup-my-action

ENTRYPOINT ["/setup-my-action"]
