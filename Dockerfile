FROM alpine:3.9

WORKDIR /go

#tsinghua alpine source
#add curl
RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.9/main/" > /etc/apk/repositories \
        && apk update \
        && apk upgrade \
        && apk add --no-cache bash \
        bash-doc \
        bash-completion curl \
        && rm -rf /var/cache/apk/* \
        && mkdir -p /go/logs && mkdir /go/conf

COPY ./go-api /go/go-api

EXPOSE 1338

VOLUME [ "/go/logs","/go/conf"]

CMD [ "/go/go-api","-log_dir=/go/logs","-config_dir=/go/conf"]