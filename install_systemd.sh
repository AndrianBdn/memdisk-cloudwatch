#!/bin/sh
set -e

NAME=memdisk-cloudwatch
VERSION=v0.9.2
BIN_URL=https://github.com/AndrianBdn/$NAME/releases/download/$VERSION/$NAME-x86_64.gz
BIN_PATH=/usr/local/bin/$NAME
TMP_BIN=/tmp/$NAME
SYSTEMD_UNIT=/etc/systemd/system/$NAME.service

curl -L $BIN_URL | gunzip -c > $TMP_BIN
chmod 755 $TMP_BIN

set +e
if ! $TMP_BIN -runonce 1
then
    echo "got error when trying $NAME"
    echo "not continuing to install"
    exit 1
fi
set -e

rm -f $BIN_PATH
cp $TMP_BIN $BIN_PATH

cat > $SYSTEMD_UNIT <<- EOM
[Unit]
Description=memdisk-cloudwatch: reporting mem&disk to cloudwatch

[Service]
Environment=HOME=/root
; HOME environment is needed when AWS credentials are stored in /root/.aws/credentials file 
ExecStart=/usr/local/bin/memdisk-cloudwatch
Restart=on-failure
RestartSec=300

[Install]
WantedBy=multi-user.target
EOM

set +e
systemctl daemon-reload
if ps aux | grep $NAME | grep -v grep
then
systemctl restart $NAME
else
systemctl enable $NAME
systemctl start $NAME
fi

touch /root/.$NAME-$VERSION
