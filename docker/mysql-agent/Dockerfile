FROM mysql/mysql-server:8.0.12

COPY bin/linux_amd64/mysql-agent /

USER mysql

ENV HOME /var/lib/mysql-agent

ENTRYPOINT ["/mysql-agent"]
