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
ENV MUNKEHUB_PORT=8080 \
    MUNKEHUB_DB=/data/munkehub.db
VOLUME /data
EXPOSE 8080
CMD ["./munkehub-api"]