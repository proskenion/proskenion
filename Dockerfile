FROM alpine:3.8

ENTRYPOINT ["/proskenion/bin/linux/proskenion"]

WORKDIR /proskenion
COPY ./bin ./bin
COPY ./example ./example
COPY ./database ./database