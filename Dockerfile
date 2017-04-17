FROM alpine:3.4
MAINTAINER Barnett <zengqg@yiyun.pro>

RUN mkdir /alert
COPY ./conf /alert/conf
COPY ./static /alert/static
COPY ./views /alert/views
COPY ./alertCenter /alert/alertCenter

WORKDIR /alert

RUN chmod 655 /alert/alertCenter

EXPOSE 8888

ENTRYPOINT ["/alert/alertCenter"]
