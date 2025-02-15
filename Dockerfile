FROM golang:1.23.1

# Укажите рабочую директорию
WORKDIR /app
COPY / .

RUN go mod tidy

RUN go build -o /build ./cmd \
    && go clean -cache -modcache


EXPOSE 8080

CMD ["/build"]