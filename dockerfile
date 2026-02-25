FROM golang:1.26.0-alpine AS builder
RUN apk add --no-cache git
WORKDIR /Dialyn
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app .

FROM scratch
COPY --from=builder /app /app
ENV RABBIT_PORT=380 RABBIT_URL=amqp://guest:guest@localhost OWN_PORT=8250
EXPOSE 8250 380
ENTRYPOINT ["/app"]