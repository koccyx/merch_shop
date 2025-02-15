FROM golang:1.23.3

WORKDIR ${GOPATH}/avito-assignment/
COPY . ${GOPATH}/avito-assignment/

ENV CONFIG_PATH=./config/dev.yaml

RUN go build -o /build ./cmd/avito_assignment \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]