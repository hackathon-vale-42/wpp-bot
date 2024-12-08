FROM golang:latest

WORKDIR /app

COPY ["./go.mod", "./go.sum", "Makefile", "./"]

COPY ["cmd", "cmd/"]

COPY ["api", "api/"]

RUN go mod download

EXPOSE 8000

RUN make build

ENTRYPOINT ["./bin/wpp-bot"]
