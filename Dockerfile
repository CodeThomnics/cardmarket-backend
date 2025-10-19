FROM golang:1.25.3-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM golang:1.25.3-alpine AS dev

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

# Add Go bin to PATH
ENV PATH="/root/go/bin:${PATH}"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE ${PORT}
CMD ["air", "-c", ".air.toml"]

FROM alpine:3.22.0 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
EXPOSE ${PORT}
CMD ["./main"]


