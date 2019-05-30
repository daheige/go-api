FROM alpine:3.9

WORKDIR /go

RUN mkdir -p /go/logs && mkdir /go/conf

COPY ./go-api /go/go-api

EXPOSE 1338

VOLUME [ "/go/logs","/go/conf"]

CMD [ "/go/go-api","-log_dir=/go/logs","-config_dir=/go/conf"]