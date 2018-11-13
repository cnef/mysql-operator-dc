FROM mysql/mysql-server:8.0.13

COPY bin/linux_amd64/mysql-agent /

USER mysql

ENV HOME /var/lib/mysql-agent

ENTRYPOINT ["/mysql-agent"]
