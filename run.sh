#!/bin/bash

if [ -e /run/secrets/ ]
then
  for i in /run/secrets/*
  do
    VARNAME=`basename ${i}`
    echo "Importing Secret ${VARNAME}"
    declare $VARNAME=`cat ${i}`
  done
fi


touch /var/www/sathelperapp/logs/errors
chmod u+rwx,g+rwx /var/www/sathelperapp/temp
chmod u+rwx,g+rwx /var/www/sathelperapp/logs

sed -i "s/upload_max_filesize = 2M/upload_max_filesize = ${MAX_UPLOAD_SIZE}/g" /etc/php/7.0/fpm/php.ini
sed -i "s/post_max_size = 8M/post_max_size = ${MAX_UPLOAD_SIZE}/g" /etc/php/7.0/fpm/php.ini

chown nginx.nginx /var/www/sathelperapp/config/config.inc.php
chown -R nginx.nginx /var/www/sathelperapp/logs/

service php7.0-fpm start

/usr/sbin/nginx & tail -f /var/log/php7.0-fpm.log & tail -f /var/www/sathelperapp/logs/errors