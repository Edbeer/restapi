FROM golang:1.18 as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download



FROM golang:1.18 as runner

COPY --from=builder ./ ./

RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app
ENV config=docker

EXPOSE 5000
EXPOSE 5555
EXPOSE 7070

ENTRYPOINT CompileDaemon --build="go build cmd/api/main.go" --command=./main

