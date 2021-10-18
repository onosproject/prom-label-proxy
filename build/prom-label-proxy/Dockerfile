FROM onosproject/golang-build:v1.0.0 as build

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile.onf go.mod go.sum /go/src/github.com/onosproject/prom-label-proxy/
COPY main.go /go/src/github.com/onosproject/prom-label-proxy/
COPY injectproxy/ /go/src/github.com/onosproject/prom-label-proxy/injectproxy/
COPY vendor/ /go/src/github.com/onosproject/prom-label-proxy/vendor/

RUN cd /go/src/github.com/onosproject/prom-label-proxy && GOFLAGS=-mod=vendor make -f Makefile.onf ${ONOS_MAKE_TARGET}

FROM alpine:3.13
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/prom-label-proxy/build/_output/prom-label-proxy /usr/local/bin/prom-label-proxy

ENTRYPOINT ["prom-label-proxy"]
