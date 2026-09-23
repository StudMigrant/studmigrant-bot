FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/studmigrant-bot ./cmd/bot

FROM debian:12-slim AS certs
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl && \
    curl -o /usr/local/share/ca-certificates/russian_trusted_root_ca.crt \
    "https://gu-st.ru/content/lending/russian_trusted_root_ca_pem.crt" && \
    curl -o /usr/local/share/ca-certificates/russian_trusted_sub_ca.crt \
    "https://gu-st.ru/content/lending/russian_trusted_sub_ca_pem.crt" && \
    update-ca-certificates

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/studmigrant-bot /app/studmigrant-bot

EXPOSE 8080

ENTRYPOINT ["/app/studmigrant-bot"]
