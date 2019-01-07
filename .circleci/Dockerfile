FROM alpine:3.8
RUN apk add --no-cache build-base curl automake autoconf libtool git zlib-dev

ENV GRPC_VERSION=1.13.0 \
    PROTOBUF_VERSION=3.6.0.1 \
    PROTOC_GEN_DOC_VERSION=1.1.0 \
    OUTDIR=/out
RUN mkdir -p /protobuf && \
    curl -L https://github.com/google/protobuf/archive/v${PROTOBUF_VERSION}.tar.gz | tar xvz --strip-components=1 -C /protobuf
RUN git clone --depth 1 --recursive -b v${GRPC_VERSION} https://github.com/grpc/grpc.git /grpc && \
    rm -rf grpc/third_party/protobuf && \
    ln -s /protobuf /grpc/third_party/protobuf
RUN cd /protobuf && \
    autoreconf -f -i -Wall,no-obsolete && \
    ./configure --prefix=/usr --enable-static=no && \
    make -j2 && make install
RUN cd grpc && \
    make -j2 plugins
RUN cd /protobuf && \
    make install DESTDIR=${OUTDIR}
RUN cd /grpc && \
    make install-plugins prefix=${OUTDIR}/usr
RUN find ${OUTDIR} -name "*.a" -delete -or -name "*.la" -delete

RUN apk add --no-cache go
ENV GOPATH=/go \
    PATH=/go/bin/:$PATH
RUN go get -u -v -ldflags '-w -s' \
        github.com/Masterminds/glide \
        github.com/satellitex/protobuf/protoc-gen-go \
        github.com/gogo/protobuf/protoc-gen-gofast \
        github.com/gogo/protobuf/protoc-gen-gogo \
        github.com/gogo/protobuf/protoc-gen-gogofast \
        github.com/gogo/protobuf/protoc-gen-gogofaster \
        github.com/gogo/protobuf/protoc-gen-gogoslick \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
        github.com/johanbrandhorst/protobuf/protoc-gen-gopherjs \
        github.com/ckaznocha/protoc-gen-lint \
        && install -c ${GOPATH}/bin/protoc-gen* ${OUTDIR}/usr/bin/

RUN mkdir -p ${GOPATH}/src/github.com/pseudomuto/protoc-gen-doc && \
    curl -L https://github.com/pseudomuto/protoc-gen-doc/archive/v${PROTOC_GEN_DOC_VERSION}.tar.gz | tar xvz --strip 1 -C ${GOPATH}/src/github.com/pseudomuto/protoc-gen-doc
RUN cd ${GOPATH}/src/github.com/pseudomuto/protoc-gen-doc && \
    make build && \
    install -c ${GOPATH}/src/github.com/pseudomuto/protoc-gen-doc/protoc-gen-doc ${OUTDIR}/usr/bin/

WORKDIR /go