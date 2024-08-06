FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/docker-exporter ./cmd/docker-exporter


FROM alpine:3 AS runtime

WORKDIR /app

COPY --from=builder /app/docker-exporter /app/docker-exporter

ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL
ARG VERSION

EXPOSE 9100

LABEL org.opencontainers.image.created=$BUILD_DATE \
      org.opencontainers.image.authors="Julien W <cefadrom1@gmail.com>" \
      org.opencontainers.image.source=$VCS_URL \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.title="docker-exporter" \
      org.opencontainers.image.description="A prometheus exporter to expose docker metrics that are normally not available through cAdvisor" \
      org.opencontainers.image.base.name="golang" \
      org.opencontainers.image.base.version="1.22.1-alpine"

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=5s CMD wget -q -O - http://localhost:${SERVER_PORT}/health || exit 1

CMD ["/app/docker-exporter"]
