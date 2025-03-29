FROM golang:1.24.1-bookworm

WORKDIR /go/src

COPY go* ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air", "-c", ".air.toml"]