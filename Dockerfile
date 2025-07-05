FROM golang:1.24

WORKDIR ${GOPATH}/micro-blog/
COPY . ${GOPATH}/micro-blog/

RUN go build -o /build ./cmd/micro-blog/ \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]