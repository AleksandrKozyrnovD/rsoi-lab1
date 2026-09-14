FROM golang:1.27.0-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum* ./
RUN go mod download

COPY src/ ./
RUN go mod tidy && go build -o personservice .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder --chmod=755 /app/personservice /personservice
COPY --from=builder /app/config.yml /app/config.yml


EXPOSE 8080

ENTRYPOINT ["/personservice"]
