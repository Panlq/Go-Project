FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR C:/Users/asus/Desktop/review/go-gin-example
COPY . C:/Users/asus/Desktop/review/go-gin-example
RUN go build .

EXPOSE 8000
ENTRYPOINT [" ./go-gin-example"]
