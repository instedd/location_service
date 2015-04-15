# See https://github.com/phusion/baseimage-docker
FROM phusion/baseimage

# Register API service in runit
RUN mkdir /etc/service/api
ADD ./etc/run /etc/service/api/run

# Register migrations on startup
RUN mkdir -p /etc/my_init.d
ADD ./etc/migrate /etc/my_init.d/migrate

# Load binary files
ADD ./etc/debian /app/bin

# Expose API port
EXPOSE 8080
