FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o checklist-apps

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/checklist-apps /checklist-apps

COPY .env /app/.env
COPY configs /app/configs
COPY templates /app/templates

ENV PORT=80

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD [ "wget", "-q", "http://localhost:80/diag/health", "-O", "/dev/null" ] || exit 1

CMD ["/checklist-apps"]
