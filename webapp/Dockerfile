FROM debian:jessie

RUN apt-key update && \
    apt-get update && \
    apt-get install -y ca-certificates && \
    mkdir /var/app -p && \
    mkdir /var/data -p

COPY app_linux_amd64 /bin/app

COPY static /var/app/static

WORKDIR /var/app

ENTRYPOINT [ "/bin/app", "-allow-share" ]
