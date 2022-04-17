FROM golang:alpine As builder

ENV GOPROXY "https://goproxy.cn,direct"

ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /app
ADD . .

RUN set -eux; \
    dpkgArch=dpkg --print-architecture | awk -F- '{ print $NF }'; \
    buildtime="$(date -u '+%Y-%m-%d %I:%M:%S%p')"; \
    goversion="$(go version)"; \
    flags="-s -w -extldflags '-static' \
    -X 'main.buildtime=$buildtime' \
    -X 'main.goversion=$goversion'"; \
    CGO_ENABLED=0 GOOS=${GOOS} GOARCH=$dpkgArch go build -ldflags "$flags" -o hello main.go


FROM alpine

WORKDIR /app
COPY --from=builder /app/hello .
COPY ./entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

CMD ["./hello", "-v"]