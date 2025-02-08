# Build stage
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git tzdata  # ติดตั้ง tzdata ใน build stage
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./
RUN go build -o /go/bin/app -v ./

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata  # ติดตั้ง tzdata ใน final stage
COPY --from=builder /go/bin/app /app
ENTRYPOINT ["/app"]

LABEL Name=project Version=0.0.1
EXPOSE 3001
