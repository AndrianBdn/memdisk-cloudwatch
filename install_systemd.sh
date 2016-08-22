#!/bin/sh
set -e

NAME=memdisk-cloudwatch
BIN_URL=https://github.com/AndrianBdn/memdisk-cloudwatch/releases/download/v0.9.2/memdisk-cloudwatch-x86_64.gz
BIN_PATH=/usr/local/bin/$NAME
TMP_BIN=/tmp/$NAME
SYSTEMD_UNIT=/etc/systemd/system/$NAME.service

curl -L $BIN_URL | gunzip -c > $TMP_BIN
chmod 755 $TMP_BIN

set +e
if ! $TMP_BIN -runonce 1
then
    echo "got error when trying memdisk-cloudwatch"
    echo "not continuing to install"
fi
set -e

rm -f $BIN_PATH
cp $TMP_BIN $BIN_PATH

cat > $SYSTEMD_UNIT <<- EOM
[Unit]
Description=memdisk-cloudwatch: reporting mem&disk to cloudwatch

[Service]
ExecStart=/usr/local/bin/memdisk-cloudwatch
Restart=on-abort

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
