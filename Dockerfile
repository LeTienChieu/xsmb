FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o myapp .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/myapp /myapp
CMD ["/myapp"]