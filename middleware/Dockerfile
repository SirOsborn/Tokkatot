# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY middleware/go.mod middleware/go.sum ./
RUN go mod download

# Copy middleware source code
COPY middleware/ .

# Build the binary (statically linked, no CGO needed)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o tokkatot .

# ── Final stage — minimal runtime image ──────────────────────────
FROM alpine:3.19

WORKDIR /app

# Install CA certificates (required for HTTPS outbound calls, e.g. web push)
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Phnom_Penh /etc/localtime \
    && echo "Asia/Phnom_Penh" > /etc/timezone

# Copy the compiled binary
COPY --from=builder /app/tokkatot .

# Copy frontend and certs (mounted as volumes in production, baked in for dev)
COPY certs ./certs
COPY frontend ./frontend

# Expose app port
EXPOSE 3000

# Runtime defaults (overridden by docker-compose environment block)
ENV ENVIRONMENT=production \
    SERVER_PORT=3000 \
    SERVER_HOST=0.0.0.0 \
    DB_SSLMODE=disable

# Run as non-root user for security
RUN addgroup -S tokkatot && adduser -S tokkatot -G tokkatot
RUN chown -R tokkatot:tokkatot /app
USER tokkatot

CMD ["./tokkatot"]
