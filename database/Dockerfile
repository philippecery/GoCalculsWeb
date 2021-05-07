FROM mariadb:10.5

RUN mkdir /etc/mysql/tls

ADD ./scripts /docker-entrypoint-initdb.d

ENV MYSQL_RANDOM_ROOT_PASSWORD true