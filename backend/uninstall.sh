#!/bin/bash

set -u

APP=code-snippets

if [[ ! -d /usr/bin/$APP ]]; then
  rm -f /usr/bin/$APP
fi

systemctl stop $APP

systemctl disable $APP

rm -rf /var/lib/$APP /etc/$APP

rm -f /etc/systemd/system/$APP.service

systemctl daemon-reload

userdel -f -r $APP

echo "Done"
exit 0
