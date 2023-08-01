# Use an official Golang image with the correct version
FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
# COPY cred.env .

RUN go build -o main

# RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

FROM debian:buster-slim

RUN apt-get update && apt-get install -y postgresql-client

WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=build /app/main .
# COPY --from=build /go/main .

# Copy the cred.env file to the working directory
# COPY cred.env .

EXPOSE 9999

ENV DB_CONNECTION_STRING=postgres://jaikash12:jaikash12@localhost/ginauth?sslmode=disable
ENV IDP_URL=http://host.docker.internal:1234

CMD ["./main"]