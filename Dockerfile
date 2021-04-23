# See https://github.com/phusion/baseimage-docker
FROM phusion/baseimage:focal-1.0.0-alpha1-amd64

RUN apt-get update -y && apt-get install --no-install-recommends -y -q golang build-essential git unzip && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ADD . /app

WORKDIR /app
RUN make all

EXPOSE 8080

# Register migrations on startup
RUN mkdir -p /etc/my_init.d
ADD ./etc/migrate /etc/my_init.d/10_migrate

# Register API service in runit
ADD ./etc/run /etc/my_init.d/99_api
