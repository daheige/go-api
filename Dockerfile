FROM docker.io/alpine

WORKDIR /go

#ENV TZ=Asia/Shanghai

#RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir /go/logs

COPY ./go-api /go/go-api
#可以不用把配置文件复制到容器中
COPY ./app.yaml /go/

EXPOSE 1338

VOLUME [ "/go/logs"]

CMD [ "/go/go-api","-log_dir=/go/logs","-config_dir=/go"]