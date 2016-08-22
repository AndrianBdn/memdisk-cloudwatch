#!/bin/sh
set -e

NAME=memdisk-cloudwatch
VERSION=v0.9.2
BIN_URL=https://github.com/AndrianBdn/memdisk-cloudwatch/releases/download/$VERSION/memdisk-cloudwatch-x86_64.gz
BIN_PATH=/usr/local/bin/$NAME
TMP_BIN=/tmp/$NAME

curl -L $BIN_URL | gunzip -c > $TMP_BIN
chmod 755 $TMP_BIN

set +e
if ! $TMP_BIN -runonce 1
then
    echo "got error when trying memdisk-cloudwatch"
    echo "not continuing to install"
    exit 1
fi
set -e

rm -f $BIN_PATH
cp $TMP_BIN $BIN_PATH

TMPCRON=$(mktemp /tmp/memdisk-cw-cron.XXXXXX)
crontab -l | grep -v $NAME > $TMPCRON || true 

echo "*/5 * * * * /usr/local/bin/$NAME -crontab 1" >> $TMPCRON
cat $TMPCRON | crontab - 
rm -f $TMPCRON

touch /root/.memdisk-cloudwatch-$VERSION

