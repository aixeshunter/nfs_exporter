FROM          alpine:latest
LABEL         maintainer="Aixes Hunter <aixeshunter@gmail.com>"

# ENV http_proxy "http://127.0.0.1:8080"

# install nfs-utils
RUN apk update && apk add --update nfs-utils && rm -rf /var/cache/apk/*

RUN rm /sbin/halt /sbin/poweroff /sbin/reboot

#RUN unset http_proxy

# start nfs_exporter
ADD nfs_exporter /usr/local/bin/nfs_exporter

RUN chmod 775 /usr/local/bin/nfs_exporter

EXPOSE 9689

ENTRYPOINT  [ "/usr/local/bin/nfs_exporter" ]