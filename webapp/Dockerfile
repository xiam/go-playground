FROM debian:jessie

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY app_linux_amd64 /bin/app

RUN mkdir /var/app -p

RUN mkdir /var/data -p

COPY static /var/app/static

WORKDIR /var/app

ENTRYPOINT [ "/bin/app" ]
