########################################
### base
FROM golang:1.22.0-alpine3.19 as base

RUN apk update --no-cache \
    && apk add --no-cache git \
    && apk add make git golangci-lint

########################################
### builder
FROM base as builder

RUN mkdir -p /usr/local/go/bin

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

########################################
### runner
FROM base as runner

COPY --from=builder /usr/local/go /usr/local/go
COPY --from=builder /go /go

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

WORKDIR /app