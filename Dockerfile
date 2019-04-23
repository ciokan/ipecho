FROM golang:1.12.1-alpine3.9 as builder

ENV GO111MODULE=on

RUN set -x \
	&& apk update \
	&& apk add --no-cache git bash curl ca-certificates

RUN mkdir -p $GOPATH/src/gitlab.com/ciokan/ipecho && \
	curl -Lo /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.1.3/dumb-init_1.1.3_amd64 && \
	chmod +x /usr/local/bin/dumb-init

WORKDIR $GOPATH/src/gitlab.com/ciokan/ipecho
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . $GOPATH/src/gitlab.com/ciokan/ipecho

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/ipecho

FROM alpine:3.9
COPY --from=builder /go/bin/ipecho /usr/local/bin/ipecho
COPY --from=builder /usr/local/bin/dumb-init /usr/local/bin/dumb-init
RUN chmod +x /usr/local/bin/dumb-init

EXPOSE 8080
ENTRYPOINT ["dumb-init", "--"]
CMD /usr/local/bin/ipecho