#!/usr/bin/env bash

cp ./gin-jail.conf /etc/fail2ban/jail.d/
cp ./gin-app.conf /etc/fail2ban/filter.d/
ln -s /etc/fail2ban/jail.d/gin-jail.conf "$(pwd)/gin-jail-link.conf"
ln -s /etc/fail2ban/filter.d/gin-app.conf "$(pwd)/gin-app-link.conf"
