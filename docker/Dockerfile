FROM alpine:latest
MAINTAINER Alexander Zaytsev "thebestzorro@yandex.ru"
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates tzdata
ADD webpw.alpine /bin/webpw
ADD index.html /root/index.html
RUN chmod 0755 /bin/webpw
EXPOSE 30080
ENTRYPOINT ["webpw"]
CMD ["-index", "/root/index.html"]