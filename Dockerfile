# Stage 1: Build the application
FROM golang:1.21 AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

# Stage 2: Create the final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=build /app/main .

CMD ["./main"]
