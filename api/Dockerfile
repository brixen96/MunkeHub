# Build stage
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /munkehub-api .

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=build /munkehub-api .
# Generic setup package for the website Download button + the helper
# exe for self-updates (build them first: scripts/build-download.ps1)
COPY download/ ./download/
COPY helper/ ./helper/
# Cache dirs (icons + 3D model assets, lazy-fetched and persisted on disk)
RUN mkdir -p /app/icons /app/mv /app/charmodels
ENV MUNKEHUB_PORT=8080 \
    MUNKEHUB_DB=/data/munkehub.db
VOLUME /data
EXPOSE 8080
CMD ["./munkehub-api"]