FROM alpine:3.4
MAINTAINER Barnett <zengqingguo@goyoo.com>

COPY ./conf /alert/conf
COPY ./static /alert/static
COPY ./views /alert/views
COPY ./alertCenter /alert/alertCenter

RUN chmod 655 /alert/alertCenter && apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

EXPOSE 8888

ENTRYPOINT ["/alert/alertCenter"]
