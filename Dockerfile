FROM golang:1.11.2-alpine AS builder

WORKDIR /go/src/github.com/chickenzord/kube-annotate
RUN apk add -U --no-cache git curl wget make && \
    curl -L -s https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 -o $GOPATH/bin/dep && \
    chmod +x $GOPATH/bin/dep
COPY Gopkg.* Makefile ./
RUN make deps
COPY . ./
RUN BUILD_OUTPUT=/bin/kube-annotate make build


FROM alpine:3.8 AS runtime

RUN apk add -U --no-cache curl wget bash
COPY --from=builder /bin/kube-annotate /bin/kube-annotate
CMD ["kube-annotate"]
