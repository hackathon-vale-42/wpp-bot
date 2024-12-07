FROM golang:latest

WORKDIR /app

COPY . ./

RUN go mod download

EXPOSE 8000

ENTRYPOINT ["./scripts/run.sh"]
