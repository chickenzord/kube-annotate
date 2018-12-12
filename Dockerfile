FROM golang:1.11.2-alpine AS builder

WORKDIR /go/src/github.com/chickenzord/kube-annotate
RUN apk add -U --no-cache git curl wget && \
    go get -u github.com/golang/dep/cmd/dep
COPY Gopkg.* ./
RUN dep ensure -v -vendor-only
COPY . ./
RUN go build -o /bin/kube-annotate .

FROM alpine:3.8 AS runtime
RUN apk add -U --no-cache curl wget bash
COPY --from=builder /bin/kube-annotate /bin/kube-annotate
CMD ["kube-annotate"]
