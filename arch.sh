#!/bin/bash


YESTERDAY=`date -d '1 day ago' '+%Y-%m-%d'`
echo $YESTERDAY
FILE="/var/www/traffic/logs/${YESTERDAY}.log"
ZIP="/var/www/traffic/logs/${YESTERDAY}.zip"
echo $FILE
echo $ZIP
zip $ZIP $FILE
rm $FILE