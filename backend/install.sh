#!/bin/bash

set -u

APP=code-snippets
SESSION_TOKEN=$(openssl rand -hex 32)"_t"

checkError() {
    code=$?
    if [[ "$code" != "0" ]]; then
        echo "FAIL: $1"
        exit $code
    fi
}

if [[ ! -f $APP ]]; then
  echo "The $APP binary is missing"
  exit 1
fi

useradd $APP
checkError "Add user"

mkdir -p /var/lib/$APP/ /etc/$APP/
checkError "Create folders"

sed "s/{GEN_TOKEN}/${SESSION_TOKEN}/" env.dist > /etc/$APP/env
checkError "Create service env file"

cp -f $APP /usr/bin/$APP
checkError "Install the binary"

cp -R $APP views /var/lib/$APP/
checkError "Install views"

chown $APP:$APP -R /usr/bin/$APP /var/lib/$APP

cp $APP.service /etc/systemd/system/$APP.service
checkError "Create SystemD service"

systemctl daemon-reload
checkError "SystemD daemon reload"

echo "Done"
exit 0
