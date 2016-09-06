FROM alpine:3.4
MAINTAINER Barnett <zengqingguo@goyoo.com>

COPY ./conf /alert/conf
COPY ./static /alert/static
COPY ./views /alert/views
COPY ./alertCenter /alert/alertCenter

CMD chmod 655 /alert/alertCenter

EXPOSE 8888

ENTRYPOINT ["/alert/alertCenter"] 
