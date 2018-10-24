FROM alpine:3.8

CMD ["/proskenion/bin/linux/proskenion"]

WORKDIR /proskenion
COPY ./bin ./bin
COPY ./config ./config
COPY ./database ./database