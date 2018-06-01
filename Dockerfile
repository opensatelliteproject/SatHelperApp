FROM ubuntu:16.04

MAINTAINER Lucas Teske <lucas@teske.com.br>

ARG DEBIAN_FRONTEND=noninteractive

RUN echo "deb http://nginx.org/packages/ubuntu/ xenial nginx" > /etc/apt/sources.list.d/nginx.list
RUN echo "deb-src http://nginx.org/packages/ubuntu/ xenial nginx" >> /etc/apt/sources.list.d/nginx.list

RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys ABF5BD827BD9BF62

RUN apt-get update && apt-get install -y --no-install-recommends nginx php-fpm php-mysql php-mbstring php-xml php-imap php-zip && rm -rf /var/lib/apt/lists/*

RUN sed -i 's/user  nginx;/user  nginx;\ndaemon off;/' /etc/nginx/nginx.conf
RUN sed -i 's/error_log .*;/error_log  \/dev\/stderr warn;/g' /etc/nginx/nginx.conf
RUN sed -i 's/access_log .*;/access_log  \/dev\/stdout main;/g' /etc/nginx/nginx.conf
RUN sed -i 's/#gzip  on;/gzip  on;/g' /etc/nginx/nginx.conf

RUN sed -i 's/www-data/nginx/g' /etc/php/7.0/fpm/pool.d/www.conf
RUN sed -i 's/;php_flag[display_errors].*/php_flag[display_errors] = On/g' /etc/php/7.0/fpm/pool.d/www.conf
RUN sed -i 's/;php_admin_flag[log_errors].*/php_admin_flag[log_errors] = On/g' /etc/php/7.0/fpm/pool.d/www.conf
RUN sed -i 's/;cgi.fix_pathinfo=1/cgi.fix_pathinfo=0/g' /etc/php/7.0/fpm/php.ini
RUN sed -i 's/display_errors = Off/display_errors = On/g' /etc/php/7.0/fpm/php.ini

RUN mkdir -p /var/www/sathelperapp

COPY . /var/www/sathelperapp/

RUN chown -R nginx.nginx /var/www/

WORKDIR /opt
COPY sathelperapp.conf /etc/nginx/conf.d/default.conf
COPY run.sh .
CMD /opt/run.sh
