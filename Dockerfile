FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/dag-sync /app/cmd

FROM alpine:latest

COPY --from=builder /app/dag-sync/cmd /bin/dag-sync

CMD ["/bin/dag-sync"]